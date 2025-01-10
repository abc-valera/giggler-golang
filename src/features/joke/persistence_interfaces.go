package joke

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/gen/gorm/gormModel"
	"github.com/abc-valera/giggler-golang/gen/gorm/gormQuery"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"gorm.io/gorm"
)

var (
	command func(persistence.IDS) iCommand
	Query   func(persistence.IDS) iQuery
)

func init() {
	switch persistence.DbVal {
	case persistence.DbVariantGorm:
		command = func(dataStore persistence.IDS) iCommand {
			return persistence.NewGormGenericCommand(persistence.GormDS(dataStore), infraGormNewDTO)
		}
		Query = func(dataStore persistence.IDS) iQuery {
			return &gormQ{gormQuery.Use(persistence.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type (
	Joke struct {
		ID          string
		Title       string
		Text        string
		Explanation string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   time.Time

		UserID string
	}

	iCommand persistence.IGenericCommandCreateUpdateDelete[Joke]

	iQuery interface {
		GetByID(ctx context.Context, id string) (Joke, error)
		GetAllByUserID(ctx context.Context, userID string, s persistence.Selector) ([]Joke, error)
	}
)

type gormQ struct{ *gormQuery.Query }

func (q gormQ) GetByID(ctx context.Context, id string) (Joke, error) {
	dto, err := q.WithContext(ctx).Joke.Where(q.Joke.ID.Eq(id)).First()
	if err != nil {
		return Joke{}, err
	}
	return ifraGormNewJoke(dto), nil
}

func (q gormQ) GetAllByUserID(ctx context.Context, userID string, s persistence.Selector) ([]Joke, error) {
	dtos, err := q.WithContext(ctx).Joke.Where(q.Joke.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}
	return infraGormNewJokes(dtos), nil
}

func infraGormNewDTO(joke Joke) *gormModel.Joke {
	return &gormModel.Joke{
		ID:          joke.ID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: joke.Explanation,
		CreatedAt:   joke.CreatedAt,
		UpdatedAt:   joke.UpdatedAt,
		DeletedAt:   gorm.DeletedAt{Time: joke.DeletedAt, Valid: true},
		UserID:      joke.UserID,
	}
}

func ifraGormNewJoke(dto *gormModel.Joke) Joke {
	return Joke{
		ID:          dto.ID,
		Title:       dto.Title,
		Text:        dto.Text,
		Explanation: dto.Explanation,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		DeletedAt:   dto.DeletedAt.Time,
		UserID:      dto.UserID,
	}
}

func infraGormNewDTOs(jokes []Joke) []*gormModel.Joke {
	var dtos []*gormModel.Joke
	for _, joke := range jokes {
		dtos = append(dtos, infraGormNewDTO(joke))
	}
	return dtos
}

func infraGormNewJokes(dtos []*gormModel.Joke) []Joke {
	var jokes []Joke
	for _, dto := range dtos {
		jokes = append(jokes, ifraGormNewJoke(dto))
	}
	return jokes
}
