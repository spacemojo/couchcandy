package couchcandy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBaseURL(t *testing.T) {

	url := createBaseURL(Session{
		Host:     "http://127.0.0.1",
		Port:     5984,
		Username: "unittest",
		Password: "test",
	})

	assert.NotEmpty(t, url)
	assert.Equal(t, "http://unittest:test@127.0.0.1:5984", url)

}

func TestCreateDatabaseURL(t *testing.T) {

	url := createDatabaseURL(Session{
		Host:     "http://127.0.0.1",
		Port:     5984,
		Username: "unit",
		Password: "test",
		Database: "teste",
	})

	assert.NotEmpty(t, url)
	assert.Equal(t, "http://unit:test@127.0.0.1:5984/teste", url)

}
