// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package snippets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/fileuploader"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

func RunUploadSamples(graphClient *graph.GraphServiceClient, largeFile string) {
	itemPath := "Documents/vacation.gif"

	UploadFileToOneDrive(graphClient, largeFile, itemPath)
	UploadAttachmentToMessage(graphClient, largeFile)
}

func UploadFileToOneDrive(graphClient *graph.GraphServiceClient, largeFile string, itemPath string) {
	// <LargeFileUploadSnippet>
	byteStream, _ := os.Open(largeFile)

	// Use properties to specify the conflict behavior
	itemUploadProperties := models.NewDriveItemUploadableProperties()
	itemUploadProperties.SetAdditionalData(map[string]any{"@microsoft.graph.conflictBehavior": "replace"})
	uploadSessionRequestBody := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	uploadSessionRequestBody.SetItem(itemUploadProperties)

	// Create the upload session
	// itemPath does not need to be a path to an existing item
	myDrive, _ := graphClient.Me().Drive().Get(context.Background(), nil)

	uploadSession, _ := graphClient.Drives().
		ByDriveId(*myDrive.GetId()).
		Items().
		ByDriveItemId("root:/"+itemPath+":").
		CreateUploadSession().
		Post(context.Background(), uploadSessionRequestBody, nil)

	// Max slice size must be a multiple of 320 KiB
	maxSliceSize := int64(320 * 1024)
	fileUploadTask := fileuploader.NewLargeFileUploadTask[models.DriveItemable](
		graphClient.RequestAdapter,
		uploadSession,
		byteStream,
		maxSliceSize,
		models.CreateDriveItemFromDiscriminatorValue,
		nil)

	// Create a callback that is invoked after each slice is uploaded
	progress := func(progress int64, total int64) {
		fmt.Printf("Uploaded %d of %d bytes\n", progress, total)
	}

	// Upload the file
	uploadResult := fileUploadTask.Upload(progress)

	if uploadResult.GetUploadSucceeded() {
		fmt.Printf("Upload complete, item ID: %s\n", *uploadResult.GetItemResponse().GetId())
	} else {
		fmt.Print("Upload failed.\n")
	}
	// </LargeFileUploadSnippet>
}

func ResumeUpload(
	fileUploadTask fileuploader.LargeFileUploadTask[models.DriveItemable],
	progress fileuploader.ProgressCallBack) {
	// <ResumeSnippet>
	fileUploadTask.Resume(progress)
	// </ResumeSnippet>
}

func UploadAttachmentToMessage(graphClient *graph.GraphServiceClient, largeFile string) {
	// <UploadAttachmentSnippet>
	// Create message
	message := models.NewMessage()
	subject := "Large attachment"
	message.SetSubject(&subject)

	savedDraft, _ := graphClient.Me().Messages().Post(context.Background(), message, nil)

	// Set up the attachment
	byteStream, _ := os.Open(largeFile)
	largeAttachment := models.NewAttachmentItem()
	attachmentType := models.FILE_ATTACHMENTTYPE
	largeAttachment.SetAttachmentType(&attachmentType)
	fileName := filepath.Base(largeFile)
	largeAttachment.SetName(&fileName)
	fileInfo, _ := byteStream.Stat()
	fileSize := fileInfo.Size()
	largeAttachment.SetSize(&fileSize)

	uploadSessionRequestBody := users.NewItemMessagesItemAttachmentsCreateUploadSessionPostRequestBody()
	uploadSessionRequestBody.SetAttachmentItem(largeAttachment)

	uploadSession, _ := graphClient.Me().
		Messages().
		ByMessageId(*savedDraft.GetId()).
		Attachments().
		CreateUploadSession().
		Post(context.Background(), uploadSessionRequestBody, nil)

	// Max slice size must be a multiple of 320 KiB
	maxSliceSize := int64(320 * 1024)
	fileUploadTask := fileuploader.NewLargeFileUploadTask[models.FileAttachmentable](
		graphClient.RequestAdapter,
		uploadSession,
		byteStream,
		maxSliceSize,
		models.CreateFileAttachmentFromDiscriminatorValue,
		nil)

	// Create a callback that is invoked after each slice is uploaded
	progress := func(progress int64, total int64) {
		fmt.Printf("Uploaded %d of %d bytes\n", progress, total)
	}

	// Upload the file
	uploadResult := fileUploadTask.Upload(progress)

	if uploadResult.GetUploadSucceeded() {
		fmt.Print("Upload complete\n")
	} else {
		fmt.Print("Upload failed.\n")
	}
	// </UploadAttachmentSnippet>
}
