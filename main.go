package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var GeneralFont *text.GoTextFace
var GameCfg *GameConfig

//go:embed LeedsUni10-12-13.ttf
var LeedsUni_ttf []byte

// TODO выделить в отдельные пакеты логику

func main(){
	ebiten.SetWindowTitle("Panic!")

	controlPlayer1 := ControlKeys{
		Close: ebiten.KeyX,
		Hand: [HAND_SIZE]ebiten.Key{
			ebiten.KeyA,
			ebiten.KeyS,
			ebiten.KeyD,
			ebiten.KeyF,
			ebiten.KeyG,
		},
		Open1: ebiten.KeyW,
		Open2: ebiten.KeyE,
	}
	controlPlayer2 := ControlKeys{
		Close: ebiten.KeyM,
		Hand: [HAND_SIZE]ebiten.Key{
			ebiten.KeyQuote,
			ebiten.KeySemicolon,
			ebiten.KeyL,
			ebiten.KeyK,
			ebiten.KeyJ,
		},
		Open1: ebiten.KeyI,
		Open2: ebiten.KeyO,
	}


	cfg := &GameConfig{
		GeneralFont: getFont(),
		Layout: getLayout(),
		ControlPlayer1: controlPlayer1,
		ControlPlayer2: controlPlayer2,
	}
	ebiten.SetWindowSize(int(cfg.Layout.W), int(cfg.Layout.H))
	GeneralFont = cfg.GeneralFont
	GameCfg = cfg
	game := NewGame(cfg)

	if err := ebiten.RunGame(game) ; err != nil {
		panic(err)
	}
}


func getFont() *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(LeedsUni_ttf))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size: 25,
	}
}


func getLayout() ScreenLayout {
	// TODO задавать размеры в зависимости от размера экрана
	l := ScreenLayout{}
	l.W = 1280
	l.H = 720

	l.CardW = l.W * 0.1
	l.CardH = l.H * 0.1

	l.A = l.CardW / 2
	l.B = l.CardW / 2
	l.C = l.CardH  * (1 + 1)

	delta_oc := Coords{
		X: 2 * l.A + l.CardW,
	}
	delta_hand := Coords{
		X: l.B + l.CardW,
	}
	mid := HAND_SIZE / 2
	
	p1 := PlayerLayout{}
	p1.Open = Coords{
		X: l.W / 2 - l.A - l.CardW / 2,
		Y: l.H / 2,
	}
	p1.Close = p1.Open.Add(delta_oc.Neg())
	p1.Hand[mid] = Coords{
		X: l.W / 2,
		Y: l.H / 2 - l.C,
	}
	for di := 1; di <= mid; di += 1 {
		p1.Hand[mid - di] = p1.Hand[mid].Add(delta_hand.Neg().Mul(float64(di)))
		if mid + di < HAND_SIZE {
			p1.Hand[mid + di] = p1.Hand[mid].Add(delta_hand.Mul(float64(di)))
		}
	}
	l.P1Layout = p1

	p2 := PlayerLayout{}
	p2.Open = Coords{
		X: l.W / 2 + l.A + l.CardW / 2,
		Y: l.H / 2,
	}
	p2.Close = p2.Open.Add(delta_oc)
	p2.Hand[mid] = Coords{
		X: l.W / 2,
		Y: l.H / 2 + l.C,
	}
	for di := 1; di <= mid; di += 1 {
		p2.Hand[mid - di] = p2.Hand[mid].Add(delta_hand.Mul(float64(di)))
		if mid + di < HAND_SIZE {
			p2.Hand[mid + di] = p2.Hand[mid].Add(delta_hand.Neg().Mul(float64(di)))
		}
	}
	l.P2Layout = p2

	return l
}



