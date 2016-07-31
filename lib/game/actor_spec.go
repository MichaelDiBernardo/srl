package game

var (
	PlayerSpec = &Spec{
		Type:    OTActor,
		Subtype: "Player",
		Name:    "DEBO",
		Traits: &Traits{
			Mover: NewActorMover,
			Stats: NewActorStats(stats{
				str: 2,
				agi: 2,
				vit: 2,
				mnd: 2,
			}),
		},
	}
	MonOrc = &Spec{
		Type:    OTActor,
		Subtype: "MonOrc",
		Name:    "ORC",
		Traits: &Traits{
			Mover: NewActorMover,
			AI:    NewRandomAI,
			Stats: NewActorStats(stats{
				str: 2,
				agi: 0,
				vit: 2,
				mnd: 0,
			}),
		},
	}
)
