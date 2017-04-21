package console

import (
	"container/list"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"github.com/nsf/termbox-go"
)

var hudBounds = consoleBounds

// Where the message panel should go.
var messagePanelBounds = math.Rect(math.Pt(1, 19), math.Pt(1, hudBounds.Max.Y))

// How many lines should show by default in the message panel.
var messagePanelNumLines = messagePanelBounds.Height()

// Where the status panel should go.
var statusPanelBounds = math.Rect(math.Pt(38, 0), math.Pt(hudBounds.Max.X, messagePanelBounds.Min.Y))

// Where the map panel should go.
var mapPanelBounds = math.Rect(math.Origin, math.Pt(statusPanelBounds.Min.X, messagePanelBounds.Min.Y))

// The screen coordinate where the player will render.
var mapPlayerPos = mapPanelBounds.Center()

// Create a new HUD.
func newHudScreen(display display) *screen {
	// Since the messagePanel may capture --more-- prompts and require a redraw
	// of itself in mid-screen render, we want to make sure the other panels
	// have already been drawn first.
	return &screen{
		display: display,
		panels: []panel{
			newHudControlPanel(),
			newMapPanel(display),
			newStatusPanel(display),
			newMessagePanel(messagePanelNumLines, display),
		},
	}
}

// Create a new targeting hud (used for picking ranged, spell targets.)
func newTargetScreen(display display) *screen {
	return &screen{
		display: display,
		panels: []panel{
			newTargetPanel(display),
			newMapPanel(display),
			newStatusPanel(display),
			newMessagePanel(messagePanelNumLines, display),
		},
	}
}

var hudKeymap = map[rune]game.Command{
	'h': game.MoveCommand{Dir: math.Pt(-1, 0)},
	'j': game.MoveCommand{Dir: math.Pt(0, 1)},
	'k': game.MoveCommand{Dir: math.Pt(0, -1)},
	'l': game.MoveCommand{Dir: math.Pt(1, 0)},
	'y': game.MoveCommand{Dir: math.Pt(-1, -1)},
	'u': game.MoveCommand{Dir: math.Pt(1, -1)},
	'b': game.MoveCommand{Dir: math.Pt(-1, 1)},
	'n': game.MoveCommand{Dir: math.Pt(1, 1)},
	'z': game.RestCommand{},
	'q': game.QuitCommand{},
	',': game.TryPickupCommand{},
	'd': game.TryDropCommand{},
	'w': game.TryEquipCommand{},
	'r': game.TryRemoveCommand{},
	'a': game.TryUseCommand{},
	'f': game.TryShootCommand{},
	'i': game.ModeCommand{Mode: game.ModeInventory},
	'>': game.AscendCommand{},
	'<': game.DescendCommand{},
	'@': game.ModeCommand{Mode: game.ModeSheet},
}

// A panel that just controls input to the HUD. Doesn't draw anything.
type hudControlPanel struct{}

func newHudControlPanel() *hudControlPanel {
	return &hudControlPanel{}
}

func (h *hudControlPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type != termbox.EventKey || tboxev.Key != 0 {
		return nocommand()
	}

	srlev := hudKeymap[tboxev.Ch]
	if srlev != 0 {
		return srlev, nil
	}

	return nocommand()
}

// Listens to nothing.
func (_ *hudControlPanel) HandleEvent(e game.Event) {
}

// Draw nothing.
func (_ *hudControlPanel) Render(g *game.Game) {
}

// A panel that handles targetting on the HUD, when the game is in targetting mode.
type targetPanel struct {
	display display
	targets []game.Target // The targets in LOS that we can cycle through.
	cur     int           // Index into targets for cur. If we're freetargeting, this will be -1.
	pos     math.Point    // Where our target cursor should be.
}

func newTargetPanel(display display) *targetPanel {
	return &targetPanel{display: display}
}

func (t *targetPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		// TODO: This is another good argument for a panel Init()
		t.clearTargets()
		return game.ModeCommand{Mode: game.ModeHud}, nil
	}

	switch tboxev.Ch {
	case '>':
		t.nextTarget()
	case '<':
		t.prevTarget()
	}

	return nocommand()
}

// Listens to nothing.
func (_ *targetPanel) HandleEvent(e game.Event) {
}

// Move cursor to current target.
func (t *targetPanel) Render(g *game.Game) {
	// TODO: There should be either something that comes back in the ModeEvent
	// or an explicit Init() from screen -> panel that lets this panel init
	// itself with targets on modeswitch, instead of this hack init check every
	// render.
	if t.targets == nil {
		t.initTargets(g)
	}

	if t.cur == -1 {
		t.display.SetCursor(t.pos.X, t.pos.Y)
	} else {
		pos := t.target().Pos.Add(mapPlayerPos)
		t.display.SetCursor(pos.X, pos.Y)
	}
}

func (t *targetPanel) initTargets(g *game.Game) {
	t.targets = g.Player.Shooter.Targets()

	if t.anyTargets() {
		t.cur = 0
	} else {
		t.cur = -1
		t.pos = mapPlayerPos
	}
}

func (t *targetPanel) clearTargets() {
	t.targets = nil
}

func (t *targetPanel) prevTarget() {
	if !t.anyTargets() {
		return
	}

	t.cur--
	if t.cur < 0 {
		t.cur = len(t.targets) - 1
	}
}

func (t *targetPanel) nextTarget() {
	if !t.anyTargets() {
		return
	}
	t.cur = (t.cur + 1) % len(t.targets)
}

func (t *targetPanel) anyTargets() bool {
	return len(t.targets) > 0
}

func (t *targetPanel) target() *game.Target {
	if t.cur == -1 {
		return nil
	} else {
		return &t.targets[t.cur]
	}
}

// Panel that renders the gameplay map.
type mapPanel struct {
	display display
}

// A glyph used to render a tile.
type glyph struct {
	Ch rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

// Glyphs used to render actors.
var actorGlyphs = map[game.Species]glyph{
	game.SpecHuman: glyph{Ch: '@', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
	game.SpecOrc:   glyph{Ch: 'o', Fg: termbox.ColorGreen, Bg: termbox.ColorBlack},
	game.SpecAnt:   glyph{Ch: 'd', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
}

// Glyphs used to render items.
var itemGlyphs = map[game.Species]glyph{
	game.SpecSword:        glyph{Ch: '|', Fg: termbox.ColorBlue, Bg: termbox.ColorBlack},
	game.SpecColt:         glyph{Ch: '{', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	game.SpecLeatherArmor: glyph{Ch: '[', Fg: termbox.ColorYellow, Bg: termbox.ColorBlack},
	game.SpecCure:         glyph{Ch: '!', Fg: termbox.ColorGreen, Bg: termbox.ColorBlack},
	game.SpecStim:         glyph{Ch: '!', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	game.SpecHyper:        glyph{Ch: '!', Fg: termbox.ColorYellow, Bg: termbox.ColorBlack},
	game.SpecRestore:      glyph{Ch: '!', Fg: termbox.ColorBlue, Bg: termbox.ColorBlack},
}

// Glyphs used to render tiles.
var featureGlyphs = map[game.FeatureType]glyph{
	"FeatWall":       glyph{Ch: '#', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	"FeatFloor":      glyph{Ch: '.', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
	"FeatClosedDoor": glyph{Ch: '+', Fg: termbox.ColorYellow, Bg: termbox.ColorBlack},
	"FeatOpenDoor":   glyph{Ch: '\'', Fg: termbox.ColorYellow, Bg: termbox.ColorBlack},
	"FeatStairsUp":   glyph{Ch: '>', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
	"FeatStairsDown": glyph{Ch: '<', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
}

// Create a new mapPanel.
func newMapPanel(display display) *mapPanel {
	return &mapPanel{display: display}
}

func (m *mapPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	return nocommand()
}

// Listens to nothing.
func (m *mapPanel) HandleEvent(e game.Event) {
}

// Render the gameplay map to the hud.
func (m *mapPanel) Render(g *game.Game) {
	center := g.Player.Pos()
	boundsdist := mapPlayerPos
	viewport := math.Rect(center.Sub(boundsdist), center.Add(boundsdist))
	maptrans := mapPanelBounds.Min.Sub(viewport.Min)
	level := g.Level

	for x := viewport.Min.X; x < viewport.Max.X; x++ {
		for y := viewport.Min.Y; y < viewport.Max.Y; y++ {
			cur := math.Pt(x, y)
			if !cur.In(level) {
				continue
			}

			tile := level.At(cur)
			hasactor := tile.Actor != nil
			isplayer := hasactor && tile.Actor.IsPlayer()
			drawpos := cur.Add(maptrans)

			// When you're blind, you may be walking on unseen tiles. So, we
			// always want to show the player, even if the tile is unseen.
			if !tile.Seen && !isplayer {
				m.display.SetCell(drawpos.X, drawpos.Y, ' ', termbox.ColorBlack, termbox.ColorBlack)
				continue
			}

			var gl glyph

			// If you're blind (see above), you may be on an unseen and/or
			// not-visible tile, but we still want to draw the player glyph.
			if hasactor && (tile.Visible || isplayer) {
				gl = actorGlyphs[tile.Actor.Spec.Species]
				if tile.Actor.Sheet.Petrified() {
					gl.Fg = termbox.ColorBlack | termbox.AttrBold
				}
			} else if !tile.Items.Empty() {
				item, stack := tile.Items.Top(), tile.Items.Len() > 1
				gl = itemGlyphs[item.Spec.Species]
				if stack {
					gl.Bg = termbox.ColorCyan
				}
			} else {
				gl = featureGlyphs[tile.Feature.Type]
				if !tile.Visible {
					gl.Fg = termbox.ColorBlack | termbox.AttrBold
				}
			}
			m.display.SetCell(drawpos.X, drawpos.Y, gl.Ch, gl.Fg, gl.Bg)
		}
	}
}

type messageLine struct {
	text string
}

type morePrompt struct {
	acked bool
}

// The message panel, where messages are rendered at the bottom of the hud.
type messagePanel struct {
	display display
	lines   *list.List
	size    int
}

// Create a new messagePanel.
func newMessagePanel(size int, display display) *messagePanel {
	return &messagePanel{
		display: display,
		lines:   list.New(),
		size:    size,
	}
}

func (m *messagePanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	return nocommand()
}

// Listens for new message events to build the message list.
func (m *messagePanel) HandleEvent(e game.Event) {
	switch ev := e.(type) {
	case game.MessageEvent:
		m.lines.PushBack(&messageLine{text: ev.Text})
	case game.MoreEvent:
		m.lines.PushBack(&morePrompt{acked: false})
	}

	if m.lines.Len() > m.size {
		m.lines.Remove(m.lines.Front())
	}
}

// Render the panel to the display.
func (m *messagePanel) Render(g *game.Game) {
	more := false
	var e *list.Element
	i := 0

	for e = m.lines.Front(); e != nil; e = e.Next() {
		switch line := e.Value.(type) {
		case *messageLine:
			m.display.Write(messagePanelBounds.Min.X, messagePanelBounds.Min.Y+i, line.text, termbox.ColorWhite, termbox.ColorBlack)
			i++
		case *morePrompt:
			if line.acked {
				continue
			}
			line.acked = true
			more = true
			break
		}
	}

	if more {
		m.display.Write(messagePanelBounds.Min.X, messagePanelBounds.Min.Y+i, "--more--", termbox.ColorWhite, termbox.ColorBlack)
		m.display.Flush()

		// Wait for player to clear prompt.
		for {
			tboxev := m.display.PollEvent()
			if tboxev.Type == termbox.EventKey && (tboxev.Key == termbox.KeyEsc || tboxev.Key == termbox.KeyEnter) {
				break
			}
		}

		// HACK: Clear the --more-- line.
		m.display.Write(messagePanelBounds.Min.X, messagePanelBounds.Min.Y+i, "        ", termbox.ColorWhite, termbox.ColorBlack)

		// Rerender the message panel.
		m.Render(g)
		m.display.Flush()
	}
}

// Panel that renders player status on hud.
type statusPanel struct {
	display display
}

// Create a new statusPanel.
func newStatusPanel(display display) *statusPanel {
	return &statusPanel{display: display}
}

func (s *statusPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	return nocommand()
}

// Listens for nothing.
func (s *statusPanel) HandleEvent(e game.Event) {
}

// Render the panel to the display.
func (s *statusPanel) Render(g *game.Game) {
	player := g.Player
	fg, bg := termbox.ColorWhite, termbox.ColorBlack
	sheet := player.Sheet

	lcol, rcol := statusPanelBounds.Min.X, statusPanelBounds.Min.X+20

	// left
	s.display.Write(lcol, statusPanelBounds.Min.Y+0, player.Spec.Name, fg, bg)
	s.display.Write(lcol, statusPanelBounds.Min.Y+1, "Human", fg, bg)

	s.display.Write(lcol, statusPanelBounds.Min.Y+3, fmt.Sprintf("%-7s%3d/%-3d", "HP", sheet.HP(), sheet.MaxHP()), fg, bg)
	s.display.Write(lcol, statusPanelBounds.Min.Y+4, fmt.Sprintf("%-7s%3d/%-3d", "MP", sheet.MP(), sheet.MaxMP()), fg, bg)
	s.display.Write(lcol, statusPanelBounds.Min.Y+5, fmt.Sprintf("%-7s%5d", "XP", player.Learner.XP()), fg, bg)
	s.display.Write(lcol, statusPanelBounds.Min.Y+6, fmt.Sprintf("%-7s%2dF", "FL", g.Progress.Floor), fg, bg)

	s.display.Write(lcol, statusPanelBounds.Min.Y+8, fmt.Sprintf("%-7s%8s", "FIGHT", sheet.Attack().Describe()), fg, bg)

	// right
	s.display.Write(rcol, statusPanelBounds.Min.Y+3, fmt.Sprintf("%-7s%3d", "STR", sheet.Stat(game.Str)), fg, bg)
	s.display.Write(rcol, statusPanelBounds.Min.Y+4, fmt.Sprintf("%-7s%3d", "AGI", sheet.Stat(game.Agi)), fg, bg)
	s.display.Write(rcol, statusPanelBounds.Min.Y+5, fmt.Sprintf("%-7s%3d", "VIT", sheet.Stat(game.Vit)), fg, bg)
	s.display.Write(rcol, statusPanelBounds.Min.Y+6, fmt.Sprintf("%-7s%3d", "MND", sheet.Stat(game.Mnd)), fg, bg)

	def, defc := sheet.Defense(), fg
	if len(def.CorrDice) > 0 {
		defc = termbox.ColorRed
	}
	s.display.Write(rcol, statusPanelBounds.Min.Y+8, fmt.Sprintf("%-7s%8s", "DEF", def.Describe()), defc, bg)

	s.printEffects(player)
}

func (s *statusPanel) printEffects(player *game.Obj) {
	x, y := statusPanelBounds.Min.X, statusPanelBounds.Min.Y+10
	linelen, w := 0, statusPanelBounds.Width()

	for _, lightspec := range effectLights {
		ticks := lightspec.cond(player)
		if ticks <= 0 {
			continue
		}
		light, col := lightspec.makelight(lightspec, ticks)

		incr := len(light) + 1
		if linelen+incr > w {
			x, y, linelen = statusPanelBounds.Min.X, y+1, 0
		}

		s.display.Write(x, y, light, col, termbox.ColorBlack)

		x += incr
		linelen += incr
	}
}

type effectLight struct {
	label        string
	defaultcolor termbox.Attribute
	cond         func(*game.Obj) int
	makelight    func(effectLight, int) (string, termbox.Attribute)
}

var effectLights = []effectLight{
	{
		label:        "Pois %2d",
		defaultcolor: termbox.ColorGreen,
		cond: func(p *game.Obj) int {
			return p.Ticker.Counter(game.EffectPoison)
		},
		makelight: makeCountingLight,
	},
	{
		label:        "Cut %2d",
		defaultcolor: termbox.ColorRed,
		cond: func(p *game.Obj) int {
			return p.Ticker.Counter(game.EffectCut)
		},
		makelight: makeCountingLight,
	},
	{
		label:        "Stun %2d",
		defaultcolor: termbox.ColorYellow,
		cond: func(p *game.Obj) int {
			return p.Ticker.Counter(game.EffectStun)
		},
		makelight: makeCountingLight,
	},
	{
		label:        "Blind",
		defaultcolor: termbox.ColorWhite,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Blind())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Slow",
		defaultcolor: termbox.ColorBlue,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Slow())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Conf",
		defaultcolor: termbox.ColorYellow,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Confused())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Fear",
		defaultcolor: termbox.ColorWhite,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Afraid())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Para",
		defaultcolor: termbox.ColorRed,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Paralyzed())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Stone",
		defaultcolor: termbox.ColorWhite,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Petrified())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Sil",
		defaultcolor: termbox.ColorBlue,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Silenced())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Curse",
		defaultcolor: termbox.ColorWhite,
		cond: func(p *game.Obj) int {
			return bool2int(p.Sheet.Cursed())
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Hyp",
		defaultcolor: termbox.ColorBlue,
		cond: func(p *game.Obj) int {
			return p.Ticker.Counter(game.EffectHyper)
		},
		makelight: makeLabelLight,
	},
	{
		label:        "Stim",
		defaultcolor: termbox.ColorYellow,
		cond: func(p *game.Obj) int {
			return p.Ticker.Counter(game.EffectStim)
		},
		makelight: makeLabelLight,
	},
}

func makeCountingLight(el effectLight, ticks int) (light string, fg termbox.Attribute) {
	return fmt.Sprintf(el.label, ticks), el.defaultcolor
}

func makeLabelLight(el effectLight, ticks int) (light string, fg termbox.Attribute) {
	return el.label, el.defaultcolor
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
