package jokeData

// Query defines type-safe query methods for Joke model
type Query[T any] interface {
	// GetByID retrieves a joke by its ID
	//
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id string) (T, error)

	// FindByUserID retrieves all jokes for a specific user
	//
	// SELECT * FROM @@table WHERE user_id=@userID
	FindByUserID(userID string) ([]T, error)

	// SearchJokes searches jokes with optional filters
	//
	// SELECT * FROM @@table
	// {{where}}
	//   {{if title != ""}} title LIKE @title {{end}}
	//   {{if userID != ""}} AND user_id=@userID {{end}}
	// {{end}}
	// ORDER BY created_at DESC
	SearchJokes(title, userID string) ([]T, error)

	// FindByTitle retrieves a joke by user ID and title
	//
	// SELECT * FROM @@table WHERE user_id=@userID AND title=@title
	FindByTitle(userID, title string) (T, error)
}
