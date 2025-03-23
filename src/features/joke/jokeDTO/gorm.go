package jokeDTO

import (
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data/dbDTO"
	"giggler-golang/src/shared/data/dbgen/gormModel"
)

func NewGorm(joke jokeModel.Joke) *gormModel.Joke {
	return &gormModel.Joke{
		ID:          joke.ID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: joke.Explanation,
		CreatedAt:   joke.CreatedAt,
		UpdatedAt:   joke.UpdatedAt,
		DeletedAt:   dbDTO.NewDeletedAt(joke.DeletedAt),
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
		DeletedAt:   dbDTO.NewDomainDeletedAt(dto.DeletedAt),
		UserID:      dto.UserID,
	}
}

func NewGorms(jokes []jokeModel.Joke) []*gormModel.Joke {
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
