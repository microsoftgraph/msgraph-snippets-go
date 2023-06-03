// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

func NewGraphClientWithClientSecret() *graph.GraphServiceClient {
	// <ClientSecretSnippet>
	cred, _ := azidentity.NewClientSecretCredential(
		"TENANT_ID",
		"CLIENT_ID",
		"CLIENT_SECRET",
		nil,
	)

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"https://graph.microsoft.com/.default"})
	// </ClientSecretSnippet>

	return graphClient
}

func NewGraphClientWithClientCertificate() *graph.GraphServiceClient {
	// <ClientCertificateSnippet>
	// Load certificate
	certFile, _ := os.Open("certificate.pem")
	info, _ := certFile.Stat()
	certBytes := make([]byte, info.Size())
	certFile.Read(certBytes)
	certFile.Close()

	certs, key, _ := azidentity.ParseCertificates(certBytes, nil)

	cred, _ := azidentity.NewClientCertificateCredential(
		"TENANT_ID",
		"CLIENT_ID",
		certs,
		key,
		nil,
	)

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"https://graph.microsoft.com/.default"})
	// </ClientCertificateSnippet>

	return graphClient
}

func NewGraphClientWithOnBehalfOf() *graph.GraphServiceClient {
	// <OnBehalfOfSnippet>
	cred, _ := azidentity.NewOnBehalfOfCredentialWithSecret(
		"TENANT_ID",
		"CLIENT_ID",
		"USER_ASSERTION_STRING",
		"CLIENT_SECRET",
		nil,
	)

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"https://graph.microsoft.com/.default"})
	// </OnBehalfOfSnippet>

	return graphClient
}

func NewGraphClientWithDeviceCode() *graph.GraphServiceClient {
	// <DeviceCodeSnippet>
	cred, _ := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: "TENANT_ID",
		ClientID: "CLIENT_ID",
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"User.Read"})
	// </DeviceCodeSnippet>

	return graphClient
}

func NewGraphClientWithInteractive() *graph.GraphServiceClient {
	// <InteractiveSnippet>
	cred, _ := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID:    "TENANT_ID",
		ClientID:    "CLIENT_ID",
		RedirectURL: "REDIRECT_URL",
	})

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"User.Read"})
	// </InteractiveSnippet>

	return graphClient
}

func NewGraphClientWithUserNamePassword() *graph.GraphServiceClient {
	// <UserNamePasswordSnippet>
	cred, _ := azidentity.NewUsernamePasswordCredential(
		"TENANT_ID",
		"CLIENT_ID",
		"USER_NAME",
		"PASSWORD",
		nil,
	)

	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		cred, []string{"User.Read"})
	// </UserNamePasswordSnippet>

	return graphClient
}
