package common

import "time"

type ProgressResponse struct {
	Lesson    string    `json:"lesson"`
	CharWPM   int       `json:"char_wpm"`
	EffWPM    int       `json:"eff_wpm"`
	Accuracy  float64   `json:"accuracy"`
	CreatedAt time.Time `json:"created_at"`
}

type CWSettingsResponse struct {
	CharWPM    int     `json:"char_wpm"`
	EffWPM     int     `json:"eff_wpm"`
	Freq       int     `json:"freq"`
	StartDelay float64 `json:"start_delay"`
}

type PageSettingsResponse struct {
	Theme     string `json:"theme"`
	Lang      string `json:"language"`
	CurLesson string `json:"cur_lesson"`
}

type AuthTokenPairResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type UserInfoResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HttpErrorResponse struct {
	Status int
	Err    string
}
