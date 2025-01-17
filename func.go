package main

import "context"

// This file will describe many type of functions

// The most basic example
func sampleFunc() {
}

// Function that requires a parameter and returns a value
func funcWithParam(ctx context.Context) error {
	return nil
}

// Functions with multiple parameter and multiple returned value
// userId and orderId are both int, so we can save time by writing it like this
// or if you're more of a classic person than we can do something like
// sampleFunc2(userId int, orderId int, userName string) etc...
func sampleFunc2(userId, orderId int, userName string) (any, error) {
	return nil, nil
}

// See oop.go for next lesson
