package game

func NewPlayer() *Obj {
	return NewObj(OTActor, "Player", &Traits{Mover: NewActorMover})
}
