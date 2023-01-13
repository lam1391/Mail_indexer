package envvariables

import (
	"log"

	"github.com/joho/godotenv"
)

func GetEnvVariables() {
	err := godotenv.Load("./variables.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

}
