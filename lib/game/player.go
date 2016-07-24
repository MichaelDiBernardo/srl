package game

// Creates a new player.
func newPlayer(eq *EventQueue) *Obj {
	return newObj(OTActor, "Player", &Traits{Mover: NewActorMover}, eq)
}
