package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)


const HAND_SIZE = 5
const noneHand = -1
var noneOpen *Deck = nil

type Player struct {
	Id int
	Close *Deck
	Open *Deck
	Hand [HAND_SIZE]*Deck
	
	// TODO MyOpen and OtherOpen


	SelectedHand int
	SelectedOpen *Deck

	Controlling ControlKeys
}


type ControlKeys struct {
	Close ebiten.Key
	Hand [HAND_SIZE]ebiten.Key
	Open1 ebiten.Key
	Open2 ebiten.Key
}


func NewPlayer(d *Deck, id int, layout PlayerLayout, controlling ControlKeys) Player {
	log.Println(layout)
	d.Center = &layout.Close
	o := &Deck{OpenNumber: 0, Center: &layout.Open}
	var h [HAND_SIZE]*Deck
	for i := range HAND_SIZE{
		h[i] = &Deck{
			OpenNumber: 1,
			Center: &layout.Hand[i],
		}
	}
	return Player{
		Id: id,
		Close: d,
		Open: o,
		Hand: h,
		SelectedHand: noneHand,
		SelectedOpen: noneOpen,
		Controlling: controlling,
	}
}


func (p *Player) Draw(screen *ebiten.Image) {
	p.Close.Draw(screen, NOT_SELECTED)
	sel := NOT_SELECTED
	if p.SelectedOpen == p.Open {
		sel = Selected(p.Id)
	}
	p.Open.Draw(screen, sel)
	for i, deck := range p.Hand {
		sel = NOT_SELECTED
		if i == p.SelectedHand {
			sel = Selected(p.Id)
		}
		deck.Draw(screen, sel)
	}
}

/* делит карты из close на остальные колоды */
func (p *Player) LayOutCards() {
	p.Close.Shuffle()
	outer:
	for i := range HAND_SIZE {
		if len(p.Hand[i].Cards) != 0 {
			// TODO мб убрать этот иф
			panic("All cards must be in Close, before call LayOutCards")
		}
		for range i + 1 {
			card := p.Close.PopTop()
			if card == nil {
				break outer
			}

			p.Hand[i].AddCard(card, false)
		}
	}

	for i := range HAND_SIZE {
		p.Hand[i].OpenNumber = 1
	}
}


func (p *Player) OpenClosed() {
	card := p.Close.GetTop()
	p.Close.DeleteTop()
	p.Open.AddCard(card, true)
}


func (p *Player) NumberCardsInHand() int {
	res := 0
	for _, h := range p.Hand {
		res += len(h.Cards)
	}
	return res
}


func (p *Player) NumberCards() int {
	return p.NumberCardsInHand() + len(p.Close.Cards) + len(p.Open.Cards)
}

func (p *Player) HaveMove(other *Deck) bool {
	o1 := p.Open.GetTop()
	o2 := other.GetTop()
	for  _, h := range p.Hand {
		if len(h.Cards) == 0 {
			continue
		}

		if h.OpenNumber == 0 {
			return true
		}

		c := h.GetTop()
		if ValDelta(*c, *o1) == 1 {
			return true
		} else if ValDelta(*c, *o2) == 1 {
			return true
		}
	}

	for i := range HAND_SIZE {
		for j := i + 1 ; j < HAND_SIZE ; j++ {
			li := len(p.Hand[i].Cards)
			lj := len(p.Hand[j].Cards)
			if ! ((li == 0) == (lj == 0)) {
				return true
			}
			if p.Hand[i].GetTop().Val == p.Hand[j].GetTop().Val {
				return true
			}
		}
	}
	return false
}

func (p *Player) SelectHand(i int) {
	log.Printf("p%d select hand %d", p.Id, i)
	if p.SelectedHand == i {
		// remove selection
		p.SelectedHand = noneHand

	} else if p.Hand[i].OpenNumber == 0 && len(p.Hand[i].Cards) > 0 {
		// open closed hand
		p.Hand[i].OpenNumber = 1

	} else if p.SelectedHand != noneHand {
		// already selected some deck
		if len(p.Hand[i].Cards) == 0 && len(p.Hand[p.SelectedHand].Cards) != 0 || 
				len(p.Hand[p.SelectedHand].Cards) != 0 && p.Hand[i].GetTop().Val == p.Hand[p.SelectedHand].GetTop().Val {
			// equal top card vals or empty deck
			card := p.Hand[p.SelectedHand].PopTop()
			p.Hand[i].AddCard(card, true)
			p.SelectedHand = noneHand
		}

	} else {
		p.SelectedHand = i
	}
}

func (p *Player) SelectOpen(open *Deck) {
	log.Printf("p%d select open %v", p.Id, open)
	if p.SelectedOpen == open {
		p.SelectedOpen = noneOpen
	} else {
		p.SelectedOpen = open
	}
}

func (p *Player) ResetSelected() {
	p.SelectedHand = noneHand
	p.SelectedOpen = noneOpen
}

func (p *Player) CanMove() bool {
	if p.SelectedHand != noneHand && p.SelectedOpen != noneOpen {
		if len(p.Hand[p.SelectedHand].Cards) == 0 {
			p.ResetSelected()
			return false
		}
		
		oc := p.SelectedOpen.GetTop()
		hc := p.Hand[p.SelectedHand].GetTop()
		return ValDelta(*oc, *hc) == 1
	}
	return false
}

func (p *Player) MakeMove() {
	log.Printf("p%d make move hand %d open top %d", p.Id, p.SelectedHand, p.SelectedOpen.GetTop()) 
	card := p.Hand[p.SelectedHand].PopTop()
	p.SelectedOpen.AddCard(card, true)
	p.ResetSelected()
}


func (p *Player) GatherOpen(o *Deck) { // TODO rename
	for range len(o.Cards) {
		p.Close.AddCard(o.PopTop(), false)
	}
	for i := range HAND_SIZE {
		for range len(p.Hand[i].Cards) {
			p.Close.AddCard(p.Hand[i].PopTop(), false)
		}
	}
}

// TODO убрать открытую колоду у игрока и положить их в игру?