package model

import "time"

type Task struct {
	ID          int       `json:"id" query:"id"`
	UUID        string    `json:"uuid" query:"uuid"`
	Name        string    `json:"name" query:"name"`
	Description string    `json:"description" query:"description"`
	IsFinished  int       `json:"isFinished" query:"isFinished"`
	DueDate     time.Time `json:"dueDate" query:"duedate"`
	OwnerId     int       `json:"ownerId" query:"ownerId"`
	AssigneeId  int       `json:"assigneeId" query:"assigneeId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
