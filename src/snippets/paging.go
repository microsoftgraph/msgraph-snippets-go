// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

// <ImportSnippet>
import (
	"context"
	"fmt"
	"log"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// </ImportSnippet>

func RunPagingSamples(graphClient *graph.GraphServiceClient) {
	IterateAllMessages(graphClient)
	IterateAllMessagesWithPause(graphClient)
}

func IterateAllMessages(graphClient *graph.GraphServiceClient) {
	// <PagingSnippet>
	headers := abstractions.NewRequestHeaders()
	headers.Add("Prefer", "outlook.body-content-type=\"text\"")

	var pageSize int32 = 10
	query := users.ItemMessagesRequestBuilderGetQueryParameters{
		Select: []string{"body", "sender", "subject"},
		Top:    &pageSize,
	}

	options := users.ItemMessagesRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: &query,
	}

	result, err := graphClient.Me().Messages().Get(context.Background(), &options)
	if err != nil {
		log.Fatalf("Error getting messages: %v\n", err)
	}

	// Initialize iterator
	pageIterator, err := graphcore.NewPageIterator[*models.Message](
		result,
		graphClient.GetAdapter(),
		models.CreateMessageCollectionResponseFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error creating page iterator: %v\n", err)
	}

	// Any custom headers sent in original request should also be added
	// to the iterator
	pageIterator.SetHeaders(headers)

	// Iterate over all pages
	err = pageIterator.Iterate(
		context.Background(),
		func(message *models.Message) bool {
			fmt.Printf("%s\n", *message.GetSubject())
			// Return true to continue the iteration
			return true
		})
	if err != nil {
		log.Fatalf("Error iterating over messages: %v\n", err)
	}
	// </PagingSnippet>
}

func IterateAllMessagesWithPause(graphClient *graph.GraphServiceClient) {
	// <ResumePagingSnippet>
	var pageSize int32 = 10
	query := users.ItemMessagesRequestBuilderGetQueryParameters{
		Select: []string{"body", "sender", "subject"},
		Top:    &pageSize,
	}

	options := users.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}

	result, err := graphClient.Me().Messages().Get(context.Background(), &options)
	if err != nil {
		log.Fatalf("Error getting messages: %v\n", err)
	}

	// Initialize iterator
	pageIterator, err := graphcore.NewPageIterator[*models.Message](
		result,
		graphClient.GetAdapter(),
		models.CreateMessageCollectionResponseFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error creating page iterator: %v\n", err)
	}

	// Pause iterating after 25
	var count, pauseAfter = 0, 25

	// Iterate over all pages
	err = pageIterator.Iterate(
		context.Background(),
		func(message *models.Message) bool {
			count++
			fmt.Printf("%d: %s\n", count, *message.GetSubject())
			// Once count = 25, this returns false,
			// Which pauses the iteration
			return count < pauseAfter
		})
	if err != nil {
		log.Fatalf("Error iterating over messages: %v\n", err)
	}

	// Pause 5 seconds
	fmt.Printf("Iterated first %d messages, pausing for 5 seconds...\n", pauseAfter)
	time.Sleep(5 * time.Second)
	fmt.Printf("Resuming iteration...\n")

	// Resume iteration
	err = pageIterator.Iterate(
		context.Background(),
		func(message *models.Message) bool {
			count++
			fmt.Printf("%d: %s\n", count, *message.GetSubject())
			// Return true to continue the iteration
			return true
		})
	if err != nil {
		log.Fatalf("Error iterating over messages: %v\n", err)
	}
	// </ResumePagingSnippet>
}

func ManuallyPageAllMessages(graphClient *graph.GraphBaseServiceClient) {
	// <ManualPagingSnippet>
	var pageSize int32 = 10
	query := users.ItemMessagesRequestBuilderGetQueryParameters{
		Top: &pageSize,
	}

	options := users.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}

	result, err := graphClient.Me().Messages().Get(context.Background(), &options)
	if err != nil {
		log.Fatalf("Error getting messages: %v\n", err)
	}

	for {
		for _, message := range result.GetValue() {
			fmt.Printf("%s\n", *message.GetSubject())
		}

		nextPageUrl := result.GetOdataNextLink()
		if nextPageUrl != nil {
			result, err = graphClient.Me().Messages().
				WithUrl(*nextPageUrl).
				Get(context.Background(), nil)
			if err != nil {
				log.Fatalf("Error getting messages: %v\n", err)
			}
		} else {
			break
		}
	}
	// </ManualPagingSnippet>
}
