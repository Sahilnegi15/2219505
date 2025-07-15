package main

import (
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

var (
	urlStore  = make(map[string]URLMapping)
	storeLock sync.RWMutex
)

func GetOrGenerateCode(custom string) string {
	if custom != "" {
		return custom
	}
	id, _ := shortid.Generate()
	return id
}

func SaveURLMapping(code string, url string, validity time.Duration) bool {
	storeLock.Lock()
	defer storeLock.Unlock()

	// Reject if custom code exists
	if _, exists := urlStore[code]; exists {
		return false
	}

	urlStore[code] = URLMapping{
		OriginalURL: url,
		ExpiresAt:   time.Now().Add(validity),
	}

	return true
}

func GetOriginalURL(code string) (string, bool) {
	storeLock.RLock()
	defer storeLock.RUnlock()

	data, exists := urlStore[code]
	if !exists || time.Now().After(data.ExpiresAt) {
		return "", false
	}

	return data.OriginalURL, true
}
