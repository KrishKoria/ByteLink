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
	assert.Truef(t, testStore.client != nil, "Expected client to be initialized")
}

func TestSetAndGet(t *testing.T) {
	initialURL := "https://www.google.com"
	shortURL := "test"
	userUUID := "test-user"
	SaveMapping(initialURL, shortURL, userUUID)
	longURL := GetLongUrl(shortURL)
	assert.Equal(t, initialURL, longURL)
}
