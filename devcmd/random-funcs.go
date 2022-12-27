package devcmd

import (
	"fmt"
	"log"

	"github.com/katasec/ark/logs"
	"github.com/katasec/ark/messaging"
)

func startDb() {

	// Start Postgres Container
	imageName := DEV_PGSQL_IMAGE_NAME
	envVars := []string{
		"POSTGRES_USER=" + DevDbDefaultUser,
		"POSTGRES_PASSWORD=" + DevDbDefaultPass,
	}
	port := "5432"

	dh.StartContainerUI(imageName, envVars, port, "arkdb", nil)
}

func testMq() {

	connectionString := d.Config.AzureConfig.MqConfig.MqConnectionString
	queueName := d.Config.AzureConfig.MqConfig.MqName

	var mq messaging.Messenger = messaging.NewAsbMessenger(connectionString, queueName)

	smessage := "hello world"
	log.Println("Sending message:" + smessage)
	err := mq.Send(smessage)
	if err != nil {
		log.Fatalf(err.Error())
	} else {
		log.Println("Message sent successfully")
	}

	log.Println("Receiving message")
	rmessage, err := mq.Receive()
	if err != nil {
		log.Fatalf(err.Error())
	} else {
		log.Println("The received message was:" + rmessage)
	}

}

func testLogging() {
	accountName := d.Config.AzureConfig.StorageConfig.LogStorageAccountName
	containerName := d.Config.AzureConfig.StorageConfig.LogsContainer
	accountKey := d.Config.AzureConfig.StorageConfig.LogStorageKey
	w := logs.NewAzureWriter(accountName, containerName, accountKey, "testfile.txt")

	_, err := fmt.Fprintln(w, "Hi this is a test")
	if err != nil {
		log.Println("error:" + err.Error())
	}
}
