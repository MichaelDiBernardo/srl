package game

var (
	PlayerSpec = &Spec{
		Type:    OTActor,
		Subtype: OSTPlayer,
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
		Type:    OTActor,
		Subtype: OSTMonster,
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
