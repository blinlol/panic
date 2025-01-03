package main

import (
	"encoding/json"
	"os"
	"testing"
	"log"

	// "github.com/stretchr/testify/assert"
)


func TestHaveMove(t *testing.T) {
	game := Game{}
	data, err := os.ReadFile("save/2025-01-03 14:37:46.240538661 +0300 MSK m=+3.286972978.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &game)
	if err != nil {
		panic(err)
	}

	log.Printf("%v\n", game)
}