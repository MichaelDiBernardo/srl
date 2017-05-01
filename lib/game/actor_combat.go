package game

import (
	"errors"
	"fmt"
	"log"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
}

// Melee combat.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *ActorFighter) Hit(other Fighter) {
	atk := f.obj.Sheet.Attack()
	hit(atk, f.obj, other.Obj())
}

// Anything that can attack at range.
type Shooter interface {
	Objgetter
	// Primarily for the player -- try to switch to "shooting mode". If the
	// player has no ranged weapon or is otherwise in no shape to shoot, this
	// will emit a message and cancel the modeswitch.
	TryShoot()
	// Return a list of targets in LOS, sorted by proximity.
	Targets() []Target
	// Given a point on the map, this will give the target info for shooting at
	// that spot. Will return ErrTargetOutOfRange or ErrNoClearShot if the shot
	// is impossible.
	Target(math.Point) (Target, error)
	// Shoot at the given point on the map. Returns the same errors under the
	// same conditions that Target() does. If an actor is hit on the way to the
	// intended target, combat rolls etc. will be run against it and the
	// projectile will stop.
	Shoot(math.Point) error
}

var ErrTargetOutOfRange = errors.New("TargetOutOfRange")
var ErrNoClearShot = errors.New("NoClearShot")

// A target in LOS of the shooter. Points in Pos and Path are relative to the
// actor, who is considered to be at origin. So, if there is a target to the
// east with one space intervening, its Pos would be (2,0) and its path would
// be {(1,0), (2,0))
type Target struct {
	Pos    math.Point
	Path   Path
	Target *Obj
}

type ActorShooter struct {
	Trait
}

func NewActorShooter(obj *Obj) Shooter {
	return &ActorShooter{Trait: Trait{obj: obj}}
}

func (s *ActorShooter) TryShoot() {
	obj := s.obj
	if obj.Equipper.Body().Shooter() == nil {
		obj.Game.Events.Message("Nothing to shoot with.")
	} else if s.obj.Sheet.Afraid() {
		msg := fmt.Sprintf("%s is too afraid to shoot!", obj.Spec.Name)
		obj.Game.Events.Message(msg)
	} else if s.obj.Sheet.Confused() {
		msg := fmt.Sprintf("%s is too confused to shoot!", obj.Spec.Name)
		obj.Game.Events.Message(msg)
	} else {
		obj.Game.SwitchMode(ModeShoot)
	}
}

func (s *ActorShooter) Targets() []Target {
	fov := s.obj.Senser.FOV()
	targets := []Target{}

	for _, p := range fov {
		tile := s.obj.Game.Level.At(p)
		victim := tile.Actor

		if victim == nil || victim == s.obj {
			continue
		}

		target, err := s.Target(tile.Pos)
		if err != nil {
			continue
		}
		targets = append(targets, target)
	}
	return targets
}

func (s *ActorShooter) Target(p math.Point) (Target, error) {
	mypos, lev := s.obj.Pos(), s.obj.Game.Level
	srange := s.obj.Sheet.Ranged().Range

	if math.EucDist(s.obj.Pos(), p) > srange {
		return Target{}, ErrTargetOutOfRange
	}

	path, ok := lev.FindLine(mypos, p, LineTest)
	if !ok {
		return Target{}, ErrNoClearShot
	}

	target := Target{
		Pos:    p,
		Path:   path,
		Target: lev.At(p).Actor,
	}
	return target, nil
}

func (s *ActorShooter) Shoot(p math.Point) error {
	target, err := s.Target(p)
	if err != nil {
		return err
	}

	attacker, level, path := s.obj, s.obj.Level, target.Path[1:]
	atk := attacker.Sheet.Ranged()
	dist := 0

	for _, pt := range path {
		dist++
		defender := level.At(pt).Actor

		atk.Hit -= dist / 2

		// Unintended targets are much harder to hit.
		if pt != target.Pos {
			atk.Hit /= 2
		}

		if hit(atk, attacker, defender) {
			break
		}
	}

	return nil
}

// hit has the attacker attack the defender with the given attack. This works
// for both melee and ranged attacks, but presumes the target is in range.
// Returns true if the target was hit; it could be for 0 damage, but as long as
// contact was made, this will return true.
func hit(atk Attack, a, d *Obj) bool {
	def := d.Sheet.Defense()

	atkroll := combatroll(a) + atk.Hit
	defroll := combatroll(d) + def.Evasion
	residual := atkroll - defroll

	aname, dname := a.Spec.Name, d.Spec.Name

	if residual <= 0 {
		msg := fmt.Sprintf("%v missed %v.", aname, dname)
		a.Game.Events.Message(msg)
		return false
	}

	crits := residual / (atk.CritDiv + def.Effects.Has(ResistCrit))

	// Calculate raw phys damage.
	droll, proll := atk.RollDamage(crits), def.RollProt()
	log.Printf("DR: %d PR: %d", droll, proll)
	dmg := math.Max(0, droll-proll)

	// Figure out how much branded damage we did.
	xdmg, poisondmg := applybs(dmg, atk.Effects, def.Effects)
	dmg += xdmg

	critstr := ""
	if crits > 0 {
		critstr = fmt.Sprintf(" %dx critical!", crits)
	}

	msg := fmt.Sprintf("%s %s %s (%d).%s", aname, atk.Verb, dname, dmg, critstr)
	a.Game.Events.Message(msg)

	if dmg <= 0 {
		return true
	}

	ispara := d.Sheet.Paralyzed()

	for effect, _ := range atk.Effects {
		switch effect {
		case BrandPoison:
			d.Ticker.AddEffect(EffectPoison, poisondmg)
		case BrandAcid:
			if OneIn(def.Effects.Resists(effect) + 1) {
				d.Ticker.AddEffect(EffectShatter, DieRoll(4, 4))
			}
		case EffectStun:
			score := atk.CritDiv - BaseCritDiv + a.Sheet.Stat(Str)
			difficulty := d.Sheet.Skill(Chi)
			resists := def.Effects.Resists(effect)
			won, _ := skillcheck(score, difficulty, resists, a, d)
			if won {
				d.Ticker.AddEffect(EffectStun, dmg)
			}
		case EffectBlind:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectBlind, DieRoll(5, 4))
			}
		case EffectConfuse:
			// TODO: Eventually remove this check and instead use a Cruel-Blow
			// style check, cruel blow should be the only ability that gives
			// confusion melee anyways.
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectConfuse, DieRoll(5, 4))
			}
		case EffectPara:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectPara, DieRoll(4, 4))
			}
		case EffectPetrify:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectPetrify, DieRoll(4, 4))
			}
		case EffectCut:
			if crits > DieRoll(1, 2) {
				d.Ticker.AddEffect(EffectCut, dmg/2)
			}
		case EffectShatter:
			score := atk.CritDiv - BaseCritDiv + a.Sheet.Stat(Str)
			won, _ := skillcheck(score, 10, 0, a, d)
			if won {
				d.Ticker.AddEffect(EffectShatter, DieRoll(4, 4))
			}
		case EffectDrainStr:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainStr, 1)
			}
		case EffectDrainAgi:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainAgi, 1)
			}
		case EffectDrainVit:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainVit, 1)
			}
		case EffectDrainMnd:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainMnd, 1)
			}
		}
	}

	if ispara {
		checkpara(d)
	}
	d.Sheet.Hurt(dmg)

	// Have to handle this outside the effect loop above because we need the
	// damage to be applied to the target first to determine if they're dead.
	if d.Sheet.Dead() &&
		atk.Effects.Has(EffectVamp) > 0 &&
		def.Effects.Resists(EffectVamp) <= 0 {
		vamp(a, d)
	}
	return true
}

// Given the base physical damage done by an attack, and the atk and def
// effects, this figures out how much extra and poison damage should be done from
// brands and slays. Poison damage is separated out because it is applied as
// damage-over-time, instead of being immediately inflicted on the target.
func applybs(basedmg int, atk Effects, def Effects) (xdmg, poisondmg int) {
	if basedmg == 0 {
		return 0, 0
	}

	slays, brands := atk.Slays(), atk.Brands()
	xdmg, poisondmg = 0, 0

	for slay, _ := range slays {
		if def.SlainBy(slay) <= 0 {
			return
		}
		xdmg += DieRoll(1, basedmg)
	}

	for brand, _ := range brands {
		raw := DieRoll(1, basedmg)
		resisted := def.ResistDmg(brand, raw)

		if brand == BrandPoison {
			poisondmg += resisted
		} else {
			xdmg += resisted
		}
	}
	return xdmg, poisondmg
}

func checkpara(obj *Obj) {
	sheet := obj.Sheet

	if OneIn(2) {
		msg := fmt.Sprintf("%s breaks out of paralysis!", obj.Spec.Name)
		obj.Game.Events.Message(msg)

		sheet.SetParalyzed(false)
		obj.Ticker.RemoveEffect(EffectPara)
	}
}

func vamp(attacker, defender *Obj) {
	a, d := attacker.Sheet, defender.Sheet
	vhp, vmp := d.MaxHP()/4, d.MaxMP()/4
	a.Heal(vhp)
	a.HealMP(vmp)
}
