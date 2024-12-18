package main

import (
	"fmt"
	"log"

	"github.com/MosinEvgeny/authservice/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("config loaded:", cfg)
}
