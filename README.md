## Go Translate

A web API to store, retrieve, and manage text Resources, with automatic and manual Translation options.

The purpose of building this api was to learn Go, building a useful feature/app.
Note that my experience is mostly on C# so some naming/structure might seem off because I'm not used to GoLang yet :D

### Running the service
1. Run `docker-compose up -d` to create/run the required containers: PostgreSQL, RabbitMQ.
   1. Change the ports or anything you need in `docker-compose.yaml`
2. config.json contains values that should work, but you can change them to your requirements
  1. If you use the `docker-compose.yaml` make sure the values you change there are changed in the Config too.
  2. If you want to use real keys, add a `config.local.json` and use it in main.go, `*.local.json` is in `.gitignore`.
3. In root folder execute `go run .`
4. To run all the tests, in root folder execute `go test ./...` - Note: some need docker running 

### Persistence/Repositories
The repository to use is the one using Gorm ORM, the other ones are just for playing around for educational purposes. Some features are not implemented on File and raw SQL repositories.
1. File Repository configuration: `"persistence": "file", "file": "YOURFILENAME.txt"`.
2. Raw SQL Repository configuration: `"persistence": "postgres", "database": { "connection_string": "YOURCONNECTIONSTRING" }`.
3. Gorm Repository configuration: `"persistence": "gorm", "database": { "connection_string": "YOURCONNECTIONSTRING" }`.

### Translation
There are two translation options: Google Translate using Google cloud services, and Fake Translator that generates random words to emulate translation.
1. Fake Translation configuration: `"translation": "fake"`
2. Google Translation configuration: `"translation": "google", "google_api_key": "YOUR_KEY"` (note: google charges for the usage).

### Queueing
Translation calling a service might take some time if you have many resources, so for performance (and education) it implements RabbitMQ
1. It's settings are in `config.json`, where you can change the settings in `"queue": { "type": "rabbitmq" ... }` for your own RabbitMQ instance.
2. There's only 1 queue implemented, and it's required to start up the application.

### Authentication
I made an authentication middleware to learn, but I haven't expanded on it much. It worked last time I tried it.
1. To bypass the authentication `config.json` should be `"auth":{ "skip_authentication": true }`
