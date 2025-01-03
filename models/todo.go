package models

import (
	"time"
)

/**
 * Todo
 */
type Todo struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (Todo) TableName() string { return "todos" }

/**
 * TodoCreation
 */
type TodoCreation struct {
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (TodoCreation) TableName() string { return Todo{}.TableName() }

/**
 * TodoUpdate
 */
type TodoUpdate struct {
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (TodoUpdate) TableName() string { return Todo{}.TableName() }
