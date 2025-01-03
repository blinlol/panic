package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


type Card struct {
	Suit Suit
	Val CardVal
}


func (c *Card) Draw(screen *ebiten.Image, x, y float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	text.Draw(screen, fmt.Sprintf("%s %s", SuitToText[c.Suit], ValToText[c.Val]), 
				GeneralFont, op)
}


func ValDelta(c1, c2 Card) uint {
	a := math.Abs(float64(c1.Val - c2.Val))
	return uint(math.Min(a, float64(N_CARDVALS) - a))
}


type Suit int
const (
	Spades Suit = iota
	Clubs
	Hearts
	Diamonds
)
const N_SUITS = 4
var ALL_SUITS [N_SUITS]Suit = [N_SUITS]Suit{Spades, Clubs, Hearts, Diamonds}

var SuitToText map[Suit]string = map[Suit]string{
	Spades: "S",
	Clubs: "C",
	Hearts: "H",
	Diamonds: "D",
}

type CardVal int
const (
	Val_A CardVal = iota
	Val_2
	Val_3
	Val_4
	Val_5
	Val_6
	Val_7
	Val_8
	Val_9
	Val_10
	Val_J
	Val_Q
	Val_K
)
const N_CARDVALS = Val_K - Val_A + 1
var ALL_CARDVALS [N_CARDVALS]CardVal = [N_CARDVALS]CardVal{
	Val_A,
	Val_2, Val_3, Val_4, Val_5, 
	Val_6, Val_7, Val_8, Val_9, 
	Val_10, Val_J, Val_Q, Val_K,
}

var ValToText map[CardVal]string

func init(){
	ValToText = make(map[CardVal]string)
	for v:=2; v<=10; v++ {
		ValToText[CardVal(v-1)] = fmt.Sprintf("%d", v)
	}
	ValToText[Val_J] = "J"
	ValToText[Val_Q] = "Q"
	ValToText[Val_K] = "K"
	ValToText[Val_A] = "A"
}