package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	http.HandleFunc("/check-image", checkImageHandler)
	fmt.Println("Server is listening on port 8129...")

	if err := http.ListenAndServe(":8129", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func checkImageHandler(w http.ResponseWriter, r *http.Request) {

}
