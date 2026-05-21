package model

import "time"

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}
