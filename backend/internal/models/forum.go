package models

import "github.com/google/uuid"

type ForumCategory struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ForumCategory) TableName() string {
	return "forum_category"
}

type ForumThread struct {
	Base
	CategoryID uuid.UUID `json:"category_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	Title      string    `json:"title"`
	IsPinned   bool      `json:"is_pinned"`
	IsLocked   bool      `json:"is_locked"`
}

func (ForumThread) TableName() string {
	return "forum_thread"
}

type ForumPost struct {
	Base
	ThreadID uuid.UUID  `json:"thread_id"`
	AuthorID uuid.UUID  `json:"author_id"`
	Body     string     `json:"body"`
	ParentID *uuid.UUID `json:"parent_id"` // for reply threading
}

func (ForumPost) TableName() string {
	return "forum_post"
}
