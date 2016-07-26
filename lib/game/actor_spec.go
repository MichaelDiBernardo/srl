package game

var (
	PlayerSpec = &Spec{
		Type:    OTActor,
		Subtype: "Player",
		Name:    "DEBO",
		Traits: &Traits{
			Mover: NewActorMover,
		},
	}
	MonOrc = &Spec{
		Type:    OTActor,
		Subtype: "MonOrc",
		Name:    "ORC",
		Traits: &Traits{
			Mover: NewActorMover,
			AI:    NewRandomAI,
		},
	}
)
