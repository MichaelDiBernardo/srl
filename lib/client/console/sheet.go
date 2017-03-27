package console

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

func newSheetScreen(display display) *screen {
	return &screen{
		display: display,
		panels: []panel{
			newSkillEditPanel(display),
			newSheetPanel(display),
		},
	}
}

type sheetPanel struct {
	display display
}

func newSheetPanel(display display) *sheetPanel {
	return &sheetPanel{display: display}
}

func (s *sheetPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	}
	return nocommand()
}

func (s *sheetPanel) HandleEvent(e game.Event) {
}

func (s *sheetPanel) Render(g *game.Game) {
	p := g.Player
	sheet := g.Player.Sheet

	s.display.Write(1, 2, fmt.Sprintf("NAME   %v", p.Spec.Name), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 3, fmt.Sprintf("RACE   %v", p.Spec.Species.Describe()), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 2, fmt.Sprintf("STR   %d", sheet.Stat(game.Str)), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 3, fmt.Sprintf("AGI   %d", sheet.Stat(game.Agi)), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 4, fmt.Sprintf("VIT   %d", sheet.Stat(game.Vit)), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 5, fmt.Sprintf("MND   %d", sheet.Stat(game.Mnd)), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(1, 7, fmt.Sprintf("TURN %13d", g.Progress.Turns), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 8, fmt.Sprintf("XP LEFT %10d", p.Learner.XP()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 9, fmt.Sprintf("XP TOTAL %9d", p.Learner.TotalXP()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 10, fmt.Sprintf("FLOOR %11dF", g.Progress.Floor), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 11, fmt.Sprintf("MIN FLOOR %7dF", 1), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(1, 12, fmt.Sprintf("MAX FLOOR %7dF", g.Progress.MaxFloor), termbox.ColorWhite, termbox.ColorBlack)

	s.display.Write(22, 7, fmt.Sprintf("FIGHT %9s", sheet.Attack().Describe()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 8, fmt.Sprintf("DEF %11s", sheet.Defense().Describe()), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 10, fmt.Sprintf("HP %12s", fmt.Sprintf("%d:%d", sheet.HP(), sheet.MaxHP())), termbox.ColorWhite, termbox.ColorBlack)
	s.display.Write(22, 11, fmt.Sprintf("MP %12s", fmt.Sprintf("%d:%d", sheet.MP(), sheet.MaxMP())), termbox.ColorWhite, termbox.ColorBlack)

	for sk := game.Melee; sk < game.NumSkills; sk++ {
		rowfmt := "%-5s %3d = %2d %3s"
		row := fmt.Sprintf(rowfmt, skillname(sk), sheet.Skill(sk), sheet.UnmodSkill(sk), extrasign(sheet.SkillMod(sk)))
		s.display.Write(40, 7+int(sk), row, termbox.ColorWhite, termbox.ColorBlack)
	}
}

type skillEditPanel struct {
	display display
	editing bool
	cur     game.SkillName
	change  *game.SkillChange
}

func newSkillEditPanel(display display) *skillEditPanel {
	return &skillEditPanel{display: display}
}

func (s *skillEditPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if !s.editing {
		switch tboxev.Ch {
		case 'i':
			s.editing = true
			return game.StartLearningCommand{}, nil
		}
	} else {
		switch tboxev.Key {
		case termbox.KeyEsc:
			s.editing = false
			s.change = nil
			return game.CancelLearningCommand{}, nil
		case termbox.KeyEnter:
			s.editing = false
			s.change = nil
			return game.FinishLearningCommand{}, nil
		}

		switch tboxev.Ch {
		case 'j':
			s.cur = (s.cur + 1) % game.NumSkills
		case 'k':
			if s.cur == game.Melee {
				s.cur = game.Song
			} else {
				s.cur--
			}
		case 'h':
			return game.UnlearnSkillCommand{Skill: s.cur}, nil
		case 'l':
			return game.LearnSkillCommand{Skill: s.cur}, nil
		}
	}
	return nocommand()
}

func (s *skillEditPanel) HandleEvent(e game.Event) {
	switch ev := e.(type) {
	case game.SkillChangeEvent:
		s.change = ev.Change
	}
}

func (s *skillEditPanel) Render(g *game.Game) {
	if !s.editing || s.change == nil {
		return
	}

	for sk := game.Melee; sk < game.NumSkills; sk++ {
		row, fg := fmt.Sprintf("%5d", s.change.Changes[sk].Cost), termbox.ColorWhite
		if sk == s.cur {
			fg = termbox.ColorBlue
		}
		s.display.Write(60, 7+int(sk), row, fg, termbox.ColorBlack)
	}
}

func extrasign(x int) string {
	s := fmt.Sprintf("%d", x)
	if x >= 0 {
		s = "+" + s
	}
	return s
}

func skillname(s game.SkillName) string {
	skillnames := []string{
		"FIGHT",
		"DODGE",
		"SHOOT",
		"SNEAK",
		"CHI",
		"SENSE",
		"MAGIC",
		"SONG",
	}
	return skillnames[s]
}
