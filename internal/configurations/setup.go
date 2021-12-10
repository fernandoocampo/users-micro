package configurations

import (
	"github.com/caarlos0/env"
)

// Application contains data related to application configuration parameters.
type Application struct {
	DryRun          bool   `env:"DRY_RUN" envDefault:"false"`
	ApplicationPort string `env:"APPLICATION_PORT" envDefault:":8080"`
	Repository      RepositoryParameters
}

// RepositoryParameters contains data related to a repository.
type RepositoryParameters struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DBNAME" envDefault:"postgres"`
}

// Load load application configuration
func Load() (Application, error) {
	cfg := Application{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	repository := RepositoryParameters{}
	if err := env.Parse(&repository); err != nil {
		return cfg, err
	}
	cfg.Repository = repository
	return cfg, nil
}
