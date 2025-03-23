package jokeModel

import (
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/data/gormgen/gormModel"
	"github.com/abc-valera/giggler-golang/src/shared/view"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
)

func NewGormDTO(joke Joke) *gormModel.Joke {
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

func FromGormDTO(dto *gormModel.Joke) Joke {
	return Joke{
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

func NewGormDTOs(jokes []Joke) []*gormModel.Joke {
	var dtos []*gormModel.Joke
	for _, joke := range jokes {
		dtos = append(dtos, NewGormDTO(joke))
	}
	return dtos
}

func FromGormDTOs(dtos []*gormModel.Joke) []Joke {
	var jokes []Joke
	for _, dto := range dtos {
		jokes = append(jokes, FromGormDTO(dto))
	}
	return jokes
}

func NewViewDTO(joke Joke) *viewgen.JokeSchema {
	return &viewgen.JokeSchema{
		ID:          joke.ID,
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: view.NewOptString(joke.Explanation),
		CreatedAt:   joke.CreatedAt,
	}
}

func NewViewDTOs(jokes []Joke) viewgen.JokesSchema {
	viewJokes := make(viewgen.JokesSchema, len(jokes))
	for i, joke := range jokes {
		viewJokes[i] = *NewViewDTO(joke)
	}
	return viewJokes
}
