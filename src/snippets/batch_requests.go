// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

import (
	"context"
	"fmt"
	"log"
	"time"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/thlib/go-timezone-local/tzlocal"
)

func RunBatchSamples(graphClient *graph.GraphServiceClient) {
	SimpleBatch(graphClient)
	DependentBatch(graphClient)
}

func SimpleBatch(graphClient *graph.GraphServiceClient) {
	// <SimpleBatchSnippet>
	// Use the request builder to generate a regular
	// request to /me
	meRequest, err := graphClient.Me().
		ToGetRequestInformation(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error creating GET /me request: %v\n", err)
	}

	now := time.Now()
	nowMidnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.Local)

	viewStart := nowMidnight.UTC().Format(time.RFC3339)
	viewEnd := nowMidnight.AddDate(0, 0, 1).UTC().Format(time.RFC3339)

	query := users.ItemCalendarViewRequestBuilderGetQueryParameters{
		StartDateTime: &viewStart,
		EndDateTime:   &viewEnd,
		Select:        []string{"subject", "id"},
	}

	// Use the request builder to generate a request
	// to /me/calendarView?startDateTime="start"&endDateTime="end"
	eventsRequest, err := graphClient.Me().
		CalendarView().
		ToGetRequestInformation(context.Background(),
			&users.ItemCalendarViewRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
	if err != nil {
		log.Fatalf("Error creating GET /me/calendarView request: %v\n", err)
	}

	// Build the batch
	batch := graphcore.NewBatchRequest(graphClient.GetAdapter())

	// Using AddBatchRequestStep adds each request as a step
	// with no specified order of execution
	meRequestItem, err := batch.AddBatchRequestStep(*meRequest)
	if err != nil {
		log.Fatalf("Error adding GET /me request to batch: %v\n", err)
	}
	eventsRequestItem, err := batch.AddBatchRequestStep(*eventsRequest)
	if err != nil {
		log.Fatalf("Error adding GET /me/calendarView request to batch: %v\n", err)
	}

	batchResponse, err := batch.Send(context.Background(), graphClient.GetAdapter())
	if err != nil {
		log.Fatalf("Error sending batch: %v\n", err)
	}

	// De-serialize response based on known return type
	user, err := graphcore.GetBatchResponseById[models.Userable](
		batchResponse, *meRequestItem.GetId(), models.CreateUserFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error reading GET /me response: %v\n", err)
	}
	fmt.Printf("Hello %s\n", *(user.GetDisplayName()))

	// For collections, must use the *CollectionResponseable class to deserialize
	events, err := graphcore.GetBatchResponseById[models.EventCollectionResponseable](
		batchResponse, *eventsRequestItem.GetId(),
		models.CreateEventCollectionResponseFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error reading GET /me/calendarView response: %v\n", err)
	}
	fmt.Printf("You have %d events on your calendar today\n", len(events.GetValue()))
	// </SimpleBatchSnippet>
}

func DependentBatch(graphClient *graph.GraphServiceClient) {
	// <DependentBatchSnippet>
	now := time.Now()
	nowMidnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.Local)
	timeZone, _ := tzlocal.RuntimeTZ()

	// 5:00 PM
	startDateTime := nowMidnight.Add(time.Hour * 17)
	// 5:30 PM
	endDateTime := startDateTime.Add(time.Minute * 30)
	graphDateTimeFormat := "2006-01-02T15:04:05"

	// Create event
	newEvent := models.NewEvent()
	subject := "File end-of-day report"
	newEvent.SetSubject(&subject)

	start := models.NewDateTimeTimeZone()
	startString := startDateTime.Format(graphDateTimeFormat)
	start.SetDateTime(&startString)
	start.SetTimeZone(&timeZone)
	newEvent.SetStart(start)

	end := models.NewDateTimeTimeZone()
	endString := endDateTime.Format(graphDateTimeFormat)
	end.SetDateTime(&endString)
	end.SetTimeZone(&timeZone)
	newEvent.SetEnd(end)

	addEventRequest, err := graphClient.Me().
		Events().
		ToPostRequestInformation(context.Background(), newEvent, nil)
	if err != nil {
		log.Fatalf("Error creating POST /me/events request: %v\n", err)
	}

	viewStart := nowMidnight.UTC().Format(time.RFC3339)
	viewEnd := nowMidnight.AddDate(0, 0, 1).UTC().Format(time.RFC3339)

	query := users.ItemCalendarViewRequestBuilderGetQueryParameters{
		StartDateTime: &viewStart,
		EndDateTime:   &viewEnd,
		Select:        []string{"subject", "id"},
	}

	// Use the request builder to generate a request
	// to /me/calendarView?startDateTime="start"&endDateTime="end"
	eventsRequest, err := graphClient.Me().
		CalendarView().
		ToGetRequestInformation(context.Background(),
			&users.ItemCalendarViewRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
	if err != nil {
		log.Fatalf("Error creating GET /me/calendarView request: %v\n", err)
	}

	// Build the batch
	batch := graphcore.NewBatchRequest(graphClient.GetAdapter())

	// Force the requests to execute in order, so that the request for
	// today's events will include the new event created.

	// First request, no dependency
	addEventRequestItem, err := batch.AddBatchRequestStep(*addEventRequest)
	if err != nil {
		log.Fatalf("Error adding POST /me/events request to batch: %v\n", err)
	}

	// Second request, depends on addEventRequestId
	eventsRequestItem, err := batch.AddBatchRequestStep(*eventsRequest)
	if err != nil {
		log.Fatalf("Error creating GET /me/calendarView request to batch: %v\n", err)
	}
	eventsRequestItem.DependsOnItem(addEventRequestItem)

	batchResponse, err := batch.Send(context.Background(), graphClient.GetAdapter())
	if err != nil {
		log.Fatalf("Error sending batch: %v\n", err)
	}

	// De-serialize response based on known return type
	event, err := graphcore.GetBatchResponseById[models.Eventable](
		batchResponse, *addEventRequestItem.GetId(),
		models.CreateEventFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error reading POST /me/events response: %v\n", err)
	}
	fmt.Printf("New event created with ID: %s\n", *(event.GetId()))

	// For collections, must use the *CollectionResponseable class to deserialize
	events, err := graphcore.GetBatchResponseById[models.EventCollectionResponseable](
		batchResponse, *eventsRequestItem.GetId(),
		models.CreateEventCollectionResponseFromDiscriminatorValue)
	if err != nil {
		log.Fatalf("Error reading GET /me/calendarView response: %v\n", err)
	}
	fmt.Printf("You have %d events on your calendar today\n", len(events.GetValue()))
	// </DependentBatchSnippet>
}
