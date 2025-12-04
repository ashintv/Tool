package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Config struct {

	// cors configs
	AllowOriginFunc  func(origin string) bool
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool

	//PORT
	Port               string
	DATABASE_DSN       string
	USER_JWT_SECRET    string
	MACHINE_JWT_SECRET string
}

//"host=localhost user=user password=password dbname=dbname port=5432"

var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}

func LoadConfig(logger *zerolog.Logger) *Config {
	Env := gin.Mode()
	log.Println("\n-----------", Env, "----------")
	fmt.Print("\n\n")
	if Env != gin.ReleaseMode {
		// dev mode return developer config
		logger.Warn().Msg("Running in Development Mode, using default config values")
		return &Config{
			AllowOriginFunc: func(origin string) bool {
				return true // allow ALL origins
			},
			AllowMethods:       []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:       []string{"*"},
			ExposeHeaders:      []string{"*"},
			AllowCredentials:   true,
			Port:               ":8080",
			DATABASE_DSN:       "host=localhost user=user password=password dbname=dbname port=5432",
			USER_JWT_SECRET:    "jwt_secret_for_users",
			MACHINE_JWT_SECRET: "jwt_secret_for_users",
		}
	}
	var Cnf Config

	Cnf.AllowOriginFunc = func(origin string) bool {
		return true // allow ALL origins
	}

	Cnf.AllowMethods = methods
	Cnf.AllowHeaders = []string{"*"}
	Cnf.ExposeHeaders = []string{"*"}
	Cnf.AllowCredentials = true

	Port := os.Getenv("PORT")
	if Port == "" {
		logger.Warn().Msg("PORT not set in environment variables, defaulting to 8080")
		Port = "8080"
	} else {
		logger.Info().Msg("Using PORT from environment variable")
		// validate port
		_, err := fmt.Sscanf(Port, "%d", new(int))
		if err != nil {
			logger.Fatal().Err(err).Msg("Invalid PORT value")
			panic("Invalid Port")
		}
	}

	Cnf.DATABASE_DSN = os.Getenv("DATABASE_DSN")
	if Cnf.DATABASE_DSN == "" {
		logger.Fatal().Msg("DATABASE_DSN not set in environment variables")
		panic("DATABASE_DSN not set")
	}
	logger.Info().Msg("Using DATABASE_DSN from environment variable")

	Cnf.USER_JWT_SECRET = os.Getenv("USER_JWT_SECRET")
	if Cnf.USER_JWT_SECRET == "" {
		logger.Fatal().Msg("USER_JWT_SECRET not set in environment variables")
		panic("USER_JWT_SECRET not set")
	}
	logger.Info().Msg("Using USER_JWT_SECRET from environment variable")

	if Cnf.MACHINE_JWT_SECRET = os.Getenv("MACHINE_JWT_SECRET"); Cnf.MACHINE_JWT_SECRET == "" {
		logger.Fatal().Msg("MACHINE_JWT_SECRET not set in environment variables")
		panic("MACHINE_JWT_SECRET not set")
	}
	logger.Info().Msg("Using MACHINE_JWT_SECRET from environment variable")
	Cnf.Port = ":" + Port
	return &Cnf

}
