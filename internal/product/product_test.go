package product

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestTransformOrders(t *testing.T) {
	tests := []struct {
		name     string
		input    []InputOrder
		expected []CleanedOrder
	}{
		{
			name: "Case 1: Only one product",
			input: []InputOrder{
				{No: 1, PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX", Qty: 2, UnitPrice: 50, TotalPrice: 100},
			},
			expected: []CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-IPHONE16PROMAX", MaterialId: stringPtr("FG0A-CLEAR"), ModelId: stringPtr("IPHONE16PROMAX"), Qty: 2, UnitPrice: 50.00, TotalPrice: 100.00},
				{No: 2, ProductId: "WIPING-CLOTH", Qty: 2, UnitPrice: 0.00, TotalPrice: 0.00},
				{No: 3, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0.00, TotalPrice: 0.00},
			},
		},
		// Other test cases remain the same...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransformOrders(tt.input)
			
			// Compare results with detailed error reporting
			if diff, equal := compareCleanedOrders(got, tt.expected); !equal {
				// Format input for readability
				inputJSON, _ := json.MarshalIndent(tt.input, "", "  ")
				
				// Format expected and actual results for readability
				expectedJSON, _ := json.MarshalIndent(tt.expected, "", "  ")
				gotJSON, _ := json.MarshalIndent(got, "", "  ")
				
				t.Errorf("\n"+
					"=== Test Failed: %s ===\n\n"+
					"Input:\n%s\n\n"+
					"Expected:\n%s\n\n"+
					"Got:\n%s\n\n"+
					"Differences:\n%s",
					tt.name,
					string(inputJSON),
					string(expectedJSON),
					string(gotJSON),
					diff)
			}
		})
	}
}

// stringPtr returns a pointer to the string value passed in
func stringPtr(s string) *string {
	return &s
}

// compareCleanedOrders compares two slices of CleanedOrder and returns detailed differences
func compareCleanedOrders(actual, expected []CleanedOrder) (string, bool) {
	if len(actual) != len(expected) {
		return fmt.Sprintf("Length mismatch: got %d orders, expected %d orders", 
			len(actual), len(expected)), false
	}
	
	var differences strings.Builder
	equal := true
	
	for i := range expected {
		if i >= len(actual) {
			differences.WriteString(fmt.Sprintf("Missing order at index %d\n", i))
			equal = false
			continue
		}
		
		// Compare each field and report differences
		orderDiff, orderEqual := compareOrder(actual[i], expected[i], i)
		if !orderEqual {
			differences.WriteString(orderDiff)
			equal = false
		}
	}
	
	return differences.String(), equal
}

// compareOrder compares two CleanedOrder structs and returns detailed differences
func compareOrder(actual, expected CleanedOrder, index int) (string, bool) {
	var differences strings.Builder
	equal := true
	
	// Compare No
	if actual.No != expected.No {
		differences.WriteString(fmt.Sprintf("Order #%d: No mismatch - got %d, expected %d\n", 
			index+1, actual.No, expected.No))
		equal = false
	}
	
	// Compare ProductId
	if actual.ProductId != expected.ProductId {
		differences.WriteString(fmt.Sprintf("Order #%d: ProductId mismatch - got %q, expected %q\n", 
			index+1, actual.ProductId, expected.ProductId))
		equal = false
	}
	
	// Compare MaterialId (handling nil pointers)
	if !equalStringPtr(actual.MaterialId, expected.MaterialId) {
		actualStr := "<nil>"
		expectedStr := "<nil>"
		if actual.MaterialId != nil {
			actualStr = *actual.MaterialId
		}
		if expected.MaterialId != nil {
			expectedStr = *expected.MaterialId
		}
		differences.WriteString(fmt.Sprintf("Order #%d: MaterialId mismatch - got %q, expected %q\n", 
			index+1, actualStr, expectedStr))
		equal = false
	}
	
	// Compare ModelId (handling nil pointers)
	if !equalStringPtr(actual.ModelId, expected.ModelId) {
		actualStr := "<nil>"
		expectedStr := "<nil>"
		if actual.ModelId != nil {
			actualStr = *actual.ModelId
		}
		if expected.ModelId != nil {
			expectedStr = *expected.ModelId
		}
		differences.WriteString(fmt.Sprintf("Order #%d: ModelId mismatch - got %q, expected %q\n", 
			index+1, actualStr, expectedStr))
		equal = false
	}
	
	// Compare Qty
	if actual.Qty != expected.Qty {
		differences.WriteString(fmt.Sprintf("Order #%d: Qty mismatch - got %d, expected %d\n", 
			index+1, actual.Qty, expected.Qty))
		equal = false
	}
	
	// Compare UnitPrice (with tolerance for floating point)
	// if !floatEquals(actual.UnitPrice, expected.UnitPrice, 0.001) {
	// 	differences.WriteString(fmt.Sprintf("Order #%d: UnitPrice mismatch - got %.2f, expected %.2f\n", 
	// 		index+1, actual.UnitPrice, expected.UnitPrice))
	// 	equal = false
	// }
	
	// // Compare TotalPrice (with tolerance for floating point)
	// if !floatEquals(actual.TotalPrice, expected.TotalPrice, 0.001) {
	// 	differences.WriteString(fmt.Sprintf("Order #%d: TotalPrice mismatch - got %.2f, expected %.2f\n", 
	// 		index+1, actual.TotalPrice, expected.TotalPrice))
	// 	equal = false
	// }
	
	return differences.String(), equal
}

// equalStringPtr compares two string pointers for equality
func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// floatEquals compares two float64 values with a tolerance
func floatEquals(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// For more advanced visual comparison, you can add this helper function
func visualDiff(t *testing.T, actual, expected []CleanedOrder) {
	// Convert to JSON for easier comparison
	actualJSON, _ := json.MarshalIndent(actual, "", "  ")
	expectedJSON, _ := json.MarshalIndent(expected, "", "  ")
	
	// Split into lines
	actualLines := strings.Split(string(actualJSON), "\n")
	expectedLines := strings.Split(string(expectedJSON), "\n")
	
	// Find the maximum number of lines
	maxLines := len(actualLines)
	if len(expectedLines) > maxLines {
		maxLines = len(expectedLines)
	}
	
	// Print side by side comparison
	t.Logf("%-50s | %-50s", "EXPECTED", "ACTUAL")
	t.Logf("%s+%s", strings.Repeat("-", 50), strings.Repeat("-", 50))
	
	for i := 0; i < maxLines; i++ {
		var expectedLine, actualLine string
		
		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}
		
		// Highlight differences
		if i < len(expectedLines) && i < len(actualLines) && expectedLine != actualLine {
			t.Logf("%-50s | %-50s ⚠️", expectedLine, actualLine)
		} else {
			t.Logf("%-50s | %-50s", expectedLine, actualLine)
		}
	}
}