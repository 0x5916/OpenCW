package common

import (
	"time"

	"github.com/google/uuid"
)

type ProgressResponse struct {
	Lesson          string     `json:"lesson"`
	CharWPM         int        `json:"char_wpm"`
	EffWPM          int        `json:"eff_wpm"`
	Accuracy        float64    `json:"accuracy"`
	CreatedAt       time.Time  `json:"created_at"`
	ClientCreatedAt *time.Time `json:"client_created_at"`
}

type CWSettingsResponse struct {
	CharWPM    int       `json:"char_wpm"`
	EffWPM     int       `json:"eff_wpm"`
	Freq       int       `json:"freq"`
	StartDelay float64   `json:"start_delay"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PageSettingsResponse struct {
	Theme     string    `json:"theme"`
	Lang      string    `json:"language"`
	CurLesson int       `json:"cur_lesson"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthTokenPairResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type UserInfoResponse struct {
	CallSign      *string   `json:"call_sign"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

type HttpErrorResponse struct {
	Status int
	Code   string
	Err    string
}

type ForumCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ForumThreadResponse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	Title      string    `json:"title"`
	IsPinned   bool      `json:"is_pinned"`
	IsLocked   bool      `json:"is_locked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ForumPostResponse struct {
	ID        uuid.UUID  `json:"id"`
	ThreadID  uuid.UUID  `json:"thread_id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	Body      string     `json:"body"`
	ParentID  *uuid.UUID `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
