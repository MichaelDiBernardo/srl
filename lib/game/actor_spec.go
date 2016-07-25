package game

var (
	MonOrc = &ActorSpec{
		Type: "MonOrc",
        Name: "ORC",
		Traits: &Traits{
			Mover: NewActorMover,
			AI:    NewRandomAI,
		},
	}
)
