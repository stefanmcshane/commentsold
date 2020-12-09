package db

import (
	"context"
	"sync"
	"testing"
)

func Test_ListInventory(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    int
	}{
		{
			name:    "inventory - success",
			args:    args{context.Background()},
			wantErr: false,
			want:    50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewCSVDatabase(tt.args.ctx, inventoryData)
			got, err := db.GetInventoryItems(tt.args.ctx, 0, 50)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Expected to find %d items, found %d", tt.want, len(got))
			}
		})
	}
}
func TestCSVDatabase_UpdateInventoryItem(t *testing.T) {
	db, _ := NewCSVDatabase(context.Background(), inventoryData)
	type fields struct {
		Inventory       []*InventoryItem
		NextInventoryID int
		InventoryLock   *sync.Mutex
	}
	type args struct {
		ctx         context.Context
		updatedItem InventoryItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "inventory update - should fail",
			args:    args{ctx: context.Background(), updatedItem: InventoryItem{ID: 1, ProductID: "9866", Quantity: 0, Color: "RED"}},
			wantErr: true,
		},
		{
			name:    "inventory update - success",
			args:    args{ctx: context.Background(), updatedItem: InventoryItem{ID: 206, ProductID: "9866", Quantity: 0, Color: "RED"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.UpdateInventoryItem(tt.args.ctx, tt.args.updatedItem); (err != nil) != tt.wantErr {
				t.Errorf("CSVDatabase.UpdateInventoryItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
