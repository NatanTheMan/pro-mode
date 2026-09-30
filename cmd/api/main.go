package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/logan-dev-x/pro-mode/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Connect()
	if err != nil {
		fmt.Println(err)
	}
}
