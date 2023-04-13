package events

import (
	"fmt"
	"solace-events-consumer/src/business"
	"time"

	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/message"
	"solace.dev/go/messaging/pkg/solace/resource"
)

// Configuration parameters
var (
	brokerConfig = config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                "tcp://192.168.0.100:55555",
		config.ServicePropertyVPNName:                    "zecomeia",
		config.AuthenticationPropertySchemeBasicUserName: "catatau",
		config.AuthenticationPropertySchemeBasicPassword: "catatau",
	}
	messagingService   solace.MessagingService
	persistentReceiver = make(map[int]solace.PersistentMessageReceiver)
)

type messageHandler struct {
	idReceiver int
}

func (h messageHandler) ProcessMessage(message message.InboundMessage) {
	fmt.Printf("Message processed by receiver %d \n", h.idReceiver)
	var messageBody string

	if payload, ok := message.GetPayloadAsString(); ok {
		messageBody = payload
	} else if payload, ok := message.GetPayloadAsBytes(); ok {
		messageBody = string(payload)
	}

	business.ProcessBookMessageEvent(messageBody)
	// fmt.Printf("Message Dump %s \n", message)
}

func ConnectReceiver(numberOfConns int) {

	if mS, err := messaging.NewMessagingServiceBuilder().FromConfigurationProvider(brokerConfig).Build(); err != nil {
		panic(err)
	} else {
		messagingService = mS
	}

	// Connect to the messaging serice
	if err := messagingService.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the broker? ", messagingService.IsConnected())

	queueName := "books"

	for i := 1; i <= numberOfConns; i++ {
		// Build a Gauranteed message receiver and bind to the given queue
		durableExclusiveQueue := resource.QueueDurableExclusive(queueName)
		if pR, err := messagingService.CreatePersistentMessageReceiverBuilder().WithMessageAutoAcknowledgement().Build(durableExclusiveQueue); err != nil {
			panic(err)
		} else {
			persistentReceiver[i] = pR
		}
		/*
			nonDurableExclusiveQueue := resource.QueueNonDurableExclusive(queueName)
			topicString := "cleovon/test/book/create"
			topic := resource.TopicSubscriptionOf(topicString)
			strategy := config.MissingResourcesCreationStrategy("CREATE_ON_START")
			if pR, err := messagingService.CreatePersistentMessageReceiverBuilder().WithMissingResourcesCreationStrategy(strategy).WithSubscriptions(topic).Build(nonDurableExclusiveQueue); err != nil {
				panic(err)
			} else {
				persistentReceiver = pR
			}
		*/
		// Handling a panic from a non existing queue
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Make sure queue name '%s' exists on the broker.\nThe following error occurred when attempting to connect to create a Persistent Message Receiver:\n%s", queueName, err)
			}
		}()

		// Start Persistent Message Receiver
		if err := persistentReceiver[i].Start(); err != nil {
			panic(err)
		}

		fmt.Printf("Persistent Receiver %d running? %t\n", i, persistentReceiver[i].IsRunning())

		// Register Message callback handler to the Message Receiver
		if regErr := persistentReceiver[i].ReceiveAsync(messageHandler{idReceiver: i}.ProcessMessage); regErr != nil {
			panic(regErr)
		}
		fmt.Printf("Bound to queue: %s\n", queueName)
	}
	//fmt.Println("\n===Interrupt (CTR+C) to handle graceful terminaltion of the subscriber===")
}

func CloseReceiver() {
	// Terminate the Persistent Receiver
	for i, pR := range persistentReceiver {
		pR.Terminate(1 * time.Second)
		fmt.Printf("\nPersistent Receiver %d Terminated? %t\n", i, pR.IsTerminated())
	}

	// Disconnect the Message Service
	messagingService.Disconnect()
	fmt.Printf("Messaging Service Disconnected? %t\n", !messagingService.IsConnected())
}
