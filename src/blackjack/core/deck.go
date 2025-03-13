package core

import (
	"math/rand/v2"
	"strings"
)

const SuiteSize = 13
const Suits = 4
const DeckSize = SuiteSize * Suits

type Deck struct {
	Cards       []Card
	idx         int
	deckSize    int
	PreviewCard func(c Card)
	source      *rand.Rand
}

// Creates `shoe` of 1+ decks, unshuffled initially
func GenerateShoe(decks int) *Deck {
	shoe := GenerateDeck()
	for i := 1; i < decks; i++ {
		additionalDeck := GenerateDeck()
		shoe.Cards = append(shoe.Cards, additionalDeck.Cards...)
	}
	shoe.source = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	shoe.deckSize = decks * DeckSize
	return shoe
}

// Creates a non-suffled full 52 card deck
func GenerateDeck() *Deck {
	all := make([]Card, 0, DeckSize)
	for suiteIdx := SuiteFirst; suiteIdx < SuiteLast; suiteIdx++ {
		for cardIdx := 0; cardIdx < SuiteSize; cardIdx++ {
			card := cards[cardIdx]
			card.Suit = Suit(suiteIdx)
			all = append(all, card)
		}
	}

	return &Deck{
		idx:      0,
		deckSize: DeckSize,
		Cards:    all,
	}
}

func (d *Deck) ToString() string {
	cardsLeft := len(d.Cards[d.idx:])
	s := make([]string, 0, cardsLeft)
	for i := d.idx; i < cardsLeft; i++ {
		s = append(s, d.Cards[i].ToString())
	}
	return strings.Join(s, " ")
}

func (d *Deck) Shuffle() *Deck {
	// swiped & modified from https://go.dev/src/math/rand/v2/rand.go
	n := d.deckSize - 1
	for i := n - 1; i > 0; i-- {
		j := int(d.source.Uint64N(uint64(i + 1)))
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}

	d.idx = 0
	return d
}

func (d *Deck) Deal() Card {
	c := d.Cards[d.idx]
	d.idx++
	if d.PreviewCard != nil {
		d.PreviewCard(c)
	}

	return c
}

func (d *Deck) Remaining() int {
	return d.deckSize - d.idx
}

func (d *Deck) EstimateRemaining() float32 {
	return float32((d.deckSize - d.idx)) / 52.0
}
