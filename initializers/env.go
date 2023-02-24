package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	err := godotenv.Load()
	return err
}
