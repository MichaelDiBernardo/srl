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
		Senser:   NewActorSenser,
		Ticker:   NewActorTicker,
	},
}

var Monsters = []*Spec{
	&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 2,
		},
		Traits: &Traits{
			Mover: NewActorMover,
			AI: NewSMAI(SMAI{
				Brain: SMAIWanderer,
				Personality: &Personality{
					Fear:        25,
					Persistence: 1000,
				},
			}),
			Fighter: NewActorFighter,
			Packer:  NewActorPacker,
			Senser:  NewActorSenser,
			Ticker:  NewActorTicker,
			Sheet: NewMonsterSheet(&MonsterSheet{
				stats: &stats{
					stats: statlist{
						Str: 2,
						Agi: 0,
						Vit: 1,
						Mnd: 0,
					},
				},
				skills: &skills{
					skills: skilllist{
						Melee:   1,
						Evasion: 1,
					},
				},
				speed:      2,
				critdivmod: 4,
				maxhp:      20,
				maxmp:      10,

				damroll:  NewDice(2, 7),
				protroll: NewDice(1, 4),

				atkeffects: NewEffects(map[Effect]int{EffectBlind: 1}),
			}),
		},
	},
	&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecAnt,
		Name:    "DRAGON",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Mover: NewActorMover,
			AI: NewSMAI(SMAI{
				Brain: SMAITerritorial,
				Personality: &Personality{
					Fear:        50,
					Persistence: 0,
				},
			}),
			Fighter: NewActorFighter,
			Packer:  NewActorPacker,
			Senser:  NewActorSenser,
			Ticker:  NewActorTicker,
			Sheet: NewMonsterSheet(&MonsterSheet{
				stats: &stats{
					stats: statlist{
						Str: 3,
						Agi: 3,
						Vit: 1,
						Mnd: 0,
					},
				},
				skills: &skills{
					skills: skilllist{
						Melee:   5,
						Evasion: 5,
					},
				},
				speed:      2,
				critdivmod: 4,
				maxhp:      20,
				maxmp:      10,

				damroll:  NewDice(1, 11),
				protroll: NewDice(2, 4),

				atkeffects: NewEffects(map[Effect]int{BrandFire: 1}),
			}),
		},
	},
}
