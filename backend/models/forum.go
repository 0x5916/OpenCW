package models

import "time"

type ForumCategory struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ForumThread struct {
	ID         uint      `json:"id"`
	CategoryID uint      `json:"category_id"`
	AuthorID   uint      `json:"author_id"`
	Title      string    `json:"title"`
	IsPinned   bool      `json:"is_pinned"`
	IsLocked   bool      `json:"is_locked"`
	CreatedAt  time.Time `json:"created_at"`
}

type ForumPost struct {
	ID        uint      `json:"id"`
	ThreadID  uint      `json:"thread_id"`
	AuthorID  uint      `json:"author_id"`
	Body      string    `json:"body"`
	ParentID  *uint     `json:"parent_id"` // for reply threading
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
