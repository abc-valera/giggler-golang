package jokeDto

import (
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormModel"
	"giggler-golang/src/shared/view/viewUtil"
	"giggler-golang/src/shared/view/viewgen"
)

func NewGorm(joke jokeModel.Joke) *gormModel.Joke {
	return &gormModel.Joke{
		ID:          joke.ID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: joke.Explanation,
		CreatedAt:   joke.CreatedAt,
		UpdatedAt:   joke.UpdatedAt,
		DeletedAt:   data.NewDtoGormDeletedAt(joke.DeletedAt),
		UserID:      joke.UserID,
	}
}

func FromGorm(dto *gormModel.Joke) jokeModel.Joke {
	return jokeModel.Joke{
		ID:          dto.ID,
		Title:       dto.Title,
		Text:        dto.Text,
		Explanation: dto.Explanation,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		DeletedAt:   data.NewDomainDeletedAt(dto.DeletedAt),
		UserID:      dto.UserID,
	}
}

func NewGormDTOs(jokes []jokeModel.Joke) []*gormModel.Joke {
	var dtos []*gormModel.Joke
	for _, joke := range jokes {
		dtos = append(dtos, NewGorm(joke))
	}
	return dtos
}

func FromGorms(dtos []*gormModel.Joke) []jokeModel.Joke {
	var jokes []jokeModel.Joke
	for _, dto := range dtos {
		jokes = append(jokes, FromGorm(dto))
	}
	return jokes
}

func NewViewDTO(joke jokeModel.Joke) *viewgen.JokeSchema {
	return &viewgen.JokeSchema{
		ID:          joke.ID,
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: viewUtil.NewOptString(joke.Explanation),
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
