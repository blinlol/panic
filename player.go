package main

import (
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/blinlol/panic/internal/queue"
)


const HAND_SIZE = 5
// const noneHand = -1
const noneOpen = 0

const controlW, controlH = 30, 30

type Player struct {
	Id int
	Close *Deck
	// Open *Deck
	Hand [HAND_SIZE]*Deck
	Open [3]*Deck

	SelectedHands *queue.Queue[int]
	SelectedOpen int

	Controlling ControlKeys

	Animations []Animation
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

	p.drawControlling(screen)

	for _, anim := range p.Animations {
		anim.Draw(screen)
	}
}


func (p *Player) drawControlling(screen *ebiten.Image) {
	deltaClose := Coords{
		X: - GameCfg.Layout.CardW / 2 - 10,
	}
	deltaOpen := Coords{
		Y: - GameCfg.Layout.CardH / 2 - 20,
	}
	deltaHand := Coords{
		Y: - GameCfg.Layout.CardH / 2 - 10,
	}

	if p.Id == 2 {
		deltaClose = deltaClose.Neg().Add(Coords{X: 10})
		deltaOpen = deltaOpen.Neg()
		deltaHand = deltaHand.Neg().Add(Coords{Y: 10})
	}

	drawButton(screen, p.Controlling.Close, p.Close.Center.Add(deltaClose))
	drawButton(screen, p.Controlling.Open1, p.Open[1].Center.Add(deltaOpen))
	drawButton(screen, p.Controlling.Open2, p.Open[2].Center.Add(deltaOpen))

	for i, h := range p.Hand {
		drawButton(screen, p.Controlling.Hand[i], h.Center.Add(deltaHand))
	}
}


func drawButton(screen *ebiten.Image, key ebiten.Key, center Coords) {
	btn := ebiten.NewImage(controlW, controlH)
	text.Draw(btn, key.String(), GameCfg.GeneralFont, nil)
	op := &ebiten.DrawImageOptions{}
	// TODO рамку для буквы
	op.GeoM.Translate(center.X - controlW / 2, center.Y - controlH / 2)
	screen.DrawImage(btn, op)
}


func (p *Player) Update(){
	i := 0
	for i < len(p.Animations) {
		finished := p.Animations[i].Update()
		if finished {
			p.Animations = slices.Delete(p.Animations, i, i+1)
		} else {
			i += 1
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

			// TODO пофиксить так, чтобы не переключалось состояние пока не разложат карты. Возможно стоит
			// if len(p.Close.Cards) == 0 {
			// 	break outer
			// }
			// p.Animations = append(p.Animations, NewAnimMove(p.Close, p.Hand[i], false))
			
		}
	}

	for i := range HAND_SIZE {
		p.Hand[i].OpenNumber = 1
	}
}


func (p *Player) OpenClosed() {
	// card := p.Close.GetTop()
	// p.Close.DeleteTop()
	// p.Open[p.Id].AddCard(card, true)

	p.Animations = append(p.Animations, NewAnimMove(p.Close, p.Open[p.Id], true))
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
		if o1 != nil && ValDelta(*c, *o1) == 1 {
			log.Printf("have move p%d valdelta 1 %v %v = 1\n", p.Id, c, o1)
			return true
		} 
		if o2 != nil && ValDelta(*c, *o2) == 1 {
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
			// open top card in hand
			h.OpenNumber += 1
			p.SelectedHands.Pop()
			return true

		} else if p.SelectedOpen != noneOpen {
			// move card from hand to selected open deck if can
			o := p.Open[p.SelectedOpen]
			if h.OpenNumber > 0 && len(o.Cards) > 0 {
				if ValDelta(*o.GetTop(), *h.GetTop()) == 1 {
					p.Animations = append(p.Animations, NewAnimMove(h, o, true))
					// card := h.PopTop()
					// o.AddCard(card, true)
					p.ResetSelected()
					return true
				}
			}
		}

	} else if p.SelectedHands.Size() == 2 {
		h_from := p.Hand[p.SelectedHands.Pop()]
		h_to := p.Hand[p.SelectedHands.Pop()]

		if len(h_to.Cards) == 0 && len(h_from.Cards) != 0 {
			// move top card to empty hand
			p.Animations = append(p.Animations, NewAnimMove(h_from, h_to, true))

			// card := h_from.PopTop()
			// h_to.AddCard(card, true)
			p.ResetSelected()

			return true

		} else if len(h_to.Cards) != 0 && len(h_from.Cards) != 0 {
			// union equal val cards in hand
			if h_from.OpenNumber > 0 && h_to.OpenNumber > 0 {
				if ValDelta(*h_to.GetTop(), *h_from.GetTop()) == 0 {
					p.Animations = append(p.Animations, NewAnimMove(h_from, h_to, true))
					// card := h_from.PopTop()
					// h_to.AddCard(card, true)
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
