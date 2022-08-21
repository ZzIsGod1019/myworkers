package test

import (
	"log"
	"myworkers/config"
	"myworkers/sqs"
	"testing"
)

func TestXxx(t *testing.T) {
	config.Init()
	message, err := sqs.ReceiveMessage()
	if err != nil {
		log.Println("error:" + err.Error())
	} else {
		log.Println("info:" + *message.Body)
	}
}
