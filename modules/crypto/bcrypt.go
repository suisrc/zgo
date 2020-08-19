package crypto

import "golang.org/x/crypto/bcrypt"

// CompareHashAndPassword compares a bcrypt hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// GenerateFromPassword returns the bcrypt hash of the password at the given
// cost. If the cost given is less than MinCost, the cost will be set to
// DefaultCost, instead. Use CompareHashAndPassword, as defined in this package,
// to compare the returned hashed password with its cleartext version.
func GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
