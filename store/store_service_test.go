package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStore = &StoreService{}

func init() {
	testStore = InitializeStoreService()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStore.client != nil, "Expected client to be initialized")
}

func TestSetAndGet(t *testing.T) {
	initialURL := "https://www.google.com"
	userUUID := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"
	SaveMapping(shortURL, initialURL, userUUID)
	longURL := GetLongUrl(shortURL)
	assert.Equal(t, initialURL, longURL)
}
