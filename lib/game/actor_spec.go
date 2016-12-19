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
			Mover:    NewActorMover,
			Fighter:  NewPlayerFighter,
			Packer:   NewActorPacker,
			Equipper: NewActorEquipper,
			Sheet:    NewPlayerSheet,
		},
	}
	MonOrc = &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits: &Traits{
			Mover:    NewActorMover,
			AI:       NewRandomAI,
			Fighter:  NewPlayerFighter,
			Packer:   NewActorPacker,
			Equipper: NewActorEquipper, // TEMP
			Sheet: NewMonsterSheet(MonsterSheet{
				str: 2,
				agi: 0,
				vit: 1,
				mnd: 0,

				melee:   1,
				evasion: 1,
				maxhp:   20,
				maxmp:   10,
			}),
		},
	}
)
