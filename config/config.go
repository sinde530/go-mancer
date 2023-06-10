package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port string `env:"PORT,default=8080"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	// err = envconfig.Process(ctx, &config)
	err = envconfig.Process(ctx, config)

	if err != nil {
		return
	}

	return
}
