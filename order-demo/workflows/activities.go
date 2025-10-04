package workflows

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// ValidateOrderActivity validates the order
func ValidateOrderActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Validating order %s...\n", order.OrderID)
	time.Sleep(1 * time.Second)

	if order.Quantity <= 0 {
		return "", fmt.Errorf("invalid quantity: %d", order.Quantity)
	}
	if order.TotalAmount <= 0 {
		return "", fmt.Errorf("invalid amount: %.2f", order.TotalAmount)
	}

	return fmt.Sprintf("Order %s is valid", order.OrderID), nil
}

// ReserveInventoryActivity reserves inventory for the order
func ReserveInventoryActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Reserving %d units of %s for order %s...\n",
		order.Quantity, order.ProductName, order.OrderID)
	time.Sleep(2 * time.Second)

	// Simulate random failure (10% chance)
	if rand.Float32() < 0.1 {
		return "", fmt.Errorf("insufficient inventory for %s", order.ProductName)
	}

	return fmt.Sprintf("Reserved %d units of %s", order.Quantity, order.ProductName), nil
}

// ReleaseInventoryActivity releases reserved inventory (compensating action)
func ReleaseInventoryActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Releasing %d units of %s for order %s...\n",
		order.Quantity, order.ProductName, order.OrderID)
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("Released %d units of %s", order.Quantity, order.ProductName), nil
}

// ProcessPaymentActivity processes payment for the order
func ProcessPaymentActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Processing payment of $%.2f for order %s...\n",
		order.TotalAmount, order.OrderID)
	time.Sleep(2 * time.Second)

	// Simulate random failure (5% chance)
	if rand.Float32() < 0.05 {
		return "", fmt.Errorf("payment gateway timeout")
	}

	return fmt.Sprintf("Payment of $%.2f processed successfully", order.TotalAmount), nil
}

// ShipOrderActivity ships the order
func ShipOrderActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Shipping order %s to customer %s...\n",
		order.OrderID, order.CustomerID)
	time.Sleep(3 * time.Second)

	trackingNumber := fmt.Sprintf("TRACK-%s-%d", order.OrderID, time.Now().Unix())
	return fmt.Sprintf("Order shipped with tracking number: %s", trackingNumber), nil
}

// SendNotificationActivity sends notification to the customer
func SendNotificationActivity(ctx context.Context, order Order) (string, error) {
	fmt.Printf("[Activity] Sending notification to customer %s for order %s...\n",
		order.CustomerID, order.OrderID)
	time.Sleep(1 * time.Second)

	return fmt.Sprintf("Notification sent to customer %s", order.CustomerID), nil
}
