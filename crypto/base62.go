package crypto

import (
	"encoding/base64"
	"strings"

	"github.com/golang/snappy"
	// "github.com/lytics/base62"
)

// Base62DecodeString returns the bytes represented by the base64 string s.
func Base62DecodeString(s string) ([]byte, error) {
	s = strings.ReplaceAll(s, "l0", "-")
	s = strings.ReplaceAll(s, "l1", "+")
	s = strings.ReplaceAll(s, "l2", "/")
	s = strings.ReplaceAll(s, "l", "=")
	s = strings.ReplaceAll(s, "-", "l")
	src, err := base64.StdEncoding.DecodeString(s)
	return src, err
}

// Base62EncodeToString returns the base64 encoding of src.
func Base62EncodeToString(src []byte) string {
	s := base64.StdEncoding.EncodeToString(src)
	s = strings.ReplaceAll(s, "l", "-")
	s = strings.ReplaceAll(s, "+", "l1")
	s = strings.ReplaceAll(s, "/", "l2")
	s = strings.ReplaceAll(s, "=", "l")
	s = strings.ReplaceAll(s, "-", "l0")
	return s
}

// Base62DecodeString2 returns the bytes represented by the base64 string s.
// with snappy
func Base62DecodeString2(s string) ([]byte, error) {
	src, err := Base62DecodeString(s)
	if err != nil {
		return nil, err
	}
	got, _ := snappy.Decode(nil, src)
	return got, nil
}

// Base62EncodeToString2 returns the base64 encoding of src.
// with snappy
func Base62EncodeToString2(src []byte) string {
	got := snappy.Encode(nil, src)
	dst := Base62EncodeToString(got)
	return dst
}

// Base62DecodeString3 returns the bytes represented by the base64 string s.
func Base62DecodeString3(s string, r rune) ([]byte, error) {
	k := string(r)
	s = strings.ReplaceAll(s, k+"0", "-")
	s = strings.ReplaceAll(s, k+"1", "+")
	s = strings.ReplaceAll(s, k+"2", "/")
	s = strings.ReplaceAll(s, k, "=")
	s = strings.ReplaceAll(s, "-", k)
	src, err := base64.StdEncoding.DecodeString(s)
	return src, err
}

// Base62EncodeToString3 returns the base64 encoding of src.
func Base62EncodeToString3(src []byte, r rune) string {
	k := string(r)
	s := base64.StdEncoding.EncodeToString(src)
	s = strings.ReplaceAll(s, k, "-")
	s = strings.ReplaceAll(s, "+", k+"1")
	s = strings.ReplaceAll(s, "/", k+"2")
	s = strings.ReplaceAll(s, "=", k)
	s = strings.ReplaceAll(s, "-", k+"0")
	return s
}
