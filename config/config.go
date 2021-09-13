package config

import "github.com/spf13/viper"

type ConfigurationManager interface {
	GetServerConfig() ServerConfig
	GetPostgresConfig() PostgresConfig
	GetKafkaConfig() KafkaConfig
}

type configurations struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type KafkaConfig struct {
	Broker string
}

func (configurationManager *configurations) GetServerConfig() ServerConfig {
	return configurationManager.Server
}

func (configurationManager *configurations) GetPostgresConfig() PostgresConfig {
	return configurationManager.Postgres
}

func (configurationManager *configurations) GetKafkaConfig() KafkaConfig {
	return configurationManager.Kafka
}

func NewConfigurationManager(configurationFile string, environment string) ConfigurationManager {

	viperInstance := viper.New()
	viperInstance.SetConfigFile(configurationFile)
	configurationError := viperInstance.ReadInConfig()

	if configurationError != nil {
		panic(configurationError)
	}

	configuration := configurations{}
	subEnvironment := viperInstance.Sub(environment)
	err := subEnvironment.Unmarshal(&configuration)
	if err != nil {
		panic(err)
	}

	return &configuration
}
