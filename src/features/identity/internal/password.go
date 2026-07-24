package internal

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/validate"
)

var errInvalidPassword = errors.Join(errutil.ErrValidation, errors.New("invalid password"))

const (
	hashMemory      uint32 = 64 * 1024
	hashIterations  uint32 = 3
	hashParallelism uint8  = 2
	hashSaltLength  uint32 = 16
	hashKeyLength   uint32 = 32
)

func NewPasswordHash(ctx context.Context, password string) (string, error) {
	if err := validate.Var(password, "required,min=2,max=32"); err != nil {
		return "", errors.Join(errInvalidPassword, err)
	}

	// Generate a cryptographically secure random salt.
	salt := generateRandomBytes(hashSaltLength)

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id variant.
	hash := argon2.IDKey([]byte(password), salt, hashIterations, hashMemory, hashParallelism, hashKeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hashMemory, hashIterations, hashParallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func PasswordVerify(ctx context.Context, pass, hashedPass string) error {
	// Extract the parameters, salt and derived key from the encoded password hash.
	salt, hash, err := passwordHashDecode(hashedPass)
	if err != nil {
		return err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(pass), salt, hashIterations, hashMemory, hashParallelism, hashKeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}

	return errInvalidPassword
}

func passwordHashDecode(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, errInvalidPassword
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, errInvalidPassword
	}
	if version != argon2.Version {
		return nil, nil, errInvalidPassword
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, errInvalidPassword
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, errInvalidPassword
	}

	return salt, hash, nil
}

func generateRandomBytes(n uint32) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
