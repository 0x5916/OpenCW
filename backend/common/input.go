package common

import (
	"opencw/models"
	"time"
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
	Theme     string `json:"theme"      binding:"required,oneof=auto dark light"`
	Lang      string `json:"language"   binding:"required"`
	CurLesson int    `json:"cur_lesson" binding:"required"`
}

func FromPageSettingsModel(obj models.PageSettings) PageSettingsInput {
	return PageSettingsInput{
		Theme:     obj.Theme,
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

type ForumPostInput struct {
}
