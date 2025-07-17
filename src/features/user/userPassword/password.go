package userPassword

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
	"giggler-golang/src/shared/otel"
)

var ErrInvalidPass = errutil.NewCode(errutil.CodeInvalidArgument, errors.New("invalid password"))

// TODO: move this to env file
const (
	hashMemory      uint32 = 64 * 1024
	hashIterations  uint32 = 3
	hashParallelism uint8  = 2
	hashSaltLength  uint32 = 16
	hashKeyLength   uint32 = 32
)

func Hash(ctx context.Context, password string) string {
	_, span := otel.Trace(ctx)
	defer span.End()

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

	return encodedHash
}

// IsReal returns true if the password matches provided hash.
func IsReal(ctx context.Context, pass, hashedPass string) bool {
	_, span := otel.Trace(ctx)
	defer span.End()

	// Extract the parameters, salt and derived key from the encoded password hash.
	salt, hash, err := hashDecode(hashedPass)
	if err != nil {
		return false
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(pass), salt, hashIterations, hashMemory, hashParallelism, hashKeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}

	return false
}

func hashDecode(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, fmt.Errorf("unsupported version: %d", version)
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
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
