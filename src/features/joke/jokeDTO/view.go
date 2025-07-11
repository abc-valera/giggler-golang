package jokeDTO

import (
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

func NewViewDTO(joke jokeModel.Joke) *viewgen.JokeSchema {
	return &viewgen.JokeSchema{
		ID:          joke.ID,
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: viewDTO.NewOptString(joke.Explanation),
		CreatedAt:   joke.CreatedAt,
	}
}

func NewViewDTOs(jokes []jokeModel.Joke) viewgen.JokesSchema {
	viewJokes := make(viewgen.JokesSchema, len(jokes))
	for i, joke := range jokes {
		viewJokes[i] = *NewViewDTO(joke)
	}
	return viewJokes
}
