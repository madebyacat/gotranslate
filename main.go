package main

import (
	"context"
	"errors"
	"fmt"
	"gotranslate/api/rest"
	"gotranslate/core/contracts"
	"gotranslate/core/messages"
	"gotranslate/core/queue"
	"gotranslate/core/repository"
	"gotranslate/core/translators"
	"gotranslate/models"
	"log"
	"time"

	"cloud.google.com/go/translate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	loadConfig()

	repo, cleanup := initializeRepository()
	defer cleanup()

	translator := initializeTranslationService()

	queueClient := initializeQueue(repo, translator)
	defer queueClient.Close()

	auth := loadAuthenticationConfig()
	var router = rest.NewRouter(repo, translator, queueClient, auth)

	router.Run("localhost:3000")
}

func loadConfig() {
	// viper.SetConfigName("config.local")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func loadAuthenticationConfig() models.AuthConfig {
	auth := models.NewAuthConfig(
		viper.GetBool("auth.skip_authentication"),
		viper.GetString("auth.jwt_key"),
		viper.GetString("auth.username"),
		viper.GetString("auth.password"))
	return auth
}

func initializeRepository() (repo contracts.ResoureRepository, cleanup func()) {
	persistenceMedium := viper.GetString("persistence")
	if persistenceMedium == "file" {
		file := viper.GetString("file")
		repo := repository.NewResourceFile(file)

		if err := repo.Init(); err != nil {
			log.Fatal(err)
		}

		return repo, func() {}
	} else if persistenceMedium == "postgres" {
		pool, err := pgxpool.New(context.Background(), viper.GetString("database.connection_string"))
		if err != nil {
			log.Fatal(err)
		}

		if err := pool.Ping(context.Background()); err != nil {
			log.Fatal(err)
		}

		repo := repository.NewResourceSql(pool)

		if err := repo.Init(); err != nil {
			log.Fatal(err)
		}

		return repo, func() { pool.Close() }
	} else if persistenceMedium == "gorm" {
		db, err := gorm.Open(postgres.Open(viper.GetString("database.connection_string")), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		repo := repository.NewResourceGorm(db)

		if err := repo.Init(); err != nil {
			log.Fatal(err)
		}

		// Configure the connection pool
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get generic database object")
		}

		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)

		return repo, func() { sqlDB.Close() }
	} else {
		panic(fmt.Errorf("no repository implemented for selection %v", persistenceMedium))
	}
}

func initializeTranslationService() contracts.Translator {
	if service := viper.GetString("translation"); service == "" {
		log.Fatal(errors.New("translation not configured"))
	} else if service == "google" {
		apiKey := viper.GetString("google_api_key")
		if apiKey == "" {
			fmt.Println("***********api key not set**********")
		}

		client, err := translate.NewClient(context.Background(), option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
		}

		return translators.NewGoogle(client)
	} else if service == "fake" {
		return &translators.Fake{}
	} else {
		log.Fatal(errors.New("unsupported translation service"))
	}

	return nil
}

func initializeQueue(repo contracts.ResoureRepository, translator contracts.Translator) contracts.QueueService {
	if queueType := viper.GetString("queue.type"); queueType == "" {
		log.Fatal(errors.New("queue not configured"))
	} else if queueType != "rabbitmq" {
		log.Fatal(errors.New("unsupported queue"))
	}

	url, queueName := viper.GetString("queue.url"), viper.GetString("queue.queue_name")
	if url == "" || queueName == "" {
		log.Fatal("****** queue configuration missing ******")
	}

	queueClient, err := queue.NewRabbitMQ(url, queueName)
	if err != nil {
		log.Fatal(err)
	}

	strategy := messages.GetMessageHandlers(repo, translator)

	queueClient.Consume(strategy)

	return queueClient
}
