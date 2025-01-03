package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardVal(t *testing.T) {
	c2 := Card{Val: Val_2, Suit: Spades}
	c3 := Card{Val: Val_3, Suit: Spades}
	c4 := Card{Val: Val_4, Suit: Spades}
	cA := Card{Val: Val_A, Suit: Spades}
	cK := Card{Val: Val_K, Suit: Spades}
	c10 := Card{Val: Val_10, Suit: Spades}
	cJ := Card{Val: Val_J, Suit: Spades}

	assert.Equal(t, uint(1), ValDelta(c2, c3))
	assert.Equal(t, uint(1), ValDelta(c3, c4))
	assert.Equal(t, uint(2), ValDelta(c2, c4))
	assert.Equal(t, uint(2), ValDelta(c4, c2))
	assert.Equal(t, uint(1), ValDelta(c2, cA))
	assert.Equal(t, uint(2), ValDelta(cK, c2))
	assert.Equal(t, uint(1), ValDelta(cA, cK))
	assert.Equal(t, uint(1), ValDelta(c10, cJ))
	assert.Equal(t, uint(5), ValDelta(c10, c2))
}