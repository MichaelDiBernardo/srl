package game

func NewPlayer() *Obj {
	player := NewObj(&Traits{Mover: NewActorMover})
	player.Type = OTActor
	player.Subtype = "Player"
	return player
}
