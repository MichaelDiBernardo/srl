package console

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

type sheetScreen struct {
	display display
	sheet   panel
}

func newSheetScreen(display display) *sheetScreen {
	return &sheetScreen{
		display: display,
		sheet:   newSheetPanel(display),
	}
}

func (s *sheetScreen) Render(g *game.Game) {
	s.sheet.Render(g)
}

func (inv *sheetScreen) Handle(ev game.Event) {
}

func (inv *sheetScreen) NextCommand() game.Command {
	for {
		tboxev := inv.display.PollEvent()
		if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
			return game.ModeCommand{Mode: game.ModeHud}
		}
	}
}

type sheetPanel struct {
	display display
}

func newSheetPanel(display display) *sheetPanel {
	return &sheetPanel{display: display}
}

func (s *sheetPanel) Handle(e game.Event) {
}

func (s *sheetPanel) Render(g *game.Game) {
	p := g.Player
	sheet := g.Player.Sheet

	s.display.Write(1, 2, fmt.Sprintf("NAME   %v", p.Spec.Name), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 3, fmt.Sprintf("RACE   %v", p.Spec.Species.Describe()), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 2, fmt.Sprintf("STR   %d", sheet.Str()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 3, fmt.Sprintf("AGI   %d", sheet.Agi()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 4, fmt.Sprintf("VIT   %d", sheet.Vit()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 5, fmt.Sprintf("MND   %d", sheet.Mnd()), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(1, 7, fmt.Sprintf("TURN %13d", 923), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 8, fmt.Sprintf("XP LEFT %10d", 152), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 9, fmt.Sprintf("XP TOTAL %9d", 4095), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 10, fmt.Sprintf("FLOOR %11dF", g.Floor), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 11, fmt.Sprintf("MIN FLOOR %7dF", 1), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 7, fmt.Sprintf("FIGHT %9s", sheet.Attack().Describe()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 8, fmt.Sprintf("DEF %11s", sheet.Defense().Describe()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 10, fmt.Sprintf("HP %12s", fmt.Sprintf("%d:%d", sheet.HP(), sheet.MaxHP())), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 11, fmt.Sprintf("MP %12s", fmt.Sprintf("%d:%d", sheet.MP(), sheet.MaxMP())), termbox.ColorWhite, termbox.ColorBlack)
}
