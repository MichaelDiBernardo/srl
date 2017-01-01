package game

// Stuff used in actor_*_test.go.
var (
	atActorSpec = &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Name:    "Hi",
		Traits: &Traits{
			Mover:    NewActorMover,
			Packer:   NewActorPacker,
			Equipper: NewActorEquipper,
		},
	}

	atItemSpec = &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: "testspec",
		Name:    "Item",
		Traits: &Traits{
			Equip: NewEquip(Equip{Slot: SlotHand}),
		},
	}
)
