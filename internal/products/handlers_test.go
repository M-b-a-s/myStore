package products

import "testing"

func TestValidateProductInput(t *testing.T) {
	tests := []struct {
		name         string
		productName  string
		priceInCents int32
		quantity     int32
		wantErr      bool
	}{
		{name: "valid input", productName: "Keyboard", priceInCents: 2999, quantity: 10},
		{name: "empty name", productName: "", priceInCents: 2999, quantity: 10, wantErr: true},
		{name: "whitespace name", productName: "  ", priceInCents: 2999, quantity: 10, wantErr: true},
		{name: "negative price", productName: "Keyboard", priceInCents: -1, quantity: 10, wantErr: true},
		{name: "negative quantity", productName: "Keyboard", priceInCents: 2999, quantity: -1, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateProductInput(test.productName, test.priceInCents, test.quantity)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateProductInput() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
