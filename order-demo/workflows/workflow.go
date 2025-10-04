package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Order represents an order
type Order struct {
	OrderID       string
	CustomerID    string
	ProductName   string
	Quantity      int
	TotalAmount   float64
	Status        string
}

// OrderWorkflow is the workflow definition
func OrderWorkflow(ctx workflow.Context, order Order) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Order workflow started", "OrderID", order.OrderID)

	// Set workflow options for activities
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Validate order
	logger.Info("Step 1: Validating order")
	var validationResult string
	err := workflow.ExecuteActivity(ctx, ValidateOrderActivity, order).Get(ctx, &validationResult)
	if err != nil {
		logger.Error("Order validation failed", "error", err)
		return "", err
	}
	logger.Info("Order validated successfully", "result", validationResult)

	// Step 2: Reserve inventory
	logger.Info("Step 2: Reserving inventory")
	var inventoryResult string
	err = workflow.ExecuteActivity(ctx, ReserveInventoryActivity, order).Get(ctx, &inventoryResult)
	if err != nil {
		logger.Error("Inventory reservation failed", "error", err)
		return "", err
	}
	logger.Info("Inventory reserved", "result", inventoryResult)

	// Step 3: Process payment
	logger.Info("Step 3: Processing payment")
	var paymentResult string
	err = workflow.ExecuteActivity(ctx, ProcessPaymentActivity, order).Get(ctx, &paymentResult)
	if err != nil {
		logger.Error("Payment failed, rolling back inventory", "error", err)
		// Compensating transaction: release inventory
		_ = workflow.ExecuteActivity(ctx, ReleaseInventoryActivity, order).Get(ctx, nil)
		return "", err
	}
	logger.Info("Payment processed", "result", paymentResult)

	// Step 4: Ship order
	logger.Info("Step 4: Shipping order")
	var shipmentResult string
	err = workflow.ExecuteActivity(ctx, ShipOrderActivity, order).Get(ctx, &shipmentResult)
	if err != nil {
		logger.Error("Shipment failed", "error", err)
		return "", err
	}
	logger.Info("Order shipped", "result", shipmentResult)

	// Step 5: Send notification
	logger.Info("Step 5: Sending notification")
	var notificationResult string
	err = workflow.ExecuteActivity(ctx, SendNotificationActivity, order).Get(ctx, &notificationResult)
	if err != nil {
		logger.Warn("Notification failed, but order is complete", "error", err)
	}

	finalStatus := fmt.Sprintf("Order %s completed successfully!", order.OrderID)
	logger.Info("Order workflow completed", "OrderID", order.OrderID)
	return finalStatus, nil
}
