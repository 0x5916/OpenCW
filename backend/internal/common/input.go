package common

import (
	"time"

	"opencw/internal/models"

	"github.com/google/uuid"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email"    binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password"   binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type VerifyEmailInput struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

type UpdateCallSignInput struct {
	CallSign string `json:"call_sign" binding:"required,max=254"`
}

type UpdateEmailInput struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=256"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=256"`
}

type CWSettingsInput struct {
	CharWPM    int      `json:"char_wpm"    binding:"required,min=5,max=50"`
	EffWPM     int      `json:"eff_wpm"     binding:"required,min=5,max=50"`
	Freq       int      `json:"freq"        binding:"required,min=300,max=2000"`
	StartDelay *float64 `json:"start_delay" binding:"required,min=0.0,max=10.0"`
}

func FromCwSettingsModel(obj models.CWSettings) CWSettingsInput {
	return CWSettingsInput{
		CharWPM:    obj.CharWPM,
		EffWPM:     obj.EffWPM,
		Freq:       obj.Freq,
		StartDelay: &obj.StartDelay,
	}
}

type PageSettingsInput struct {
	Lang      string `json:"language"   binding:"required"`
	CurLesson int    `json:"cur_lesson" binding:"required"`
}

func FromPageSettingsModel(obj models.PageSettings) PageSettingsInput {
	return PageSettingsInput{
		Lang:      obj.Lang,
		CurLesson: obj.CurLesson,
	}
}

type ProgressInput struct {
	Lesson          int        `json:"lesson"            binding:"required"`
	CharWPM         int        `json:"char_wpm"          binding:"required,min=5,max=50"`
	EffWPM          int        `json:"eff_wpm"           binding:"required,min=5,max=50"`
	Accuracy        *float64   `json:"accuracy"          binding:"required,min=0.0,max=1.0"`
	ClientCreatedAt *time.Time `json:"client_created_at"`
}

// Forum categories are a fixed set: general, help, showcase, feedback.
type CreateThreadInput struct {
	Category string `json:"category" binding:"required,oneof=general help showcase feedback"`
	Title    string `json:"title"    binding:"required,min=3,max=200"`
	Body     string `json:"body"     binding:"required,min=1,max=10000"`
}

// CreateReplyInput's ParentID, when set, must reference a live reply in the
// same thread, allowing nested replies.
type CreateReplyInput struct {
	Body     string     `json:"body"      binding:"required,min=1,max=10000"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// ListThreadsQuery is bound from query parameters of GET /v1/forum/threads.
type ListThreadsQuery struct {
	Category string `form:"category" binding:"omitempty,oneof=general help showcase feedback"`
	Limit    int    `form:"limit"    binding:"omitempty,min=1,max=100"`
	Cursor   string `form:"cursor"`
}
