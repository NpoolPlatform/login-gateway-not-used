package client

import (
	"fmt"

	msgcli "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/client"
	constant "github.com/NpoolPlatform/login-gateway/pkg/message/const"
	msg "github.com/NpoolPlatform/login-gateway/pkg/message/message"

	"github.com/streadway/amqp"
)

type client struct {
	*msgcli.Client
	consumers map[string]<-chan amqp.Delivery
}

var myClients = map[string]*client{}

func Init() error {
	_myClient, err := msgcli.New(constant.ServiceName)
	if err != nil {
		return err
	}

	err = _myClient.DeclareQueue(msg.QueueExample)
	if err != nil {
		return err
	}

	sampleClient := &client{
		Client:    _myClient,
		consumers: map[string]<-chan amqp.Delivery{},
	}
	examples, err := _myClient.Consume(msg.QueueExample)
	if err != nil {
		return fmt.Errorf("fail to construct example consume: %v", err)
	}
	sampleClient.consumers[msg.QueueExample] = examples

	myClients[constant.ServiceName] = sampleClient

	return nil
}
