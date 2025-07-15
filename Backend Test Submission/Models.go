package main

import "time"

type URLMapping struct {
	OriginalURL string
	ExpiresAt   time.Time
}
