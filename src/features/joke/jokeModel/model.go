package jokeModel

import "time"

type Joke struct {
	ID          string
	Title       string
	Text        string
	Explanation *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time

	UserID string
}
