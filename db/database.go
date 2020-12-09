package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gocarina/gocsv"
)

type InventoryDB interface {
}

type CSVDatabase struct {
	Inventory       []*InventoryItem
	NextInventoryID int
	InventoryLock   sync.Mutex
}

// NewCSVDatabase sets up new mock db from csv strings
func NewCSVDatabase(ctx context.Context, inventoryCSV string) (*CSVDatabase, error) {
	if inventoryCSV == "" {
		return nil, errors.New("Must provide inventory string for csv database")
	}

	db := &CSVDatabase{}

	// Setup inventory
	var invItems []*InventoryItem
	err := gocsv.UnmarshalString(inventoryCSV, &invItems)
	if err != nil {
		e := fmt.Sprintf("Error reading from csv db - %s", err.Error())
		lg(ctx).Error(e)
		return nil, errors.New(e)
	}

	if len(invItems) != 0 {
		db.NextInventoryID = invItems[len(invItems)-1].ID + 1
	}
	db.Inventory = invItems

	return db, nil
}
