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
			Sheet: NewPlayerSheet,
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
				vit: 1,
				mnd: 0,
			}),
			Sheet: NewMonsterSheet(MonsterSheet{
				melee:   1,
				evasion: 1,
				maxhp:   20,
				maxmp:   10,
			}),
		},
	}
)
