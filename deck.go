package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


type Deck struct {
	Cards []*Card     // TODO мб сделать  закрытым атрибутом?
	Center *Coords
	OpenNumber int
}

type Selected int
const (
	NOT_SELECTED Selected = iota
	FIRST_SELECTED
	SECOND_SELECTED
)


func NewDeck() Deck {
	deck := Deck{
		Cards: make([]*Card, 0),
	}
	for _, suit := range ALL_SUITS {
		for _, val := range ALL_CARDVALS {
			deck.Cards = append(deck.Cards, &Card{Suit: suit, Val: val})
		}
	}
	return deck
}


func MergeDecks(d1, d2 Deck, openNumber int) Deck {
	return Deck{
		Cards: slices.Concat(d1.Cards, d2.Cards),
		OpenNumber: openNumber,
	}
}


func (d *Deck) Draw(screen *ebiten.Image, selected Selected) {
	x := d.Center.X - GameCfg.Layout.CardW / 2
	y := d.Center.Y - GameCfg.Layout.CardH / 2

	if len(d.Cards) == 0 {
		op := &text.DrawOptions{}
		op.GeoM.Translate(x, y)
		text.Draw(screen, "empty", GeneralFont, op)
	} else {
		op := &text.DrawOptions{}
		op.GeoM.Translate(x, y+25)
		text.Draw(screen, fmt.Sprintf("(%d)", len(d.Cards)), GeneralFont, op)
		// TODO рисовать нижние карты тоже
		if d.OpenNumber == 0 {
			op := &text.DrawOptions{}
			op.GeoM.Translate(x, y)
			text.Draw(screen, "close", GeneralFont, op)
		} else {
			d.Cards[len(d.Cards) - 1].Draw(screen, x, y)	
		}
	}

	if selected == FIRST_SELECTED {
		im := ebiten.NewImage(10, 5)
		im.Fill(color.RGBA{255, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y - 5)
		screen.DrawImage(im, op)
	} else if selected == SECOND_SELECTED {
		im := ebiten.NewImage(10, 5)
		im.Fill(color.RGBA{0, 255, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y + GameCfg.Layout.CardH + 5)
		screen.DrawImage(im, op)
	}
}


func (d *Deck) Shuffle() {
	rand.Shuffle(
		len(d.Cards),
		func (i, j int) {
			d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
		},
	)
}


func (d Deck) Split(firstSize int) (deck1, deck2 Deck) {
	deck1.Cards = make([]*Card, firstSize)
	copy(deck1.Cards, d.Cards[:firstSize])

	secondSize := len(d.Cards) - firstSize
	deck2.Cards = make([]*Card, secondSize)
	copy(deck2.Cards, d.Cards[firstSize:])
	return
}


func (d *Deck) GetTop() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	return d.Cards[len(d.Cards) - 1]
}


func (d *Deck) PopTop() *Card {
	top := d.GetTop()
	d.DeleteTop()
	return top
}


func (d *Deck) DeleteTop() (ok bool) {
	if len(d.Cards) == 0 {
		return
	}
	if d.OpenNumber > 0 {
		d.OpenNumber -= 1
	}
	d.Cards = d.Cards[:len(d.Cards) - 1]
	ok = true
	return
}


func (d *Deck) AddCard(card *Card, isOpen bool) {
	if isOpen {
		d.OpenNumber += 1
	}
	d.Cards = append(d.Cards, card)
}
