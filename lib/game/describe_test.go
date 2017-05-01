package game

import (
	"testing"
)

type atkDescTest struct {
	atk  Attack
	want string
}

var atkDescTests = []atkDescTest{
	{atk: Attack{}, want: "(+0)"},
	{atk: Attack{Hit: -2}, want: "(-2)"},
	{atk: Attack{Hit: 3}, want: "(+3)"},
	{atk: Attack{Hit: 0, Damroll: NewDice(2, 3)}, want: "(+0,2d3)"},
	{atk: Attack{Hit: -2, Damroll: NewDice(3, 5)}, want: "(-2,3d5)"},
}

func TestDescribeAttack(t *testing.T) {
	for i, test := range atkDescTests {
		if d, w := test.atk.Describe(), test.want; d != w {
			t.Errorf(`atk.Describe() test %d: was "%s", want "%s"`, i, d, w)
		}
	}
}

type defDescTest struct {
	def  Defense
	want string
}

var defDescTests = []defDescTest{
	{def: Defense{}, want: "[+0]"},
	{def: Defense{Evasion: 2}, want: "[+2]"},
	{def: Defense{Evasion: -3}, want: "[-3]"},
	{def: Defense{ProtDice: []Dice{NewDice(1, 3)}}, want: "[+0,1-3]"},
	{def: Defense{ProtDice: []Dice{NewDice(1, 3), NewDice(2, 4)}}, want: "[+0,3-11]"},
	{def: Defense{Evasion: 1, ProtDice: []Dice{NewDice(1, 3)}}, want: "[+1,1-3]"},
	{def: Defense{Evasion: -1, ProtDice: []Dice{NewDice(1, 3)}}, want: "[-1,1-3]"},
	{def: Defense{ProtDice: []Dice{NewDice(1, 3), NewDice(2, 4)}, CorrDice: []Dice{NewDice(2, 4)}}, want: "[+0,0-9]"},
}

func TestDescribeDefense(t *testing.T) {
	for i, test := range defDescTests {
		if d, w := test.def.Describe(), test.want; d != w {
			t.Errorf(`def.Describe() test %d: was "%s", want "%s"`, i, d, w)
		}
	}
}
