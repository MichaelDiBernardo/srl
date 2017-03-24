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
		Learner:  NewActorLearner,
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
			Dropper: NewItemDropper(&ItemDropper{
				num: 2,
			}),
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
						Chi: 10,
					},
				},
				speed: 2,
				maxhp: 20,
				maxmp: 10,

				attacks: []*MonsterAttack{
					{
						Attack: Attack{
							Melee:   0,
							Damroll: NewDice(2, 7),
							CritDiv: 4,
							Effects: Effects{},
							Verb:    "hits",
						},
						P: 1,
					},
					{
						Attack: Attack{
							Melee:   3,
							Damroll: NewDice(1, 11),
							CritDiv: 2,
							Effects: Effects{},
							Verb:    "stabs",
						},
						P: 1,
					},
				},
				defense: Defense{
					Evasion:  1,
					ProtDice: []Dice{NewDice(1, 4)},
					Effects:  NewEffects(map[Effect]int{}),
				},
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
			Dropper: NewItemDropper(&ItemDropper{
				num:   4,
				boost: 3,
			}),
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
					skills: skilllist{},
				},
				speed: 2,
				maxhp: 20,
				maxmp: 10,

				attacks: []*MonsterAttack{
					{
						Attack: Attack{
							Melee:   3,
							Damroll: NewDice(2, 9),
							CritDiv: 4,
							Effects: Effects{},
							Verb:    "claws",
						},
						P: 3,
					},
					{
						Attack: Attack{
							Melee:   0,
							Damroll: NewDice(2, 15),
							CritDiv: 4,
							Effects: NewEffects(map[Effect]int{BrandFire: 1}),
							Verb:    "bites",
						},
						P: 3,
					},
				},
				defense: Defense{
					Evasion:  5,
					ProtDice: []Dice{NewDice(2, 4)},
					Effects:  NewEffects(map[Effect]int{}),
				},
			}),
		},
	},
}
