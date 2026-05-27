package model

import "time"

type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

func (t *Task) IsEmpty() bool {
	return t.ID == 0 && t.Text == "" && t.Due == time.Time{}
}
