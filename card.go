package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/blinlol/panic/resources/cards"
)

type Card struct {
	Suit Suit
	Val CardVal
}


func (c *Card) Draw(screen *ebiten.Image, x, y float64, tetha float64, isOpen bool) {
	// op := &text.DrawOptions{}
	// op.GeoM.Translate(x, y)
	// text.Draw(screen, fmt.Sprintf("%s %s", SuitToText[c.Suit], ValToText[c.Val]), 
	// 			GeneralFont, op)

	cw, ch := GameCfg.Layout.CardW, GameCfg.Layout.CardH
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cw / 2, -ch / 2)
	op.GeoM.Rotate(tetha)
	op.GeoM.Translate(cw / 2, ch / 2)
	op.GeoM.Translate(x, y)
	if isOpen {
		screen.DrawImage(images[c.Suit].SubImage(rectVal[c.Val]).(*ebiten.Image), op)
	} else {
		screen.DrawImage(imageBack, op)
	}
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

const NUM_CARDS = N_CARDVALS * N_SUITS

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


var (
	images map[Suit]*ebiten.Image = make(map[Suit]*ebiten.Image)
	rectVal map[CardVal]image.Rectangle = make(map[CardVal]image.Rectangle)
	imageBack *ebiten.Image
)

func init(){
	tmp := func (suit Suit, src []byte) {
		img, _, err := image.Decode(bytes.NewReader(src))
		if err != nil {
			log.Fatal(err)
		}
		images[suit] = ebiten.NewImageFromImage(img)
	}

	suitToSrc := map[Suit][]byte {
		Clubs: cards.ImageClubsSrc,
		Hearts: cards.ImageHeartsSrc,
		Diamonds: cards.ImageDiamondsSrc,
		Spades: cards.ImageSpadesSrc,
	}
	for suit, src := range suitToSrc {
		tmp(suit, src)
	}

	size := images[Spades].Bounds().Size()
	wc := size.X / 5
	hc := size.Y / 3

	for i, val := range ALL_CARDVALS {
		x0, y0 := (i % 5) * wc, (i / 5) * hc
		x1, y1 := x0 + wc, y0 + hc
		rectVal[val] = image.Rect(x0, y0, x1, y1)
	}

	img, _, err := image.Decode(bytes.NewReader(cards.ImageCardBackSrc))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImg := ebiten.NewImageFromImage(img) 
	s := ebitenImg.Bounds().Size()
	halfRect := image.Rect(0, 0, s.X / 2, s.Y)
	imageBack = ebitenImg.SubImage(halfRect).(*ebiten.Image)
}
