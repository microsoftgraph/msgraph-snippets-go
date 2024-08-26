// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

import (
	"context"
	"log"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

func RunRequestSamples(graphClient *graph.GraphServiceClient) {
	// Create a new message
	msg := models.NewMessage()
	subject := "Temporary"
	msg.SetSubject(&subject)
	tempMessage, err := graphClient.Me().Messages().Post(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Error creating message: %v\n", err)
	}
	messageId := tempMessage.GetId()

	filterValue := "resourceProvisioningOptions/Any(x:x eq 'Team')"
	query := groups.GroupsRequestBuilderGetQueryParameters{
		Filter: &filterValue,
	}

	options := groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}
	// Get a team to update
	teams, err := graphClient.Groups().Get(context.Background(), &options)
	if err != nil {
		log.Fatalf("Error getting teams: %v\n", err)
	}
	teamId := teams.GetValue()[0].GetId()

	MakeReadRequest(graphClient)
	MakeSelectRequest(graphClient)
	MakeListRequest(graphClient)
	MakeItemByIdRequest(graphClient, *messageId)
	MakeExpandRequest(graphClient, *messageId)
	MakeDeleteRequest(graphClient, *messageId)
	MakeCreateRequest(graphClient)
	MakeUpdateRequest(graphClient, *teamId)
	MakeHeadersRequest(graphClient)
	MakeQueryParametersRequest(graphClient)
}

func MakeReadRequest(graphClient *graph.GraphServiceClient) models.Userable {
	// <ReadRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me
	result, _ := graphClient.Me().Get(context.Background(), nil)
	// </ReadRequestSnippet>

	return result
}

func MakeSelectRequest(graphClient *graph.GraphServiceClient) models.Userable {
	// <SelectRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me?$select=displayName,jobTitle

	// import github.com/microsoftgraph/msgraph-sdk-go/users
	query := users.UserItemRequestBuilderGetQueryParameters{
		Select: []string{"displayName", "jobTitle"},
	}

	options := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}

	result, _ := graphClient.Me().Get(context.Background(), &options)
	// </SelectRequestSnippet>

	return result
}

func MakeListRequest(graphClient *graph.GraphServiceClient) models.MessageCollectionResponseable {
	// <ListRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me/messages?
	// $select=subject,sender&$filter=subject eq 'Hello world'

	// import github.com/microsoftgraph/msgraph-sdk-go/users
	filterValue := "subject eq 'Hello world'"
	query := users.ItemMessagesRequestBuilderGetQueryParameters{
		Select: []string{"subject", "sender"},
		Filter: &filterValue,
	}

	options := users.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}

	result, _ := graphClient.Me().Messages().
		Get(context.Background(), &options)
	// </ListRequestSnippet>

	return result
}

func MakeItemByIdRequest(graphClient *graph.GraphServiceClient, messageId string) models.Messageable {
	// <ItemByIdRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me/messages/{message-id}
	// messageId is a string containing the id property of the message
	result, _ := graphClient.Me().Messages().
		ByMessageId(messageId).Get(context.Background(), nil)
	// </ItemByIdRequestSnippet>

	return result
}

func MakeExpandRequest(graphClient *graph.GraphServiceClient, messageId string) models.Messageable {
	// <ExpandRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me/messages/{message-id}?$expand=attachments

	// import github.com/microsoftgraph/msgraph-sdk-go/users
	expand := users.ItemMessagesMessageItemRequestBuilderGetQueryParameters{
		Expand: []string{"attachments"},
	}

	options := users.ItemMessagesMessageItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &expand,
	}
	// messageId is a string containing the id property of the message
	result, _ := graphClient.Me().Messages().
		ByMessageId(messageId).Get(context.Background(), &options)
	// </ExpandRequestSnippet>

	return result
}

func MakeDeleteRequest(graphClient *graph.GraphServiceClient, messageId string) error {
	// <DeleteRequestSnippet>
	// DELETE https://graph.microsoft.com/v1.0/me/messages/{message-id}
	// messageId is a string containing the id property of the message
	err := graphClient.Me().Messages().
		ByMessageId(messageId).Delete(context.Background(), nil)
	// </DeleteRequestSnippet>

	return err
}

func MakeCreateRequest(graphClient *graph.GraphServiceClient) models.Calendarable {
	// <CreateRequestSnippet>
	// POST https://graph.microsoft.com/v1.0/me/calendars

	calendar := models.NewCalendar()
	name := "Volunteer"
	calendar.SetName(&name)

	result, _ := graphClient.Me().Calendars().Post(context.Background(), calendar, nil)
	// </CreateRequestSnippet>

	return result
}

func MakeUpdateRequest(graphClient *graph.GraphServiceClient, teamId string) {
	// <UpdateRequestSnippet>
	// PATCH https://graph.microsoft.com/v1.0/teams/{team-id}

	funSettings := models.NewTeamFunSettings()
	allowGiphy := true
	funSettings.SetAllowGiphy(&allowGiphy)
	giphyRating := models.STRICT_GIPHYRATINGTYPE
	funSettings.SetGiphyContentRating(&giphyRating)

	team := models.NewTeam()
	team.SetFunSettings(funSettings)

	graphClient.Teams().ByTeamId(teamId).Patch(context.Background(), team, nil)
	// </UpdateRequestSnippet>
}

func MakeHeadersRequest(graphClient *graph.GraphServiceClient) models.EventCollectionResponseable {
	// <HeadersRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me/events

	// import abstractions "github.com/microsoft/kiota-abstractions-go"
	headers := abstractions.NewRequestHeaders()
	headers.Add("Prefer", "outlook.timezone=\"Pacific Standard Time\"")

	// import github.com/microsoftgraph/msgraph-sdk-go/users
	options := users.ItemEventsRequestBuilderGetRequestConfiguration{
		Headers: headers,
	}

	result, _ := graphClient.Me().Events().Get(context.Background(), &options)
	// </HeadersRequestSnippet>

	return result
}

func MakeQueryParametersRequest(graphClient *graph.GraphServiceClient) models.EventCollectionResponseable {
	// <QueryParametersRequestSnippet>
	// GET https://graph.microsoft.com/v1.0/me/calendarView?
	// startDateTime=2023-06-14T00:00:00Z&endDateTime=2023-06-15T00:00:00Z

	startDateTime := "2023-06-14T00:00:00"
	endDateTime := "2023-06-15T00:00:00Z"

	// import github.com/microsoftgraph/msgraph-sdk-go/users
	query := users.ItemCalendarViewRequestBuilderGetQueryParameters{
		StartDateTime: &startDateTime,
		EndDateTime:   &endDateTime,
	}

	options := users.ItemCalendarViewRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}

	result, _ := graphClient.Me().CalendarView().Get(context.Background(), &options)
	// </QueryParametersRequestSnippet>

	return result
}
