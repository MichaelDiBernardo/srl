package game

// Creates a new player.
func newPlayer(eq *EventQueue) *Obj {
    spec := &ActorSpec{
        Type: "Player",
        // TODO: ActorSpec should just be a spec, and all objs can just point
        // to them instead of having to copy everything over.
        Name: "DEBO",
		Traits: &Traits{
			Mover: NewActorMover,
		},
    }
	return newActor(spec, eq)
}
