package data

import "testing"

func TestValidateProduct(t *testing.T) {
	prod := Product{
		Name:  "Tea",
		Price: 1.23,
		SKU:   "abc-def-ghi",
	}
	err := prod.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
