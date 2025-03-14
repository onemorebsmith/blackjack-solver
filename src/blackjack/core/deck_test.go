package core

import "testing"

func ValidateDeck(t *testing.T, d *Deck, deckCount int) {
	{
		bySuite := map[Suit]int{}
		for _, v := range d.Cards {
			bySuite[v.Suit]++
		}

		if len(bySuite) != Suits {
			t.Fatalf("Incorrect suits in deck, got %d, expected %d", len(bySuite), Suits)
		}
		for k, v := range bySuite {
			if v != SuiteSize*deckCount {
				t.Fatalf("Suite %s has %d cards, not %d", SuitToString(k), v, SuiteSize)
			}
		}
	}

	{
		byName := map[string]int{}
		for _, v := range d.Cards {
			byName[v.Name]++
		}

		if len(byName) != SuiteSize {
			t.Fatalf("Incorrect number of cards in deck, got %d, expected %d", len(byName), Suits)
			t.Fail()
		}
		for k, v := range byName {
			if v != Suits*deckCount {
				t.Fatalf("Card %s has %d cards, not %d", k, v, SuiteSize)
			}
		}
	}
}

func TestDeckGeneration(t *testing.T) {
	d := GenerateDeck()
	if len(d.Cards) != DeckSize {
		t.Fail()
	}
	ValidateDeck(t, d, 1)
}

func TestDeckSuffle(t *testing.T) {
	d := GenerateDeck()
	if len(d.Cards) != DeckSize {
		t.Fail()
	}
	d.Shuffle()
	ValidateDeck(t, d, 1)
}

func TestShoeGeneration(t *testing.T) {
	for i := 1; i < 20; i++ {
		shoe := GenerateShoe(i)
		ValidateDeck(t, shoe, i)
	}
}
