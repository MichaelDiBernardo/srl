package game

const (
	// Player species.
	SpecHuman = "human"

	// Monster species.
	SpecOrc = "orc"
	SpecAnt = "ant"
)

var PlayerSpec = &Spec{
	Family:  FamActor,
	Genus:   GenPlayer,
	Species: SpecHuman,
	Name:    "DEBO",
	Traits: &Traits{
		Mover:    NewActorMover,
		Fighter:  NewActorFighter,
		Packer:   NewActorPacker,
		Equipper: NewActorEquipper,
		User:     NewActorUser,
		Sheet:    NewPlayerSheet,
	},
}

var Monsters = []*Spec{
	&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Gen: Gen{
			Depths:    []int{1},
			GroupSize: 2,
		},
		Traits: &Traits{
			Mover:   NewActorMover,
			AI:      NewRandomAI,
			Fighter: NewActorFighter,
			Packer:  NewActorPacker,
			Sheet: NewMonsterSheet(MonsterSheet{
				str: 2,
				agi: 0,
				vit: 1,
				mnd: 0,

				melee:      1,
				evasion:    1,
				critdivmod: 4,
				maxhp:      20,
				maxmp:      10,

				damroll:  NewDice(2, 7),
				protroll: NewDice(1, 4),
			}),
		},
	},
	&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecAnt,
		Name:    "ANT",
		Gen: Gen{
			Depths:    []int{1},
			GroupSize: 2,
		},
		Traits: &Traits{
			Mover:   NewActorMover,
			AI:      NewRandomAI,
			Fighter: NewActorFighter,
			Packer:  NewActorPacker,
			Sheet: NewMonsterSheet(MonsterSheet{
				str: 3,
				agi: 3,
				vit: 1,
				mnd: 0,

				melee:      5,
				evasion:    5,
				critdivmod: 4,
				maxhp:      20,
				maxmp:      10,

				damroll:  NewDice(1, 11),
				protroll: NewDice(2, 4),
			}),
		},
	},
}
