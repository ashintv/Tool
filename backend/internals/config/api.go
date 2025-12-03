package config

import (
	"github.com/gin-gonic/gin"
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

func LoadConfig() *Config {
	Env := gin.Mode()
	if Env != gin.ReleaseMode {

		// dev mode return developer config
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
	// add a function for validate data
	return &Cnf

}
