package core

import (
	"fmt"
)

type Suit int

const (
	SuitUnknown Suit = iota
	SuitClubs
	SuitSpades
	SuitDiamonds
	SuitHearts
)

const (
	SuiteFirst = SuitClubs
	SuiteLast  = SuitHearts + 1
)

type Card struct {
	Name  string
	Value int
	Suit  Suit
}

var cards = []Card{
	{Name: "2", Value: 2},
	{Name: "3", Value: 3},
	{Name: "4", Value: 4},
	{Name: "5", Value: 5},
	{Name: "6", Value: 6},
	{Name: "7", Value: 7},
	{Name: "8", Value: 8},
	{Name: "9", Value: 9},
	{Name: "10", Value: 10},
	{Name: "J", Value: 10},
	{Name: "Q", Value: 10},
	{Name: "K", Value: 10},
	{Name: "A", Value: 11},
}

func SuitToString(s Suit) string {
	switch s {
	case SuitClubs:
		return `♣`
	case SuitSpades:
		return `♠`
	case SuitDiamonds:
		return `♦`
	case SuitHearts:
		return `♥`
	}
	return `?`
}

func (c *Card) ToString() string {
	return fmt.Sprintf("%s%s", SuitToString(c.Suit), c.Name)
}
