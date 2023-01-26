package main

import (
	"fmt"
	"log"

	"github.com/fazrithe/pkk-backend-v2/initializers"
	"github.com/fazrithe/pkk-backend-v2/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Institution{})
	fmt.Println("? Migration complete")
}
