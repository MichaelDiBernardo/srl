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
			User:     NewActorUser,
			Sheet:    NewPlayerSheet,
		},
	}

	atItemSpec = &Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: "testspec",
		Name:    "Item",
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{Slot: SlotHand}),
		},
	}
)
