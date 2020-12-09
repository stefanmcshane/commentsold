package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/stefanmcshane/commentsold/db"
)

func logErrorAndReturn(ctx context.Context, w http.ResponseWriter, message string, httpStatusCode int) {
	lg(ctx).Errorf(message)
	http.Error(w, message, httpStatusCode)
}

func signin(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *dbHandler) listInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	inv, err := h.db.GetInventoryItems(ctx, 0, 100)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error reading from database - %s", err.Error()), http.StatusServiceUnavailable)
		return
	}

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error encoding inventory response - %s", err.Error()), http.StatusServiceUnavailable)
		return
	}

}

func readInventoryIDFromURL(urlPath string) (int, error) {
	splitURL := strings.Split(urlPath, "/api/inventory/")
	validURLParts := []string{}
	for _, v := range splitURL {
		if v != "" {
			validURLParts = append(validURLParts, v)
		}
	}
	if len(validURLParts) == 0 {
		return -1, errors.New("Unable to parse inventory id from the URL")
	}

	potentialID := strings.Split(validURLParts[0], "/")
	if len(potentialID) == 0 {
		return -1, errors.New("Unable to parse inventory id from the URL")
	}

	ID, err := strconv.Atoi(potentialID[0])
	if err != nil {
		return -1, errors.New("Provided ID must be integer")
	}
	return ID, nil
}

func (h *dbHandler) manageInventoryItemWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getInventoryItemByID(w, r)
		break
	case http.MethodPut:
		if strings.Contains(r.URL.Path, "adjust") {
			h.updateInventoryItemQuantity(w, r)
			break
		}
		h.updateInventoryItem(w, r)
		break
	default:
		logErrorAndReturn(r.Context(), w, "", http.StatusMethodNotAllowed)
		break
	}
}

func (h *dbHandler) getInventoryItemByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID, err := readInventoryIDFromURL(r.URL.Path)
	if err != nil {
		logErrorAndReturn(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	inv, err := h.db.GetInventoryItemByID(ctx, ID)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error reading inventory item %d from database - %s", ID, err.Error()), http.StatusServiceUnavailable)
		return
	}

	err = json.NewEncoder(w).Encode(inv)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error encoding inventory response - %s", err.Error()), http.StatusServiceUnavailable)
		return
	}

}
func (h *dbHandler) updateInventoryItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		logErrorAndReturn(ctx, w, "Must use put request", http.StatusBadRequest)
		return
	}

	ID, err := readInventoryIDFromURL(r.URL.Path)
	if err != nil {
		logErrorAndReturn(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedItem *db.InventoryItem
	err = json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error decoding inventory item - %s", err.Error()), http.StatusBadRequest)
		return
	}
	if updatedItem == nil {
		logErrorAndReturn(ctx, w, "Please provide updated inventory item in body", http.StatusBadRequest)
		return
	}
	updatedItem.ID = ID

	err = h.db.UpdateInventoryItem(ctx, *updatedItem)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error updating inventory item %d in database - %s", ID, err.Error()), http.StatusServiceUnavailable)
		return
	}
}

type adjustQuantity struct {
	Adjustment int `json:"adjustment"`
}

func (h *dbHandler) updateInventoryItemQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		logErrorAndReturn(ctx, w, "Must use put request to adjust quantity", http.StatusBadRequest)
		return
	}

	ID, err := readInventoryIDFromURL(r.URL.Path)
	if err != nil {
		logErrorAndReturn(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	var adjustment adjustQuantity
	err = json.NewDecoder(r.Body).Decode(&adjustment)
	if err != nil {
		logErrorAndReturn(ctx, w, "Error decoding adjustment quantity", http.StatusBadRequest)
		return
	}

	if adjustment.Adjustment == 0 {
		// no adjustment required
		w.WriteHeader(http.StatusOK)
		return
	}

	err = h.db.UpdateInventoryItemStock(ctx, ID, adjustment.Adjustment)
	if err != nil {
		logErrorAndReturn(ctx, w, fmt.Sprintf("Error updating inventory item quantity %d in database - %s", ID, err.Error()), http.StatusServiceUnavailable)
		return
	}
}
