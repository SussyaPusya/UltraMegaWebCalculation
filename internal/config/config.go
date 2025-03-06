package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port string `yaml:"PORT" env:"PORT" env-default:"8080"`
	EnvConfig
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}

type EnvConfig struct {
	TIME_ADDITION_MS        int `yaml:"TIME_ADDITION_MS" env:"TIME_ADDITION_MS" env-default:"5"`
	TIME_SUBTRACTION_MS     int `yaml:"TIME_SUBTRACTION_MS" env:"TIME_SUBTRACTION_MS" env-default:"5"`
	TIME_MULTIPLICATIONS_MS int `yaml:"TIME_MULTIPLICATIONS_MS" env:"TIME_MULTIPLICATIONS_MS" env-default:"5"`
	TIME_DIVISIONS_MS       int `yaml:"TIME_DIVISIONS_MS" env:"TIME_DIVISIONS_MS" env-default:"5"`
	COMPUTING_POWER         int `yaml:"COMPUTING_POWER" env:"COMPUTING_POWER" env-default:"3"`
}
