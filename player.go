package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/blinlol/panic/internal/queue"
)


const HAND_SIZE = 5
// const noneHand = -1
const noneOpen = 0

type Player struct {
	Id int
	Close *Deck
	// Open *Deck
	Hand [HAND_SIZE]*Deck
	Open [3]*Deck

	SelectedHands *queue.Queue[int]
	SelectedOpen int

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
	o := NewEmptyDeck()
	o.Center = &layout.Open
	open := [3]*Deck{}
	open[id] = &o
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
		Hand: h,
		Open: open,
		SelectedHands: queue.NewQueue[int](2),
		SelectedOpen: noneOpen,
		Controlling: controlling,
	}
}


// must be called after Players creation
func ExchangeOpenDecks(p1, p2 *Player) {
	if p1.Id != 1 || p2.Id != 2 {
		panic("players must be in right order")
	}

	p1.Open[p2.Id] = p2.Open[p2.Id]
	p2.Open[p1.Id] = p1.Open[p1.Id]
}


func (p *Player) Draw(screen *ebiten.Image) {
	p.Close.Draw(screen, true)

	p.Open[p.Id].Draw(screen, false)

	if p.SelectedOpen != noneOpen {
		p.Open[p.SelectedOpen].DrawSelection(screen, Selected(p.Id))
	}

	for i, deck := range p.Hand {
		deck.Draw(screen, true)

		if p.SelectedHands.Contain(i) {
			deck.DrawSelection(screen, Selected(p.Id))
		}
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
	p.Open[p.Id].AddCard(card, true)
}


func (p *Player) NumberCardsInHand() int {
	res := 0
	for _, h := range p.Hand {
		res += len(h.Cards)
	}
	return res
}


func (p *Player) NumberCards() int {
	return p.NumberCardsInHand() + len(p.Close.Cards) + len(p.Open[p.Id].Cards)
}


func (p *Player) HaveMove() bool {
	o1 := p.Open[1].GetTop()
	o2 := p.Open[2].GetTop()
	for  _, h := range p.Hand {
		if len(h.Cards) == 0 {
			continue
		}

		if h.OpenNumber == 0 {
			log.Printf("have move p%d OpenNumber == 0\n", p.Id)
			return true
		}

		c := h.GetTop()
		if ValDelta(*c, *o1) == 1 {
			log.Printf("have move p%d valdelta 1 %v %v = 1\n", p.Id, c, o1)
			return true
		} else if ValDelta(*c, *o2) == 1 {
			log.Printf("have move p%d valdelta 2 %v %v = 1\n", p.Id, c, o2)
			return true
		}
	}

	for i := range HAND_SIZE {
		for j := range HAND_SIZE {
			if i == j {
				continue
			}

			li := len(p.Hand[i].Cards)
			lj := len(p.Hand[j].Cards)

			if li == 0 && lj != 0 && p.Hand[j].OpenNumber != len(p.Hand[j].Cards) {
				log.Printf("have move p%d can open %d %d\n", p.Id, i, j)
				return true
			}
			if li > 0 && lj > 0 && p.Hand[i].GetTop().Val == p.Hand[j].GetTop().Val {
				log.Printf("have move p%d eq vals %d %d\n", p.Id, i, j)
				return true
			}
		}
	}
	return false
}

func (p *Player) SelectHand(i int) {
	log.Printf("p%d select hand %d", p.Id, i)
	if p.SelectedHands.Contain(i) {
		p.SelectedHands.Remove(i)
	} else if p.SelectedHands.Size() == 0 || p.SelectedOpen == noneOpen {
		p.SelectedHands.Push(i)
	} else if p.SelectedHands.Size() == 1 && p.SelectedOpen != noneOpen {
		p.SelectedHands.Pop()
		p.SelectedHands.Push(i)
	}
}


func (p *Player) SelectOpen(open int) {
	log.Printf("p%d select open %v", p.Id, open)
	if p.SelectedOpen == open {
		p.SelectedOpen = noneOpen
	} else {
		p.SelectedOpen = open
	}
}


func (p *Player) ResetSelected() {
	p.SelectedHands.Clear()
	p.SelectedOpen = noneOpen
}


func (p *Player) MakeMove() bool {
	// return true if made move

	if p.SelectedHands.Size() == 1 {
		h := p.Hand[p.SelectedHands.Top()]

		if h.OpenNumber == 0 && len(h.Cards) > 0 {
			h.OpenNumber += 1
			p.SelectedHands.Pop()
			return true

		} else if p.SelectedOpen != noneOpen {
			o := p.Open[p.SelectedOpen]
			if h.OpenNumber > 0 {
				if ValDelta(*o.GetTop(), *h.GetTop()) == 1 {
					card := h.PopTop()
					o.AddCard(card, true)
					p.ResetSelected()
					return true
				}
			}
		}

	} else if p.SelectedHands.Size() == 2 {
		h_from := p.Hand[p.SelectedHands.Pop()]
		h_to := p.Hand[p.SelectedHands.Pop()]

		if len(h_to.Cards) == 0 && len(h_from.Cards) != 0 {
			card := h_from.PopTop()
			h_to.AddCard(card, true)
			p.ResetSelected()
			return true

		} else if len(h_to.Cards) != 0 && len(h_from.Cards) != 0 {
			if h_from.OpenNumber > 0 && h_to.OpenNumber > 0 {
				if ValDelta(*h_to.GetTop(), *h_from.GetTop()) == 0 {
					card := h_from.PopTop()
					h_to.AddCard(card, true)
					p.ResetSelected()
					return true
				}
			}
		}
	}
	return false
}


func (p *Player) GatherOpen(o_ind int) {
	open := p.Open[o_ind]
	for range len(open.Cards) {
		p.Close.AddCard(open.PopTop(), false)
	}
	for i := range HAND_SIZE {
		for range len(p.Hand[i].Cards) {
			p.Close.AddCard(p.Hand[i].PopTop(), false)
		}
	}
}
