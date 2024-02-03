package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	MAX_ITEMS = 80000
)

var rdb *redis.Client

func init() {
	// Initialize the Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Add your Redis password here if applicable
		DB:       0,  // Use the default Redis database
	})

	// Set the Redis key to the value set by the environment variable MAX_ITEMS
	err := rdb.Set(context.Background(), "items_available", MAX_ITEMS, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the counter in Redis is 0
	itemsAvailable, err := rdb.Get(context.Background(), "items_available").Int()
	if err != nil {
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		return
	}
	if itemsAvailable == 0 {
		http.Error(w, "No items available", http.StatusOK)
		return
	}

	for i := 0; i < 10; i++ {
		// Decrement the counter in Redis by one
		err := decrementCounter()
		if err != nil && i == 9 {
			http.Error(w, "Failed to process order", http.StatusOK)
			return
		}
	}

	// Send a response
	fmt.Fprint(w, "Order received successfully")
}

func decrementCounter() error {
	// cheaper way to get a most likely unique value
	// better if can append with client id
	lockValue := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Acquire lock
	lock := rdb.SetNX(context.Background(), "items_available_lock", lockValue, 30*time.Millisecond)
	if lock.Err() != nil {
		return lock.Err()
	}
	if !lock.Val() {
		return fmt.Errorf("Failed to acquire lock")
	}

	// Decrement the counter in Redis by one
	err := rdb.Decr(context.Background(), "items_available").Err()
	if err != nil {
		return err
	}

	if rdb.Get(context.Background(), "items_available_lock").Val() == lockValue {
		// Release lock
		err = rdb.Del(context.Background(), "items_available_lock").Err()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Trying to release lock set by another client")
	}

	return nil
}

func main() {
	http.HandleFunc("/order", orderHandler)

	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
