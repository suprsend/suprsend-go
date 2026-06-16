package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func contextCancellationExample() {
	contextTimeoutExample()
	contextAbortExample()
	contextExpiredExample()
}

func contextTimeoutExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	// ================= Fetch existing tenant by ID
	tenant1, err := suprClient.Tenants.Get(ctx, "__tenant_id__")
	if err != nil {
		checkErr(err)
		return
	}
	log.Println(tenant1)
}

func contextAbortExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // External context abort simulation
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()
	defer cancel()
	// ================= Fetch existing tenant by ID
	tenant1, err := suprClient.Tenants.Get(ctx, "__tenant_id__")
	if err != nil {
		checkErr(err)
		return
	}
	log.Println(tenant1)
}

func contextExpiredExample() {
	suprClient, err := getSuprsendClient()
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately to simulate expired context
	// ================= Fetch existing tenant by ID
	tenant1, err := suprClient.Tenants.Get(ctx, "__tenant_id__")
	if err != nil {
		checkErr(err)
		return
	}
	log.Println(tenant1)
}

func checkErr(err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("Request timed out")
	} else if errors.Is(err, context.Canceled) {
		log.Println("Request canceled")
	} else {
		log.Println(err)
	}
}
