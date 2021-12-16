package main

import "testing"

func TestShoeGeneration(t *testing.T) {
	for i := 1; i < 20; i++ {
		shoe := GenerateShoe(i)
		ValidateDeck(t, shoe, i)
	}
}
