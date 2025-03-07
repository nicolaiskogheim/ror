package msswitchboardrabbitmqdefinitions

import (
	"fmt"
	"github.com/NorskHelsenett/ror/cmd/switchboard/switchboardconnections"

	"github.com/NorskHelsenett/ror/pkg/messagebuscontracts"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	QueueName          = "ms-switchboard"
	MsSwitchBoardQueue amqp.Queue
)

func InitOrDie() {
	queueArgs := amqp.Table{
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		"kind":            "VulnerabilityEvent",
	}
	var err error
	MsSwitchBoardQueue, err = switchboardconnections.RabbitMQConnection.GetChannel().QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		queueArgs, // arguments
	)
	if err != nil {
		args := [...]any{QueueName, err}
		msg := fmt.Sprintf("could not declare queue %s,", args)
		rlog.Fatal(msg, err)
	}

	routingKey := messagebuscontracts.Workqueue_Switchboard_Post
	err = switchboardconnections.RabbitMQConnection.GetChannel().QueueBind(
		QueueName,                                // queue name
		routingKey,                               // routing key
		messagebuscontracts.ExchangeRorResources, // exchange
		false,
		amqp.Table{
			"kind": "VulnerabilityEvent",
		},
	)
	if err != nil {
		args := [...]any{QueueName, routingKey, err}
		msg := fmt.Sprintf("could not bind queue  %s,", args)
		rlog.Fatal(msg, err)
	}
}
