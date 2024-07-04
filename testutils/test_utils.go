package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	pgtestcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SpinUpContainer(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	pgContainer, err := pgtestcontainers.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		pgtestcontainers.WithDatabase("testdb"),
		pgtestcontainers.WithUsername("testuser"),
		pgtestcontainers.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	teardown := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	return db, teardown
}
