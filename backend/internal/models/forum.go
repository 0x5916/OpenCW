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
	CategoryID uuid.UUID `json:"category_id" gorm:"not null;index:idx_forum_threads_category_sort,priority:1"`
	AuthorID   uuid.UUID `json:"author_id" gorm:"not null;index"`
	Title      string    `json:"title" gorm:"not null"`
	IsPinned   bool      `json:"is_pinned" gorm:"not null;index:idx_forum_threads_category_sort,priority:2,sort:desc"`
	IsLocked   bool      `json:"is_locked" gorm:"not null"`
}

func (ForumThread) TableName() string {
	return "forum_thread"
}

type ForumPost struct {
	Base
	ThreadID uuid.UUID  `json:"thread_id" gorm:"not null;index:idx_forum_posts_thread_sort,priority:1"`
	AuthorID uuid.UUID  `json:"author_id" gorm:"not null;index"`
	Body     string     `json:"body" gorm:"not null"`
	ParentID *uuid.UUID `json:"parent_id" gorm:"index"` // for reply threading
}

func (ForumPost) TableName() string {
	return "forum_post"
}
