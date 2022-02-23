package common

import (
	"os"

	"github.com/wonderivan/logger"
)

func MustPassEvn(env string) string {
	v := os.Getenv(env)
	if v == "" {
		logger.Error("Must pass the enviromnet variable:", v)
		os.Exit(1)
	}
	return v
}
