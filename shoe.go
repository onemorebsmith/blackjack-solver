package main

func GenerateShoe(decks int) *Deck {
	shoe := GenerateDeck()
	for i := 1; i < decks; i++ {
		additionalDeck := GenerateDeck()
		shoe.Cards = append(shoe.Cards, additionalDeck.Cards...)
	}

	return shoe
}
