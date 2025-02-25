package db

import (
	"time"
)

type Subscriber struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSubscriber(email string) *Subscriber {
	now := time.Now()
	return &Subscriber{
		Email:     email,
		CreatedAt: now,
	}
}
