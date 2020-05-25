package main

// Configuration структура конфигурационного файла
type Configuration struct {
	Application struct {
		Listen    string `yaml:"listen"`
		SecretKey string `yaml:"secret_key"`
		DB        string `yaml:"db"`
	}
	// see also: https://github.com/streadway/amqp/blob/master/_examples/simple-consumer/consumer.go
	// Rabbitmq struct {
	// 	Enabled bool `yaml:"enabled"`
	// }
	Sentry struct {
		Enabled bool   `yaml:"enabled"`
		DSN     string `yaml:"dsn"`
	}
}
