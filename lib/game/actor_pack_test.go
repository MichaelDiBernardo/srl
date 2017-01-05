package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

func TestTryPickupNoItemsOnGround(t *testing.T) {
	g := newTestGame()
	taker := g.NewObj(atActorSpec)
	g.Level.Place(taker, math.Pt(1, 1))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size > 0 {
		t.Errorf(`TryPickup() on empty square gave inven size %d; want 0`, size)
	}
}

func TestTryPickupOneItemOnGround(t *testing.T) {
	g := newTestGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)

	g.Level.Place(taker, math.Pt(1, 1))
	g.Level.Place(item, math.Pt(1, 1))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size != 1 {
		t.Errorf(`TryPickup() on 1-item square gave inven size %d; want 1`, size)
	}
	if size := g.Level.At(math.Pt(1, 1)).Items.Len(); size != 0 {
		t.Errorf(`TryPickup() on 1-item square left %d items on ground; want 0`, size)
	}
}

func TestTryPickupFromStack(t *testing.T) {
	g := newTestGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)
	item2 := g.NewObj(atItemSpec)

	g.Level.Place(taker, math.Pt(1, 1))
	g.Level.Place(item, math.Pt(1, 1))
	g.Level.Place(item2, math.Pt(1, 1))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size != 0 {
		t.Errorf(`TryPickup() on stack took something instead of opening menu; took %d things`, size)
	}
	if size := g.Level.At(math.Pt(1, 1)).Items.Len(); size != 2 {
		t.Errorf(`TryPickup() took from ground instead of opening menu; left %d things`, size)
	}
	if size := g.Events.Len(); size != 1 {
		t.Errorf(`TryPickup() pushed wrong # of events to queue; found %d, want 1`, size)
	}

	e, ok := g.Events.Next().(ModeEvent)
	if !ok {
		t.Error(`TryPickup pushed wrong event type to queue.`)
	}
	if e.Mode != ModePickup {
		t.Errorf(`TryPickup switched to mode %v, want %v`, e.Mode, ModePickup)
	}

}

func TestPickupOutOfBounds(t *testing.T) {
	g := newTestGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)

	g.Level.Place(taker, math.Pt(1, 1))
	g.Level.Place(item, math.Pt(1, 1))

	taker.Packer.TryPickup()
	taker.Packer.Pickup(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Out-of-bounds Pickup switched to mode %v; want %v`, mode, ModeHud)
	}

}

func TestTryDropWithNothingInInventory(t *testing.T) {
	g := newTestGame()

	packer := g.NewObj(atActorSpec)
	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryDrop w no items switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryDropWithFullGround(t *testing.T) {
	g := newTestGame()

	packer := g.NewObj(atActorSpec)
	g.Level.Place(packer, math.Pt(1, 1))

	packer.Tile.Items.capacity = 0
	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryDrop w full ground switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryDrop(t *testing.T) {
	g := newTestGame()

	packer := g.NewObj(atActorSpec)
	g.Level.Place(packer, math.Pt(1, 1))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeDrop {
		t.Errorf(`TryDrop switched to mode %v, want %v`, mode, ModeDrop)
	}
}

func TestDrop(t *testing.T) {
	g := newTestGame()

	packer := g.NewObj(atActorSpec)
	g.Level.Place(packer, math.Pt(1, 1))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()
	packer.Packer.Drop(0)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Dropping switched mode to %v, want %v`, mode, ModeHud)
	}

	if dropped := packer.Tile.Items.Top(); dropped != item {
		t.Errorf(`Dropped item was %v; want %v`, dropped, item)
	}
}

func TestDropOutOfBounds(t *testing.T) {
	g := newTestGame()

	packer := g.NewObj(atActorSpec)
	g.Level.Place(packer, math.Pt(1, 1))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()
	packer.Packer.Drop(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Dropping switched mode to %v, want %v`, mode, ModeHud)
	}
}
