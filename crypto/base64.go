package crypto

import (
	"encoding/base64"
)

// Base64DecodeString returns the bytes represented by the base64 string s.
// RFC 4648.
func Base64DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64EncodeToString returns the base64 encoding of src.
// RFC 4648.
func Base64EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64DecodeStringURL returns the bytes represented by the base64 string s.
func Base64DecodeStringURL(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Base64EncodeToStringURL returns the base64 encoding of src.
func Base64EncodeToStringURL(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

// Base64DecodeStringMIME returns the bytes represented by the base64 string s.
// RFC 2045
// func Base64DecodeStringMIME(s string) ([]byte, error) {
// 	// var buf bytes.Buffer
// 	// var buf = bytes.NewBuffer(make([]byte, len(s)*4))
// 	// for _, b := range s {
// 	// 	buf.WriteRune(b)
// 	// }
// 	// enc := base64.StdEncoding
// 	// dbuf := make([]byte, enc.DecodedLen(len(s)))
// 	// n, err := enc.Decode(dbuf, []byte(s))
// 	// return dbuf[:n], err
// 	return base64.StdEncoding.DecodeString(s)
// }
