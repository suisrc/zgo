package middleware

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPath(t *testing.T) {
	path := JoinPath("s1", "23", "/s3")
	log.Println(path)
	path = JoinPath("s1", "/23/", "/s3")
	log.Println(path)
	path = JoinPath("/s1/", "/23/", "/s3")
	log.Println(path)

	assert.NotNil(t, nil)
}

func TestNull(t *testing.T) {
	var s1 interface{}
	s2 := s1.(string)
	log.Println(s2)
	assert.Nil(t, nil)
}
