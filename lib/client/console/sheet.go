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
	s.display.Write(1, 2, fmt.Sprintf("NAME   %v", "Debo"), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 3, fmt.Sprintf("RACE   %v", "Human"), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 2, fmt.Sprintf("STR   %d", 2), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 3, fmt.Sprintf("AGI   %d", 2), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 4, fmt.Sprintf("VIT   %d", 2), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 5, fmt.Sprintf("MND   %d", 2), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(1, 7, fmt.Sprintf("TURN %13d", 923), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 8, fmt.Sprintf("XP LEFT %10d", 152), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 9, fmt.Sprintf("XP TOTAL %9d", 4095), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 10, fmt.Sprintf("FLOOR %12d", 2), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 11, fmt.Sprintf("MIN FLOOR %8d", 1), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 7, fmt.Sprintf("FIGHT %9s", "(+1,2d7)"), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 8, fmt.Sprintf("DEF %11s", "[+2,1-4]"), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 10, fmt.Sprintf("HP %12s", "20:20"), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 11, fmt.Sprintf("MP %12s", "20:20"), termbox.ColorWhite, termbox.ColorBlack)
}
