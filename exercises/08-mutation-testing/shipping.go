package testing

// EligibleForFreeShipping informa se o total do pedido alcançou o mínimo.
func EligibleForFreeShipping(orderTotal, minimum int) bool {
	return orderTotal >= minimum
}
