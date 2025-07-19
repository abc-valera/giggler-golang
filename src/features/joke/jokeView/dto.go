package jokeView

import (
	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

func NewJokeModel(joke *jokeData.Joke) *viewgen.JokeSchema {
	return &viewgen.JokeSchema{
		ID:          joke.ID,
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: viewDTO.NewOptString(joke.Explanation),
		CreatedAt:   joke.CreatedAt,
	}
}

func NewJokeModels(jokes []*jokeData.Joke) viewgen.JokesSchema {
	viewJokes := make(viewgen.JokesSchema, len(jokes))
	for i, joke := range jokes {
		viewJokes[i] = *NewJokeModel(joke)
	}
	return viewJokes
}
