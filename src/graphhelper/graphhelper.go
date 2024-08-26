// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package graphhelper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graphdebug "github.com/jasonjoh/msgraph-sdk-go-debug-logger"
	khttp "github.com/microsoft/kiota-http-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	auth "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

func NewUserGraphServiceClient(logger *log.Logger) (*graph.GraphServiceClient, error) {
	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	scopes := strings.Split(os.Getenv("GRAPH_USER_SCOPES"), ",")
	debug, err := strconv.ParseBool(os.Getenv("ENABLE_GRAPH_LOG"))
	if err != nil {
		debug = false
	}

	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		TenantID: tenantId,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts(credential, scopes, []string{"graph.microsoft.com"})
	if err != nil {
		return nil, err
	}

	if debug {
		return NewDebugGraphServiceClient(authProvider, logger)
	} else {
		adapter, err := graph.NewGraphRequestAdapter(authProvider)
		if err != nil {
			return nil, err
		}

		client := graph.NewGraphServiceClient(adapter)
		return client, nil
	}
}

func NewDebugGraphServiceClient(authProvider *auth.AzureIdentityAuthenticationProvider, logger *log.Logger) (*graph.GraphServiceClient, error) {
	showTokens, err := strconv.ParseBool(os.Getenv("GRAPH_LOG_TOKENS"))
	if err != nil {
		showTokens = false
	}
	showPayloads, err := strconv.ParseBool(os.Getenv("GRAPH_LOG_PAYLOADS"))
	if err != nil {
		showPayloads = false
	}

	clientOptions := graph.GetDefaultClientOptions()
	middleware := graphcore.GetDefaultMiddlewaresWithOptions(&clientOptions)
	debugMiddleware := graphdebug.NewGraphDebugLogMiddleware(logger, showTokens, showPayloads)
	allMiddleware := append(middleware, debugMiddleware)
	httpClient := khttp.GetDefaultClient(allMiddleware...)

	adapter, err := graph.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider, nil, nil, httpClient)
	if err != nil {
		return nil, err
	}

	client := graph.NewGraphServiceClient(adapter)
	return client, nil
}
