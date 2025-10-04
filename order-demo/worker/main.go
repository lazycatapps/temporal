package main

import (
	"flag"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"order-demo/workflows"
)

func main() {
	// Parse command line arguments
	var temporalAddress string
	flag.StringVar(&temporalAddress, "temporal-address", "localhost:7233", "Temporal server address")
	flag.Parse()

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: temporalAddress,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, "order-queue", worker.Options{})

	// Register workflow and activities
	w.RegisterWorkflow(workflows.OrderWorkflow)
	w.RegisterActivity(workflows.ValidateOrderActivity)
	w.RegisterActivity(workflows.ReserveInventoryActivity)
	w.RegisterActivity(workflows.ReleaseInventoryActivity)
	w.RegisterActivity(workflows.ProcessPaymentActivity)
	w.RegisterActivity(workflows.ShipOrderActivity)
	w.RegisterActivity(workflows.SendNotificationActivity)

	// Start worker
	log.Printf("Starting Order Processing Worker (connected to %s)...\n", temporalAddress)
	log.Println("Worker is listening on task queue: order-queue")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
