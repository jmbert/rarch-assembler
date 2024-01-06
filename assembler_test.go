package main

import (
	"fmt"
	"testing"
)

func TestAssemble(t *testing.T) {
	bytes, err := Assemble("tests/test1.s")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("[")
	for _, b := range bytes {
		fmt.Printf("%.2X ", b)
	}
	fmt.Printf("]\n")
}
