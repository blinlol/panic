package main

import (
	"fmt"
	"image/color"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	GS_START GameState = iota
	GS_STARTING
	GS_DIVIDE
	GS_NEW_ROUND
	GS_LAYOUT
	GS_OPEN
	GS_CHECK_CARDS
	GS_PLAY
	GS_FASTER
)

type Game struct {
	Config *GameConfig
	// TODO rename eg to Alice & Bob
	Player1 Player
	Player2 Player

	Keys            []ebiten.Key
	JustPressedKeys []ebiten.Key

	State      GameState
	Animations []Animation

	// flag true if player open top card
	openStateFlag [2]bool
}

func NewGame(cfg *GameConfig) *Game {
	deck := NewDeck()
	deck.Shuffle()
	d1, d2 := deck.Split(len(deck.Cards) / 2)
	return &Game{
		Config:  cfg,
		Player1: NewPlayer(&d1, 1, cfg.Layout.P1Layout, cfg.ControlPlayer1),
		Player2: NewPlayer(&d2, 2, cfg.Layout.P2Layout, cfg.ControlPlayer2),
	}
}

func (g *Game) Update() error {
	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])
	g.JustPressedKeys = inpututil.AppendJustPressedKeys(g.JustPressedKeys[:0])

	// ctrl+w -> exit
	if slices.Contains(g.Keys, ebiten.KeyControl) &&
		slices.Contains(g.Keys, ebiten.KeyW) {
		return ebiten.Termination
	}

	switch g.State {

	case GS_START:
		g.State = GS_STARTING
		g.Animations = append(g.Animations,
			NewAnimStart(60*1, func() { g.State = GS_DIVIDE }))

	case GS_STARTING:

	case GS_DIVIDE:
		// deck already divided if game is new
		g.State = GS_NEW_ROUND
		g.Animations = append(g.Animations,
			NewAnimNewRound(60*1, func() { g.State = GS_LAYOUT }),
		)

	case GS_NEW_ROUND:

	case GS_LAYOUT:
		g.Player1.LayOutCards()
		g.Player2.LayOutCards()
		g.State = GS_OPEN

	case GS_OPEN:
		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Close) && !g.openStateFlag[0] {
			g.Player1.OpenClosed()
			g.openStateFlag[0] = true
		}
		if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Close) && !g.openStateFlag[1] {
			g.Player2.OpenClosed()
			g.openStateFlag[1] = true
		}
		if g.openStateFlag[0] && g.openStateFlag[1] {
			g.State = GS_CHECK_CARDS
			g.openStateFlag = [2]bool{false, false}
		}

	case GS_CHECK_CARDS:
		n1 := g.Player1.NumberCardsInHand()
		n2 := g.Player2.NumberCardsInHand()
		if n1 == 0 || n2 == 0 {
			g.State = GS_FASTER
		} else {
			if g.Player1.HaveMove(g.Player2.Open) || g.Player2.HaveMove(g.Player1.Open) {
				g.State = GS_PLAY
			} else {
				g.State = GS_OPEN
			}
		}

	case GS_PLAY:
		for i := range HAND_SIZE {
			if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Hand[i]) {
				g.Player1.SelectHand(i)
			}
			if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Hand[i]) {
				g.Player2.SelectHand(i)
			}
		}

		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open1) {
			g.Player1.SelectOpen(g.Player1.Open)
		}
		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open2) {
			g.Player1.SelectOpen(g.Player2.Open)
		}

		if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open1) {
			g.Player2.SelectOpen(g.Player1.Open)
		}
		if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open2) {
			g.Player2.SelectOpen(g.Player2.Open)
		}

		if g.Player1.CanMove() {
			g.Player1.MakeMove()
			g.State = GS_CHECK_CARDS
		}
		if g.Player2.CanMove() {
			g.Player2.MakeMove()
			g.State = GS_CHECK_CARDS
		}

	case GS_FASTER:
		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open1) ||
			slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open2) {
			g.Player1.GatherOpen(g.Player1.Open)
			g.Player2.GatherOpen(g.Player2.Open)
		} else {
			g.Player1.GatherOpen(g.Player2.Open)
			g.Player2.GatherOpen(g.Player1.Open)
		}

		if g.Player1.NumberCards() == 0 || g.Player2.NumberCards() == 0 {
			winner := g.Player1
			if g.Player2.NumberCards() == 0 {
				winner = g.Player2
			}
			log.Printf("winner %d\n", winner.Id)
			return ebiten.Termination
		}
		g.State = GS_NEW_ROUND

	default:
		return fmt.Errorf("unknown GameState in Game.Update() %d", g.State)
	}

	i := 0
	for i < len(g.Animations) {
		finished := g.Animations[i].Update()
		if finished {
			g.Animations = slices.Delete(g.Animations, i, i+1)
		} else {
			i += 1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTable(screen)
	for _, anim := range g.Animations {
		anim.Draw(screen)
	}
}

func (g *Game) DrawTable(screen *ebiten.Image) {
	green := color.RGBA{0, 200, 0, 255}
	screen.Fill(green)

	g.Player1.Draw(screen)
	g.Player2.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeigth int) (screenWidth, screenHeight int) {
	screenWidth = outsideWidth
	screenHeight = outsideHeigth
	return
}

// TODO заменить константные продолжительности анимаций на зависимость от фпс
// TODO switch g.State заменить на что-то похожее на классы вершин в графе
// TODO перекладывание между колодами в руке
// TODO открытие верхней карты в руке по нажатию
// TODO одновременное открытие
// TODO объединение одинаковых карт в руке
// TODO валидация ошибочных нажатий
