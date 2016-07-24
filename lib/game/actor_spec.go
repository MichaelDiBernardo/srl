package game

var (
	MonOrc = &ActorSpec{
		Type: "MonOrc",
		Traits: &Traits{
			Mover: NewActorMover,
			AI:    NewRandomAI,
		},
	}
)
