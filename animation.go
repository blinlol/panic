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
	callback func()

}

func NewAnimMove(callback func()) *AnimMove {
	return &AnimMove{
		callback: callback,
	}
}

func (a *AnimMove) Update() (finished bool) {
	return
}

func (a *AnimMove) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "AnimMove")
}




// type AnimOpenClosed struct {
// 	duration uint
// 	callback func()
// }

// func NewAnimOpenClosed(duration uint, callback func()) *AnimOpenClosed {
// 	return &AnimOpenClosed{
// 		duration: duration,
// 		callback: callback,
// 	}
// }

// func (a *AnimOpenClosed) Update() (finished bool) {
// 	a.duration -= 1
// 	if a.duration <= 0 {
// 		a.callback()
// 		finished = true
// 	}
// 	return
// }

// func (a *AnimOpenClosed) Draw(screen *ebiten.Image) {
// 	ebitenutil.DebugPrint(screen, fmt.Sprintf("AnimOpenClosed %d", a.duration))
// }