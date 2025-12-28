package services

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetLogger() *zerolog.Logger {
	// Pretty console + colors
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC822,
	})
	
	return &log.Logger
}
