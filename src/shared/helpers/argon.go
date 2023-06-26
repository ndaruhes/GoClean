package helpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"go-clean/models/requests"
	"strings"

	"golang.org/x/crypto/argon2"
)

var passConfig = &requests.PasswordConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func GeneratePassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, passConfig.Time, passConfig.Memory, passConfig.Threads, passConfig.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, passConfig.Memory, passConfig.Time, passConfig.Threads, b64Salt, b64Hash)
	return full, nil
}

func ComparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &passConfig.Memory, &passConfig.Time, &passConfig.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	passConfig.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, passConfig.Time, passConfig.Memory, passConfig.Threads, passConfig.KeyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
