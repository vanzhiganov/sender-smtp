package main

// Configuration структура конфигурационного файла
type Configuration struct {
	Application struct {
		Listen       string `yaml:"listen"`
		SecretKey    string `yaml:"secret_key"`
		TemplateFile string `yaml:"template_file"`
	}
	Postgresql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
	// see also: https://github.com/streadway/amqp/blob/master/_examples/simple-consumer/consumer.go
	// Rabbitmq struct {
	// 	Enabled bool `yaml:"enabled"`
	// }
	Sentry struct {
		Enabled bool   `yaml:"enabled"`
		DSN     string `yaml:"dsn"`
	}
	SMTP struct {
		Server string `yaml:"server"`
		Port   int    `yaml:"port"`
		UseTLS bool   `yaml:"use_tls"`
		UseSSL bool   `yaml:"use_ssl"`
		Sender struct {
			Login    string `yaml:"login"`
			Password string `yaml:"password"`
		}
	}
}
