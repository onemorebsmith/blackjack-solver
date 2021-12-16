package main

import (
	"fmt"
	"strings"
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
	Name       string
	Value      int
	CountValue int
	Suit       Suit
}

var cards = []Card{
	{Name: "2", Value: 2, CountValue: 1},
	{Name: "3", Value: 3, CountValue: 1},
	{Name: "4", Value: 4, CountValue: 1},
	{Name: "5", Value: 5, CountValue: 1},
	{Name: "6", Value: 6, CountValue: 1},
	{Name: "7", Value: 7, CountValue: 0},
	{Name: "8", Value: 8, CountValue: 0},
	{Name: "9", Value: 9, CountValue: 0},
	{Name: "10", Value: 10, CountValue: -1},
	{Name: "J", Value: 10, CountValue: -1},
	{Name: "Q", Value: 10, CountValue: -1},
	{Name: "K", Value: 10, CountValue: -1},
	{Name: "A", Value: 11, CountValue: -1},
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

func PrintHand(c []Card) string {
	s := make([]string, 0, len(c))
	for _, v := range c {
		s = append(s, v.ToString())
	}
	return strings.Join(s, " ")
}
