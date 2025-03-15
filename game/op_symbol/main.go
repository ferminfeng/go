package main

import (
	"log"

	"testGo/game/op_symbol/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Tactical Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
