package main

import (
	"fmt"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	_, err = config.Load()
	if err != nil {
		fmt.Println(err)
	}
}
