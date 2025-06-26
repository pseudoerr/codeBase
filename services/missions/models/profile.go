package models

type Profile struct {
	TotalPoints  int      `json:"total_points"`
	Level        string   `json:"level"`
	Achievements []string `json:"achievements"`
}
