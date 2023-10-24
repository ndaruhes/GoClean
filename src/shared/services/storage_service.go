package services

import (
	"context"
	"fmt"
	"github.com/filestack/filestack-go/client"
	"log"
	"os"
)

func FilestackUploadFile() {
	file, err := os.Open("")
	if err != nil {
		log.Fatal("cannot read the file")
	}
	defer file.Close()

	cli, err := client.NewClient("ANHeLjd1kRGGIWTmt1Cl8z")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}
	fileLink, err := cli.Upload(context.Background(), file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(fileLink.AsString())

}
