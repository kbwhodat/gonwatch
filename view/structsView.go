package view

type Choices struct {
	choice string
}

var ChoiceList = []Choices{
	{choice: "recently watched"},
	{choice: "trending"},
	{choice: "movies"},
	{choice: "series"},
	{choice: "anime"},
	{choice: "sports"},
}

var TrendingChoiceList = []Choices{
	{choice: "movie"},
	{choice: "tv"},
}
