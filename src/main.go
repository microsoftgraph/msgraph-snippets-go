// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sdksnippets/graphhelper"
	"sdksnippets/snippets"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Microsoft Graph Go SDK Snippets")
	fmt.Println()

	logger := log.New(os.Stdout, "graph-debug: ", log.Ldate|log.Ltime)

	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	userClient, err := graphhelper.NewUserGraphServiceClient(logger)
	if err != nil {
		log.Fatalf("Error creating user client: %v\n", err)
	}

	user, err := userClient.Me().Get(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error getting user: %v\n", err)
	}

	fmt.Printf("Hello %s!\n", *user.GetDisplayName())

	var choice int64 = -1

	for {
		fmt.Println("Please choose one of the following options:")
		fmt.Println("0. Exit")
		fmt.Println("1. Run batch samples")

		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			choice = -1
		}

		switch choice {
		case 0:
			// Exit the program
			fmt.Println("Goodbye...")
		case 1:
			// Display access token
			snippets.RunBatchSamples(userClient)
		default:
			fmt.Println("Invalid choice! Please try again.")
		}

		if choice == 0 {
			break
		}
	}
}
