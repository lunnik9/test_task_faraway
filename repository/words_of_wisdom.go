package repository

import "math/rand"

type WordsOfWisdom struct {
	words []string
}

func NewWordsOfWisdom() *WordsOfWisdom {
	return &WordsOfWisdom{
		words: []string{
			"Chuck Norris doesn’t read books. He stares them down until he gets the information he wants.",
			"Time waits for no man. Unless that man is Chuck Norris.",
			"Chuck Norris can divide by zero.",
			"Chuck Norris beat the sun in a staring contest.",
			"Chuck Norris can hear sign language.",
			"Chuck Norris makes onions cry.",
			"Chuck Norris can slam a revolving door.",
			"Chuck Norris's computer has no backspace button, Chuck Norris doesn't make mistakes.",
			"Chuck Norris can build a snowman out of rain.",
			"Chuck Norris can strangle you with a cordless phone.",
			"Chuck Norris can do a wheelie on a unicycle.",
			"Chuck Norris stands faster than anyone can run.",
			"When Chuck Norris enters a room, he doesn’t turn the lights on. He turns the dark off.",
			"The only time Chuck Norris was wrong was when he thought he had made a mistake.",
			"Chuck Norris can play the violin with a piano.",
			"Chuck Norris’s shadow has been on the ‘best-dressed’ list twice.",
			"Chuck Norris doesn’t flush the toilet, he scares the crap out of it.",
			"Chuck Norris can speak braille.",
			"Chuck Norris can do push-ups with his beard.",
			"Chuck Norris doesn’t wear a watch. HE decides what time it is.",
			"Chuck Norris doesn’t churn butter. He roundhouse kicks the cows and the butter comes straight out.",
			"When Chuck Norris does a pushup, he isn’t lifting himself up, he’s pushing the Earth down.",
			"Chuck Norris tells Simon what to do.",
			"Chuck Norris can tie his shoes with his feet.",
			"Chuck Norris can set ants on fire with a magnifying glass. At night.",
			"Chuck Norris has a diary. It's called the Guinness Book of World Records.",
			"Chuck Norris can light a fire by rubbing two ice cubes together.",
			"Chuck Norris makes snow angels in concrete.",
			"Chuck Norris’s blood type is AK-47.",
			"Chuck Norris can kill two stones with one bird.",
		},
	}
}

func (w *WordsOfWisdom) GetQuote() string {
	return w.words[rand.Intn(len(w.words))]
}
