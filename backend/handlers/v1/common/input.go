package common

import (
	"opencw/models"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=16"`
	Email    string `json:"email"    binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateEmailInput struct {
	Email string `json:"email" binding:"required,email,max=254"`
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
	Lesson   string   `json:"lesson"   binding:"required"`
	CharWPM  int      `json:"char_wpm" binding:"required,min=5,max=50"`
	EffWPM   int      `json:"eff_wpm"  binding:"required,min=5,max=50"`
	Accuracy *float64 `json:"accuracy" binding:"required,min=0.0,max=1.0"`
}
