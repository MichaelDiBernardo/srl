package game

import (
	"testing"
)

type schedulerTestPair struct {
	actor string
	speed int
}

type schedulerTest struct {
	actors []schedulerTestPair
	want   []string
}

func TestScheduling(t *testing.T) {
	tests := []schedulerTest{
		{[]schedulerTestPair{{"A", 1}}, []string{"A", "A", "A", "A", "A"}},
		{[]schedulerTestPair{{"A", 2}}, []string{"A", "A", "A", "A", "A"}},
		{[]schedulerTestPair{{"A", 3}}, []string{"A", "A", "A", "A", "A"}},
		{[]schedulerTestPair{{"A", 4}}, []string{"A", "A", "A", "A", "A"}},
		{[]schedulerTestPair{{"A", 1}, {"B", 2}}, []string{"B", "A", "B", "A", "B", "B"}},
		{[]schedulerTestPair{{"A", 1}, {"B", 2}}, []string{"B", "A", "B", "A", "B", "B", "A", "B"}},
		{[]schedulerTestPair{{"A", 2}, {"B", 4}}, []string{"B", "A", "B", "B", "A", "B", "B", "A", "B", "B", "A"}},
	}

	g := newTestGame()

	for ti, test := range tests {
		s := NewScheduler()
		for _, actorspec := range test.actors {
			s.Add(lTestSpd(g, actorspec.actor, actorspec.speed))
		}

		actual := make([]string, 0)
		for i := 0; i < len(test.want); i++ {
			actual = append(actual, s.Next().Spec.Name)
		}

		for si, want := range test.want {
			if actual[si] != want {
				t.Errorf("TestScheduler: Test %d -- got %v, want %v", ti, actual, test.want)
				break
			}
		}
	}
}

func TestRemoveFromSchedule(t *testing.T) {
	// TODO
}

func lTestSpd(g *Game, name string, spd int) *Obj {
	return g.NewObj(&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Name:    name,
		Traits: &Traits{
			Sheet: NewMonsterSheet(&MonsterSheet{speed: spd}),
		},
	})
}
