package passworder

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/argon2"
)

var errPasswordDontMatch = errutil.NewCodeMessage(errutil.CodeInvalidArgument, "passwords don't match")

const (
	hashMemory      uint32 = 64 * 1024
	hashIterations  uint32 = 3
	hashParallelism uint8  = 2
	hashSaltLength  uint32 = 16
	hashKeyLength   uint32 = 32
)

// Hash returns hash of the provided password
func Hash(ctx context.Context, password string) (string, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(hashSaltLength)
	if err != nil {
		return "", err
	}

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

// Check checks if provided password matches provided hash,
// if matches returns nil, else returns error
func Check(ctx context.Context, providedPassword, realHash string) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	// Extract the parameters, salt and derived key from the encoded password hash.
	salt, hash, err := hashDecode(realHash)
	if err != nil {
		return errutil.NewInternalErr(err)
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(providedPassword), salt, hashIterations, hashMemory, hashParallelism, hashKeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}
	return errPasswordDontMatch
}

func hashDecode(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, errPasswordDontMatch
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, errPasswordDontMatch
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

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
