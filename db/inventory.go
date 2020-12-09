package db

import (
	"context"
	"errors"
	"fmt"
)

// InventoryItem represents the base information for an item stored in the inventory
type InventoryItem struct {
	ID             int    `csv:"id"`
	ProductID      string `csv:"product_id"`
	Quantity       int    `csv:"quantity"`
	Color          string `csv:"color"`
	Size           string `csv:"size"`
	PriceCents     int    `csv:"price_cents"`
	SalePriceCents int    `csv:"sale_price_cents"`
}

// GetInventoryItems returns the number of inventory items based on the start and end numbers for paginating
// Will succeed as long as end isnt greater that the total number of records
func (db *CSVDatabase) GetInventoryItems(ctx context.Context, start int, end int) ([]*InventoryItem, error) {
	if start < 0 {
		e := fmt.Sprintf("Start value too low, please specify a higher number than %d", start)
		lg(ctx).Error(e)
		return nil, errors.New(e)
	}
	if end >= len(db.Inventory) {
		e := fmt.Sprintf("End value too high, please specify a higher number than %d", len(db.Inventory))
		lg(ctx).Error(e)
		return nil, errors.New(e)
	}

	return db.Inventory[start:end], nil
}

// GetInventoryItemByID returns the number of inventory items based on the start and end numbers for paginating
// Will succeed as long as end isnt greater that the total number of records
func (db *CSVDatabase) GetInventoryItemByID(ctx context.Context, ID int) (*InventoryItem, error) {
	for _, i := range db.Inventory {
		if i.ID == ID {
			return i, nil
		}
	}

	e := fmt.Sprintf("Inventory item by ID %d doesnt exist", ID)
	lg(ctx).Error(e)
	return nil, errors.New(e)
}

// AddInventoryItem will create a new inventory item in the csv database
func (db *CSVDatabase) AddInventoryItem(ctx context.Context, newItem InventoryItem) error {
	db.InventoryLock.Lock()
	defer db.InventoryLock.Unlock()

	newItem.ID = db.NextInventoryID
	db.NextInventoryID++
	db.Inventory = append(db.Inventory, &newItem)
	return nil
}

// UpdateInventoryItem will find the item by the provided ID, and override it in the db
func (db *CSVDatabase) UpdateInventoryItem(ctx context.Context, updatedItem InventoryItem) error {
	db.InventoryLock.Lock()
	defer db.InventoryLock.Unlock()

	i, err := db.GetInventoryItemByID(ctx, updatedItem.ID)
	if err != nil {
		e := fmt.Sprintf("Unable to update inventory item - %s", err.Error())
		lg(ctx).Error(e)
		return errors.New(e)
	}

	*i = updatedItem

	return nil
}

// UpdateInventoryItemStock will find the item by the provided ID, and increment or decrement based on the specified adjustment amount
func (db *CSVDatabase) UpdateInventoryItemStock(ctx context.Context, updatedItemID int, adjustment int) error {
	db.InventoryLock.Lock()
	defer db.InventoryLock.Unlock()

	i, err := db.GetInventoryItemByID(ctx, updatedItemID)
	if err != nil {
		e := fmt.Sprintf("Unable to update inventory item - %s", err.Error())
		lg(ctx).Error(e)
		return errors.New(e)
	}

	newValue := i.Quantity + adjustment
	if newValue < 0 {
		e := fmt.Sprintf("Not enough inventory in stock for item %d", updatedItemID)
		lg(ctx).Error(e)
		return errors.New(e)
	}

	fmt.Println(newValue, i.Quantity)

	i.Quantity = newValue

	return nil
}
