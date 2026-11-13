package checkout

import "testing"

func TestInventoryAvailable(t *testing.T) {
	tests := []struct {
		name              string
		total, sold, held int
		quantity          int
		want              bool
	}{
		{name: "last ticket available", total: 1, quantity: 1, want: true},
		{name: "last ticket already held", total: 1, held: 1, quantity: 1, want: false},
		{name: "sold and held fill quota", total: 10, sold: 6, held: 4, quantity: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inventoryAvailable(tt.total, tt.sold, tt.held, tt.quantity); got != tt.want {
				t.Fatalf("inventoryAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
