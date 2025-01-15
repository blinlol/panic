package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"slices"
	"time"

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
	GS_CHECK_WIN
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
	game := &Game{
		Config:  cfg,
		Player1: NewPlayer(&d1, 1, cfg.Layout.P1Layout, cfg.ControlPlayer1),
		Player2: NewPlayer(&d2, 2, cfg.Layout.P2Layout, cfg.ControlPlayer2),
	}
	ExchangeOpenDecks(&game.Player1, &game.Player2)
	return game
}

func (g *Game) Update() error {
	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])
	g.JustPressedKeys = inpututil.AppendJustPressedKeys(g.JustPressedKeys[:0])

	// ctrl+w -> exit
	if slices.Contains(g.Keys, ebiten.KeyControl) &&
		slices.Contains(g.Keys, ebiten.KeyW) {
		return ebiten.Termination
	}
	if slices.Contains(g.JustPressedKeys, ebiten.Key0) {
		g.save()
	}

	switch g.State {

	case GS_START:
		g.State = GS_STARTING
		g.Animations = append(g.Animations,
			NewAnimStart(5*1, func() { g.State = GS_DIVIDE }))

	case GS_STARTING:

	case GS_DIVIDE:
		// deck already divided if game is new
		g.State = GS_NEW_ROUND
		g.Animations = append(g.Animations,
			NewAnimNewRound(5*1, func() { g.State = GS_LAYOUT }),
		)

	case GS_NEW_ROUND:

	case GS_LAYOUT:
		g.Player1.LayOutCards()
		g.Player2.LayOutCards()
		g.State = GS_OPEN

	case GS_OPEN:
		if len(g.Player1.Close.Cards) == 0 {
			g.openStateFlag[0] = true
		} else if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Close) && !g.openStateFlag[0] {
			g.Player1.OpenClosed()
			g.openStateFlag[0] = true
		}

		if len(g.Player2.Close.Cards) == 0 {
			g.openStateFlag[1] = true
		} else if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Close) && !g.openStateFlag[1] {
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
			g.Player1.ResetSelected()
			g.Player2.ResetSelected()
			g.State = GS_FASTER
		} else {
			if g.Player1.HaveMove() || g.Player2.HaveMove() {
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
			g.Player1.SelectOpen(1)
		}
		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open2) {
			g.Player1.SelectOpen(2)
		}

		if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open1) {
			g.Player2.SelectOpen(1)
		}
		if slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open2) {
			g.Player2.SelectOpen(2)
		}

		somethingChanged := g.Player1.MakeMove()
		somethingChanged = g.Player2.MakeMove() || somethingChanged
		if somethingChanged {
			g.State = GS_CHECK_CARDS
		}

	case GS_FASTER:
		if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open1) ||
			slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open2) {
			g.Player1.GatherOpen(1)
			g.Player2.GatherOpen(2)
			g.State = GS_CHECK_WIN
		} else if slices.Contains(g.JustPressedKeys, g.Player1.Controlling.Open2) ||
					slices.Contains(g.JustPressedKeys, g.Player2.Controlling.Open1) {
			g.Player1.GatherOpen(2)
			g.Player2.GatherOpen(1)
			g.State = GS_CHECK_WIN
		}

	case GS_CHECK_WIN:
		if g.Player1.NumberCards() == 0 || g.Player2.NumberCards() == 0 {
			winner := g.Player1
			if g.Player2.NumberCards() == 0 {
				winner = g.Player2
			}
			log.Printf("winner %d\n", winner.Id)
			return ebiten.Termination
		}
		g.State = GS_NEW_ROUND
		g.Animations = append(g.Animations,
			NewAnimNewRound(60*1, func() { g.State = GS_LAYOUT }),
		)

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

	g.Player1.Update()
	g.Player2.Update()

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


func (g *Game) save() {
	marshaled, err := json.MarshalIndent(
		g, "", "    ",
	)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	fname := fmt.Sprintf("%v.json", now)
	err = os.WriteFile("save/" + fname, marshaled, 0666)
	if err != nil {
		panic(err)
	}
}

// TODO заменить константные продолжительности анимаций на зависимость от фпс (мб это неправильно, потому что нужнт зависеть от рпс)
// TODO switch g.State заменить на что-то похожее на классы вершин в графе
// TODO валидация ошибочных нажатий
// TODO подумать над тем, чтобы в плеере был метод апдейт который делает все, что происходит тут, связанное с игроком
