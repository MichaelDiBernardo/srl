package game

// A thing that drops items when it dies.
type Dropper interface {
	DropItems()
}

type ItemDropper struct {
	Trait
	// Drop [1..num] items...
	num int
	// with gen depth boosted by 'boost'.
	boost int
}

func NewItemDropper(spec *ItemDropper) func(*Obj) Dropper {
	return func(o *Obj) Dropper {
		d := &ItemDropper{}
		*d = *spec
		d.obj = o
		return d
	}
}

func (i *ItemDropper) DropItems() {
	if i.num == 0 {
		return
	}

	g := i.obj.Game
	num := RandInt(0, i.num) + 1

	groups := Generate(num, g.Floor, 2, Items, g)

	for _, group := range groups {
		for _, item := range group {
			g.Level.Place(item, i.obj.Pos())
		}
	}
}
