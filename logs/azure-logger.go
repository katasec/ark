package logs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// AzureWriter is an implementation of io.Writer. The logs.NewAzureWriter(..) func returns
// an AzureWriter that can be used with standard cmds like fmt.Fprintf() to write logs
type AzureWriter struct {
	ctx context.Context
	w   *appendblob.Client
}

func NewAzureWriter(accountName string, containerName string, accountKey string, fileName string) *AzureWriter {

	// Setup context
	ctx := context.Background()

	// Generate containerl url
	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	// Create credential
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Could not create shared key for account: %s\n", accountName)
		log.Println(err.Error())
	}
	// Create an container client with creds and url
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)
	if err != nil {
		log.Printf("Could not create container client key for url: %s\n", containerURL)
		log.Println(err.Error())
	}
	// Create appendClient from container client
	appendClient := containerClient.NewAppendBlobClient(fileName)

	// Create file if not exists
	_, err = appendClient.Create(ctx, nil)
	if err != nil {
		log.Println("Error creating 0-size append blob")
		log.Fatal(err.Error())
	}

	return &AzureWriter{
		w:   appendClient,
		ctx: ctx,
	}

}
func (e AzureWriter) Write(data []byte) (int, error) {

	tmStamp := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	message := fmt.Sprintf("%s  %s", tmStamp, string(data))

	body := streaming.NopCloser(strings.NewReader(message))

	_, err := e.w.AppendBlock(e.ctx, body, nil)
	if err != nil {
		log.Println("Error appending blob:" + err.Error())
		return 0, handleError(err)
	}

	return len([]byte(message)), nil
}

func handleError(err error) error {
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
