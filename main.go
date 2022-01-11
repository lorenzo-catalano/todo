package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task      string `json:"task" gorm:"not null;default:null"`
	Completed bool   `json:"completed"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Task{})

	router := gin.Default()
	router.Delims("{[{", "}]}")
	router.GET("/tasks", getRecipes(db))
	router.GET("/tasks/:id", getTaskById(db))
	router.PUT("/tasks", putTask(db))
	router.POST("/tasks/:id", postTask(db))
	router.DELETE("/tasks/:id", deleteTaskById(db))

	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("site/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run("localhost:8080")
}

func getRecipes(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var recipes []Task
		db.Find(&recipes)
		c.IndentedJSON(http.StatusOK, recipes)
	}

}

func getTaskById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var recipe Task
		db.Preload("Subtasks").First(&recipe, id)
		c.IndentedJSON(http.StatusOK, recipe)
	}
}

func postTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var taskToUpdate Task
		if err := c.BindJSON(&taskToUpdate); err != nil {
			return
		}
		var tid, _ = strconv.ParseUint(c.Param("id"), 10, 32)
		taskToUpdate.ID = uint(tid)
		result := db.Model(&taskToUpdate).Update("Completed", taskToUpdate.Completed)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, "error")
		} else {
			c.IndentedJSON(http.StatusOK, taskToUpdate)
		}
	}
}

func putTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTask Task
		if err := c.BindJSON(&newTask); err != nil {
			return
		}
		result := db.Create(&newTask)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, "error")
		} else {
			c.IndentedJSON(http.StatusCreated, newTask)
		}
	}
}
func deleteTaskById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&Task{}, id)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
