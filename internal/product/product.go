package product

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type Float64TwoDecimal float64

func (f Float64TwoDecimal) MarshalJSON() ([]byte, error) {
	// Format the float to two decimal places
	return json.Marshal(fmt.Sprintf("%.2f", f))
}

type InputOrder struct {
	No                int               `json:"no"`
	PlatformProductId string            `json:"platformProductId"`
	Qty               int               `json:"qty"`
	UnitPrice         Float64TwoDecimal `json:"unitPrice"`
	TotalPrice        Float64TwoDecimal `json:"totalPrice"`
}

type CleanedOrder struct {
	No         int               `json:"no"`
	ProductId  string            `json:"productId"`
	MaterialId *string           `json:"materialId,omitempty"`
	ModelId    *string           `json:"modelId,omitempty"`
	Qty        int               `json:"qty"`
	UnitPrice  Float64TwoDecimal `json:"unitPrice"`
	TotalPrice Float64TwoDecimal `json:"totalPrice"`
}

// Function to parse productId and return materialToModel items
func ParseProductId(productId string, qty int, unitPrice Float64TwoDecimal, totalQty int) [][5]string {
	var materialToModelItems [][5]string
	var recursive_result [][5]string
	stringId := []rune(productId)

	foundFirstCharId := false
	isModelId := false
	isMaterialId := false
	newMaterialOnly := []rune{}
	newMaterialId := []rune{}
	newModelId := []rune{}
	newQuantity := fmt.Sprintf("%d", qty)
	foundQuantity := false
	counter := 0
	for i := range utf8.RuneCountInString(string(stringId)) {
		if stringId[i] == '/' {

			//recursive?
			qtyInt, err := strconv.Atoi(newQuantity)
			if err != nil {
				qtyInt = 0
			}
			// println("qtyInt is", qtyInt)
			recursive_result = ParseProductId(string(stringId[i+1:]), qty, unitPrice, qtyInt+qty)

			break
		}
		// Check for prefix
		// println("counter :", counter)
		// println("qty :", foundQuantity)
		if !foundFirstCharId {
			if stringId[i] == 'F' {
				foundFirstCharId = true
			} else {
				continue
			}
		}

		// Check for '-'
		if stringId[i] == '-' {
			counter++
		}

		// If found '-' 2 times, we get to the model Id
		if counter == 1 {
			isMaterialId = true
		}
		if counter == 2 {
			isModelId = true
			isMaterialId = false
		}
		if isMaterialId {
			newMaterialOnly = append(newMaterialOnly, stringId[i])
		}
		if !isModelId {
			newMaterialId = append(newMaterialId, stringId[i])
		}
		if isModelId {
			if foundQuantity {

				// println("new quantity is", stringId[i])
				newQuantity = string(stringId[i]) // This may need adjustment
				qtyInt, err := strconv.Atoi(newQuantity)
				if err == nil {
					totalQty += qtyInt - 1
				}
				// println("new found quantity is", newQuantity)
			}
			// Check suffix of quantity in model Id
			if stringId[i] != '*' && !foundQuantity {

				newModelId = append(newModelId, stringId[i])
			} else { // Found quantity
				// println("Found quantity")
				foundQuantity = true
				continue
			}

		}

	}
	// println("new material id is", string(newMaterialId))
	// println("new model id is", string(newModelId))
	// println("new quantity is", newQuantity)
	// println("new material only is", string(newMaterialOnly))
	// println("recursive count is", totalQty)
	// Create a new item to return
	materialToModelItems = append(materialToModelItems, [5]string{
		string(newMaterialId),  // Convert rune slice to string
		string(newModelId[1:]), // Convert rune slice to string
		string(newQuantity),    // Convert rune slice to string
		string(newMaterialOnly),
		string(strconv.Itoa(totalQty)),
	})

	materialToModelItems = append(materialToModelItems, recursive_result...)
	return materialToModelItems
}

func TransformOrders(input []InputOrder) []CleanedOrder {
	totalQty := 0
	var cleanedOrders []CleanedOrder
	productNo := 1
	materialCount := make(map[string]int)

	for _, order := range input {
		productId := order.PlatformProductId
		productQty := order.Qty
		productUnitPrice := order.UnitPrice

		// Call the parsing function
		materialToModelItems := ParseProductId(productId, productQty, productUnitPrice, productQty)
		recursiveRange := 0
		recursiveCount := 0
		// println(len(materialToModelItems))
		if len(materialToModelItems) > productNo {
			recursiveRange = len(materialToModelItems) - productNo
			recursiveCount = recursiveRange
		}
		// Count the number of products derived from the split

		// Calculate the new unit price based on the number of products
		newUnitPrice := productUnitPrice
		// println("total quantity by last index", materialToModelItems[productNo+recursiveRange-2][4])
		// println("productNo is", productNo)
		// println("this num ", productNo+recursiveRange-2)

		// Process the materialToModelItems as needed
		for _, item := range materialToModelItems {
			newQty, _ := strconv.Atoi(item[2])
			// println("productNo before is", productNo)
			// println("new quantity is", newQty)
			newMaterialId := item[0]
			newModelId := item[1]
			if recursiveCount > 0 {
				totalNumber, err := strconv.Atoi(materialToModelItems[productNo+recursiveRange-1][4])
				if err == nil {

				}
				newUnitPrice = order.UnitPrice / Float64TwoDecimal(totalNumber)
				recursiveCount--
			}
			cleanedOrder := CleanedOrder{
				No:         productNo,
				ProductId:  item[0] + "-" + item[1],
				MaterialId: &newMaterialId,
				ModelId:    &newModelId,
				Qty:        newQty,
				UnitPrice:  newUnitPrice,                             // Use the adjusted unit price
				TotalPrice: newUnitPrice * Float64TwoDecimal(newQty), // Calculate total price
			}
			materialCount[item[3][1:]] += newQty
			cleanedOrders = append(cleanedOrders, cleanedOrder)
			productNo++
			totalQty += newQty

		}
	}

	var nilString *string
	wipingCloth := CleanedOrder{
		No:         productNo,
		ProductId:  "WIPING-CLOTH",
		MaterialId: nilString,
		ModelId:    nilString,
		Qty:        totalQty,
		UnitPrice:  0,
		TotalPrice: 0,
	}
	productNo++
	cleanedOrders = append(cleanedOrders, wipingCloth)

	// Material order
	for key, value := range materialCount {
		cleanedOrder := CleanedOrder{
			No:         productNo,
			ProductId:  key + "-CLEANNER",
			MaterialId: nilString,
			ModelId:    nilString,
			Qty:        value,
			UnitPrice:  0,
			TotalPrice: 0,
		}
		cleanedOrders = append(cleanedOrders, cleanedOrder)
		productNo++
	}

	return cleanedOrders
}
