package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

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

	// Create sample order
	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())
	order := workflows.Order{
		OrderID:     orderID,
		CustomerID:  "CUST-12345",
		ProductName: "Laptop",
		Quantity:    2,
		TotalAmount: 2999.98,
		Status:      "pending",
	}

	// Start workflow execution
	workflowOptions := client.StartWorkflowOptions{
		ID:        "order-workflow-" + orderID,
		TaskQueue: "order-queue",
	}

	fmt.Println("=================================================")
	fmt.Printf("Temporal Server: %s\n", temporalAddress)
	fmt.Printf("Starting Order Workflow for Order: %s\n", order.OrderID)
	fmt.Printf("Customer: %s\n", order.CustomerID)
	fmt.Printf("Product: %s (Quantity: %d)\n", order.ProductName, order.Quantity)
	fmt.Printf("Total Amount: $%.2f\n", order.TotalAmount)
	fmt.Println("=================================================")

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.OrderWorkflow, order)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	fmt.Printf("\nWorkflow started with ID: %s\n", we.GetID())
	fmt.Printf("RunID: %s\n", we.GetRunID())
	fmt.Printf("\nYou can view the workflow in Temporal UI:\n")
	fmt.Printf("http://localhost:8080/namespaces/default/workflows/%s\n\n", we.GetID())

	// Wait for workflow completion
	fmt.Println("Waiting for workflow to complete...")
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}

	fmt.Println("\n=================================================")
	fmt.Printf("Workflow Result: %s\n", result)
	fmt.Println("=================================================")
}
