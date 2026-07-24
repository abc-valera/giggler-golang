module giggler-golang

go 1.26.5

require (
	github.com/danielgtaylor/huma/v2 v2.39.1
	github.com/go-playground/validator/v10 v10.30.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.54.0
	gorm.io/cli/gorm v0.2.4
	gorm.io/driver/postgres v1.6.2
	gorm.io/gorm v1.31.2
)

// Enable a custom package fork override so it's easy to swap back to the original if needed.
replace gorm.io/cli/gorm => github.com/abc-valera/gorm-cli v0.0.0-20260425221554-cbb053786aa2

require (
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	golang.org/x/exp v0.0.0-20260727155853-b88d891fe743 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)
