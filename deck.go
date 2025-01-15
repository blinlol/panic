package main

import (
	"bytes"
	// "fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/blinlol/panic/resources/cards"
)


type Deck struct {
	Cards []*Card
	Center *Coords
	OpenNumber int

	tethas [NUM_CARDS]float64
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
	for i := range deck.tethas {
		deck.tethas[i] = randTetha()
	}

	return deck
}


func NewEmptyDeck() Deck {
	deck := Deck{}
	for i := range deck.tethas {
		deck.tethas[i] = randTetha()
	}
	return deck
}


func randTetha() float64 {
	return (rand.Float64() - 0.5) * math.Pi / 6.
}

func MergeDecks(d1, d2 Deck, openNumber int) Deck {
	return Deck{
		Cards: slices.Concat(d1.Cards, d2.Cards),
		OpenNumber: openNumber,
	}
}


func (d *Deck) Draw(screen *ebiten.Image, evenly bool) {
	x, y := d.getLeftUpperCornerXY()

	if len(d.Cards) == 0 {
		// op := &text.DrawOptions{}
		// op.GeoM.Translate(x, y)
		// text.Draw(screen, "empty", GeneralFont, op)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(imageEmpty, op)

	} else {
		// op := &text.DrawOptions{}
		// op.GeoM.Translate(x, y - 25)
		// text.Draw(screen, fmt.Sprintf("(%d)", len(d.Cards)), GeneralFont, op)

		xx, yy := x, y
		for i, card := range d.Cards {
			tetha, isOpen := 0., i >= len(d.Cards) - d.OpenNumber
			if !evenly {
				tetha = d.tethas[i]
			} else {
				yy += 2
			}
			card.Draw(screen, xx, yy, tetha, isOpen)
		}

		// if d.OpenNumber == 0 {
		// 	// op := &text.DrawOptions{}
		// 	// op.GeoM.Translate(x, y)
		// 	// text.Draw(screen, "close", GeneralFont, op)

		// 	op := &ebiten.DrawImageOptions{}
		// 	op.GeoM.Translate(x, y)
		// 	screen.DrawImage(imageBack, op)
		// } else {
		// 	d.Cards[len(d.Cards) - 1].Draw(screen, x, y)	
		// }
	}
}


func (d *Deck) DrawSelection(screen *ebiten.Image, selected Selected) {
	x, y := d.getLeftUpperCornerXY()
	if selected == FIRST_SELECTED {
		im := ebiten.NewImage(int(GameCfg.Layout.CardW), 5)
		im.Fill(color.RGBA{255, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y - 5)
		screen.DrawImage(im, op)
	} else if selected == SECOND_SELECTED {
		im := ebiten.NewImage(int(GameCfg.Layout.CardW), 5)
		im.Fill(color.RGBA{255, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y + GameCfg.Layout.CardH + 5)
		screen.DrawImage(im, op)
	}
}


func (d *Deck) getLeftUpperCornerXY() (float64, float64){
	x := d.Center.X - GameCfg.Layout.CardW / 2
	y := d.Center.Y - GameCfg.Layout.CardH / 2
	return x, y
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
	if card == nil {
		panic("add nil card to deck")
	}
	if isOpen {
		d.OpenNumber += 1
	}
	d.Cards = append(d.Cards, card)
}


var (
	imageEmpty *ebiten.Image
)

func init(){
	img, _, err := image.Decode(bytes.NewReader(cards.ImageEmptyDeckSrc))
	if err != nil {
		log.Fatal(err)
	}
	imageEmpty = ebiten.NewImageFromImage(img)
}