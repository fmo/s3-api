package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	playerId := r.URL.Query().Get("playerId")
	if playerId == "" {
		http.Error(w, "playerId is required", http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

	svc := s3.New(sess)
	bucket := os.Getenv("S3_BUCKET")
	key := fmt.Sprintf("players/%s.png", playerId)

	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	response := map[string]string{"imgUrl": ""}
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, fmt.Sprintf("failed to check if object exists: %v", err), http.StatusInternalServerError)
		return
	}

	imgUrl := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), bucket, key)
	response["imgUrl"] = imgUrl
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
