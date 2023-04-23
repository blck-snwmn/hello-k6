package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type order struct {
	ID       string `json:"id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Shipped  bool   `json:"shipped"`
}

type orderStore struct {
	orders map[string]order
	mux    sync.RWMutex
}

func (s *orderStore) addOrder(o order) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.orders[o.ID] = o
}

func (s *orderStore) updateOrder(o order) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	oo, ok := s.orders[o.ID]
	if !ok {
		return false
	}
	oo.Shipped = o.Shipped
	s.orders[o.ID] = oo
	return true
}

func (s *orderStore) getOrders() []order {
	s.mux.RLock()
	defer s.mux.RUnlock()
	orders := make([]order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

var store = &orderStore{orders: map[string]order{}}

func main() {
	http.HandleFunc("/orders", ordersHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getOrders(w, r)
	case "POST":
		addOrder(w, r)
	case "PUT":
		updateOrder(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	var order order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Printf("failed to add order: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.addOrder(order)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	var order order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Printf("failed to update order: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ok := store.updateOrder(order)
	if !ok {
		fmt.Println("not found")
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	orders := store.getOrders()
	json.NewEncoder(w).Encode(orders)
}
