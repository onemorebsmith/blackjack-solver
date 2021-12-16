package blackjack

var H17Rules = []RuleShorthand{
	{DealerCard: 2, PlayerDoublesOn: []int{10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 12}},
	{DealerCard: 3, PlayerDoublesOn: []int{9, 10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 12}},
	{DealerCard: 4, PlayerDoublesOn: []int{9, 10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8}},
	{DealerCard: 5, PlayerDoublesOn: []int{9, 10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8}},
	{DealerCard: 6, PlayerDoublesOn: []int{9, 10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8}},
	{DealerCard: 7, PlayerDoublesOn: []int{10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16}},
	{DealerCard: 8, PlayerDoublesOn: []int{10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16}},
	{DealerCard: 9, PlayerDoublesOn: []int{10, 11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15, 16}},
	{DealerCard: 10, PlayerDoublesOn: []int{11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16}},
	{DealerCard: 11, PlayerDoublesOn: []int{11}, PlayerHitsOn: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16}},

	{Soft: true, DealerCard: 2, PlayerDoublesOn: []int{18}, PlayerHitsOn: []int{13, 14, 15, 16, 17}},
	{Soft: true, DealerCard: 3, PlayerDoublesOn: []int{18, 17}, PlayerHitsOn: []int{13, 14, 15, 16}},
	{Soft: true, DealerCard: 4, PlayerDoublesOn: []int{18, 17, 16, 15}, PlayerHitsOn: []int{13, 14}},
	{Soft: true, DealerCard: 5, PlayerDoublesOn: []int{18, 17, 16, 15, 14, 13}, PlayerHitsOn: []int{}},
	{Soft: true, DealerCard: 6, PlayerDoublesOn: []int{19, 18, 17, 16, 15, 14, 13}, PlayerHitsOn: []int{}},
	{Soft: true, DealerCard: 7, PlayerHitsOn: []int{17, 16, 15, 14, 13}},
	{Soft: true, DealerCard: 8, PlayerHitsOn: []int{17, 16, 15, 14, 13}},
	{Soft: true, DealerCard: 9, PlayerHitsOn: []int{18, 17, 16, 15, 14, 13}},
	{Soft: true, DealerCard: 10, PlayerHitsOn: []int{18, 17, 16, 15, 14, 13}},
	{Soft: true, DealerCard: 11, PlayerHitsOn: []int{18, 17, 16, 15, 14, 13}},
}

var H17Splits = []SplitRule{
	{PlayerCard: 11, DealerUpcard: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
	{PlayerCard: 10, DealerUpcard: []int{}},
	{PlayerCard: 9, DealerUpcard: []int{2, 3, 4, 5, 6, 8, 9}},
	{PlayerCard: 8, DealerUpcard: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
	{PlayerCard: 7, DealerUpcard: []int{2, 3, 4, 5, 6, 7}},
	{PlayerCard: 6, DealerUpcard: []int{2, 3, 4, 5, 6}},
	{PlayerCard: 5, DealerUpcard: []int{}},
	{PlayerCard: 4, DealerUpcard: []int{5, 6}},
	{PlayerCard: 3, DealerUpcard: []int{2, 3, 4, 5, 6, 7}},
	{PlayerCard: 2, DealerUpcard: []int{2, 3, 4, 5, 6, 7}},
}
