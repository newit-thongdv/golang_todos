package main

import (
	"net/http"
	"time"

	"golang-todos/database"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func main() {
	database.Connect("root:admin@123@tcp(localhost:3306)/golang?parseTime=true")
	database.Migrate()

	now := time.Now().UTC()

	item := TodoItem{
		Id:          0,
		Title:       "Item 1",
		Description: "Description item 1",
		Status:      "doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": item,
		})
	})

	r.Run(":3001") // listen and serve on 0.0.0.0:3001
}
