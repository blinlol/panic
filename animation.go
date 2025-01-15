package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Animation interface {
	Draw(*ebiten.Image)
	// returns whether the animation is finished
	Update() bool
}


type AnimStart struct {
	duration uint
	callback func()
}

func NewAnimStart(duration uint, callback func()) *AnimStart {
	return &AnimStart{
		duration: duration,
		callback: callback,
	}
}

func (a *AnimStart) Update() (finished bool) {
	a.duration -= 1
	if a.duration <= 0 {
		a.callback()
		finished = true
	}
	return
}

func (a *AnimStart) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("AnimStart %d", a.duration))
}



type AnimNewRound struct {
	duration uint
	callback func()
}

func NewAnimNewRound(duration uint, callback func()) *AnimNewRound {
	return &AnimNewRound{
		duration: duration,
		callback: callback,
	}
}

func (a *AnimNewRound) Update() (finished bool) {
	a.duration -= 1
	if a.duration <= 0 {
		a.callback()
		finished = true
	}
	return
}

func (a *AnimNewRound) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("AnimNewRound %d", a.duration))
}



type AnimMove struct {
	from *Deck
	to *Deck
	card *Card
	v Coords
	t float64
	duration float64
	isOpen bool
}

func NewAnimMove(from, to *Deck, isOpen bool) *AnimMove {
	speed := 30.
	sub := to.Center.Add(from.Center.Neg())
	dist := to.Center.Distance(*from.Center)
	v := sub.Mul(speed / dist)
	
	card := from.PopTop()
	duration := dist / speed
	if card == nil {
		duration = 0
	}

	return &AnimMove{
		from: from,
		to: to,
		card: card,
		v: v,
		t: 0,
		duration: duration,
		isOpen: isOpen,
	}
}

func (a *AnimMove) Update() (finished bool) {
	a.t += 1
	if a.t > a.duration {
		finished = true
		a.to.AddCard(a.card, a.isOpen)
	}
	return
}

func (a *AnimMove) Draw(screen *ebiten.Image) {
	p := a.from.Center.Add(a.v.Mul(a.t))
	a.card.Draw(screen, p.X, p.Y, 0, a.isOpen)
}
