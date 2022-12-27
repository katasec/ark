package devcmd

import (
	"log"

	"github.com/katasec/ark/messaging"
)

func Start() {
	//d.RefreshConfig()

	testMq()
}

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
