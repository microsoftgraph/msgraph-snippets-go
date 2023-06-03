// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

// <ImportSnippet>
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	auth "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// </ImportSnippet>

func NewGraphClientForUsGov() *graph.GraphServiceClient {
	// <NationalCloudSnippet>
	// Create the InteractiveBrowserCredential using details
	// from app registered in the Azure AD for US Government portal
	credential, _ := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		ClientID: "YOUR_CLIENT_ID",
		TenantID: "YOUR_TENANT_ID",
		ClientOptions: policy.ClientOptions{
			// https://login.microsoftonline.us
			Cloud: cloud.AzureGovernment,
		},
		RedirectURL: "YOUR_REDIRECT_URL",
	})

	// Create the authentication provider
	authProvider, _ := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential,
		[]string{"https://graph.microsoft.us/.default"})

	// Create a request adapter using the auth provider
	adapter, _ := graph.NewGraphRequestAdapter(authProvider)

	// Set the service root to the
	// Microsoft Graph for US Government L4 endpoint
	// NOTE: The API version must be included in the URL
	adapter.SetBaseUrl("https://graph.microsoft.us/v1.0")

	// Create a Graph client using request adapter
	graphClient := graph.NewGraphServiceClient(adapter)
	// </NationalCloudSnippet>

	return graphClient
}
