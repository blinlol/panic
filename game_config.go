package main

import "github.com/hajimehoshi/ebiten/v2/text/v2"

type GameConfig struct {
	Layout ScreenLayout
	GeneralFont *text.GoTextFace

	ControlPlayer1 ControlKeys
	ControlPlayer2 ControlKeys
}


type ScreenLayout struct {
	P1Layout PlayerLayout
	P2Layout PlayerLayout

	W, H float64
	A, B, C float64
	CardH, CardW float64
}


type PlayerLayout struct {
	Hand [HAND_SIZE]Coords
	Close Coords

	// coords of centers for open decks
	Open Coords
}

