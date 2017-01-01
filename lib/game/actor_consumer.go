package game

// A thing that can use single-use things.
type Consumer interface {
	Objgetter
	Use(consumable Consume)
}
