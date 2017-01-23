package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

const (
	FOVRadiusMax = 8
	FOVRadius    = 4
	ScentRadius  = FOVRadius
)

type Field []math.Point

// Senses all the things an actor can sense.
type Senser interface {
	Objgetter
	CalcFields()
	FOV() Field
	CanSee(other *Obj) bool
}

type ActorSenser struct {
	Trait
	fov Field
}

func NewActorSenser(obj *Obj) Senser {
	return &ActorSenser{Trait: Trait{obj: obj}}
}

// Calculates this actor's fields. If it's the player, they get FOV + scent.
// Monsters just get FOV. Instead of making two different Senser
// implementations to do this, we use simple branch in this method.  If this is
// the player, calculating fields will also update the relevant bits on the map
// (scent, visibility.)
func (a *ActorSenser) CalcFields() {
	// If we're doing sense and scent, we calculate whichever has the bigger
	// radius and then use a subset of the points for each.
	obj := a.obj
	player := obj.IsPlayer()
	sightrad, scentrad := obj.Sheet.Sight(), 0
	if player {
		scentrad = ScentRadius
	}

	// Calculate field.
	radius := math.Max(sightrad, scentrad)
	field := a.field(radius)

	// Sight.
	a.fov = trimfield(field, sightrad, radius)
	if !obj.IsPlayer() {
		return
	}

	// Scent, if you're the player. The player is smelly.
	scent := trimfield(field, scentrad, radius)

	// Update these fields on the level.
	obj.Level.UpdateVis(a.fov)
	obj.Level.UpdateScent(scent)
}

func (a *ActorSenser) FOV() Field {
	return a.fov
}

func (a *ActorSenser) CanSee(other *Obj) bool {
	pos := other.Pos()
	for _, pt := range a.fov {
		if pos == pt {
			return true
		}
	}
	return false
}

// Calculates a "field-of-vision" type "field" around the actor.
func (a *ActorSenser) field(radius int) Field {
	// TODO: This is basically a direct translation of fcrawl's raycasting FOV
	// algorithm. I didn't try at all to make it less pythony and more go-ey.
	// Should replace with something less churny or just a totally different
	// algorithm, there's enough of these written in C that should be a lot
	// less impedence-mismatchy to directly translate.
	fov := newPointSet()
	fov.Add(math.Origin)

	pos, level := a.Obj().Pos(), a.Obj().Level

	// Light begins casting in all directions.
	light := make(map[math.Point]pointset)
	light[math.Origin] = newPointSetL(math.ChebyEdge(1))

	for r := 0; r < radius; r++ {
		edge := math.ChebyEdge(r)
		for _, cpt := range edge {
			li, pt := light[cpt], pos.Add(cpt)
			if li == nil || !pt.In(level) || level.At(pt).Feature.Opaque {
				continue
			}

			for dp, _ := range li {
				cur := cpt.Add(dp)
				next := li.Intersect(adj45dirs(dp))
				light[cur] = light[cur].Union(next)
				fov.Add(cur)
			}
		}
	}

	// fov is currently relative to 0,0 as center, and has not yet been
	// translated to the map coords. We also take this opportunity to coerce
	// the set into a slice.
	transfov := make(Field, 0, len(fov))
	for p, _ := range fov {
		tpt := p.Add(pos)
		if tpt.In(level) {
			transfov = append(transfov, tpt)
		}
	}
	return transfov
}

func adj45dirs(d math.Point) pointset {
	dirscircle := []math.Point{
		math.Pt(-1, 0), math.Pt(-1, -1), math.Pt(0, -1), math.Pt(1, -1),
		math.Pt(1, 0), math.Pt(1, 1), math.Pt(0, 1), math.Pt(-1, 1),
	}

	index := -1
	for i, cpt := range dirscircle {
		if cpt == d {
			index = i
			break
		}
	}

	dirs := newPointSet()
	for i := -1; i <= 1; i++ {
		dirind := (index + i) % 8
		if dirind < 0 {
			dirind = 8 + dirind
		}
		dirs.Add(dirscircle[dirind])
	}
	return dirs
}

type pointset map[math.Point]bool

func newPointSet() pointset {
	return make(map[math.Point]bool)
}

func newPointSetL(pts []math.Point) pointset {
	ps := newPointSet()
	for _, pt := range pts {
		ps.Add(pt)
	}
	return ps
}

func (ps pointset) Add(pt math.Point) {
	ps[pt] = true
}

func (ps pointset) Union(other pointset) pointset {
	union := newPointSet()
	for k, v := range ps {
		union[k] = v
	}
	for k, v := range other {
		union[k] = v
	}
	return union
}

func (ps pointset) Intersect(other pointset) pointset {
	intersection := newPointSet()
	for k, v := range ps {
		if other[k] {
			intersection[k] = v
		}
	}
	return intersection
}

func trimfield(field Field, r int, max int) []math.Point {
	// Premature optimization -- if the radius we're asking to trim to is the
	// max, don't bother searching for a trim point.
	if r == max {
		return field
	}
	i := 0
	for _, p := range field {
		i++
		if math.ChebyDist(math.Origin, p) > r {
			break
		}
	}
	return field[:i]
}
