package game

const (
	// Player species.
	SpecHuman = "human"

	// Monster species.
	SpecOrc = "orc"
)

var (
	PlayerSpec = &Spec{
		Family:  FamActor,
		Genus:   GenPlayer,
		Species: SpecHuman,
		Name:    "DEBO",
		Traits: &Traits{
			Mover:   NewActorMover,
			Fighter: NewActorFighter,
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
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits: &Traits{
			Mover:   NewActorMover,
			AI:      NewRandomAI,
			Fighter: NewActorFighter,
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
