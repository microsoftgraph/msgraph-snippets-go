// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

// <ImportSnippet>
import (
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	khttp "github.com/microsoft/kiota-http-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// </ImportSnippet>

func NewGraphClientWithChaosHandler(credential azcore.TokenCredential, scopes []string) *graph.GraphServiceClient {
	// <ChaosHandlerSnippet>
	// tokenCredential is one of the credential classes from azidentity
	// scopes is an array of permission scope strings
	authProvider, _ := authentication.NewAzureIdentityAuthenticationProviderWithScopes(credential, scopes)

	// Get default middleware from SDK
	defaultClientOptions := graph.GetDefaultClientOptions()
	defaultMiddleWare := graphcore.GetDefaultMiddlewaresWithOptions(&defaultClientOptions)

	// Add chaos handler to default middleware
	chaosHandler := khttp.NewChaosHandler()
	allMiddleware := append(defaultMiddleWare, chaosHandler)

	// Create an HTTP client with the middleware
	httpClient := khttp.GetDefaultClient(allMiddleware...)

	// Create the adapter
	// Passing nil values causes the adapter to use default implementations
	adapter, _ :=
		graph.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
			authProvider, nil, nil, httpClient)

	graphClient := graph.NewGraphServiceClient(adapter)
	// </ChaosHandlerSnippet>

	return graphClient
}

func NewGraphClientWithProxy(scopes []string) *graph.GraphServiceClient {
	// <ProxySnippet>
	proxyAddress := "http://proxy-url"
	proxyUrl, _ := url.Parse(proxyAddress)

	// Setup proxy for the token credential from azidentity
	authClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	credential, _ := azidentity.NewClientSecretCredential(
		"CLIENT_ID", "TENANT_ID", "SECRET", &azidentity.ClientSecretCredentialOptions{
			// Pass the proxied client to the credential
			ClientOptions: policy.ClientOptions{
				Transport: authClient,
			},
		})

	// scopes is an array of permission scope strings
	authProvider, _ := authentication.NewAzureIdentityAuthenticationProviderWithScopes(credential, scopes)

	// Get default middleware from SDK
	defaultClientOptions := graph.GetDefaultClientOptions()
	defaultMiddleWare := graphcore.GetDefaultMiddlewaresWithOptions(&defaultClientOptions)

	// Create an HTTP client with the middleware
	httpClient, _ := khttp.GetClientWithProxySettings(proxyAddress, defaultMiddleWare...)
	// For authenticated proxy, use
	// khttp.GetClientWithAuthenticatedProxySettings(
	//     proxyAddress, "user", "password", defaultMiddleWare...)

	// Create the adapter
	// Passing nil values causes the adapter to use default implementations
	adapter, _ :=
		graph.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
			authProvider, nil, nil, httpClient)

	graphClient := graph.NewGraphServiceClient(adapter)
	// </ProxySnippet>

	return graphClient
}
