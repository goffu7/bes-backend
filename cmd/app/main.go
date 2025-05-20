package main

import (
	"bes-backend/internal/product"
	"encoding/json"
	"fmt"
)

func main() {
	// Sample input for demonstration
	input := []product.InputOrder{
		{
			No:                1,
			PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
			Qty:               1,
			UnitPrice:         160,
			TotalPrice:        160,
		},
		{
			No:                2,
			PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
			Qty:               1,
			UnitPrice:         50,
			TotalPrice:        50,
		},
	}

	// Transform the orders
	cleanedOrders := product.TransformOrders(input)

	// Marshal the cleaned orders to JSON for output
	outputJSON, err := json.MarshalIndent(cleanedOrders, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Print the output
	fmt.Println(string(outputJSON))
}
