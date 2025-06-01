package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"TaipeiCityDashboardBE/app/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Question struct {
	ID      int           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string        `gorm:"column:question" json:"question"`
	Options pq.StringArray `gorm:"column:options;type:text[]" json:"options"`
}
type Answer struct {
	ID     int    `gorm:"column:id" json:"id"`
	Answer int `gorm:"column:answer" json:"choice"`
}

func CreateQuestion(c *gin.Context) {
	// Example of creating a question with auto-generated ID and options
	// var req struct {
	// 	Title   string   `gorm:"column:question" json:"title"`
	// 	Options []string `json:"options"`
	// }
	var question Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received question: %+v\n", question)

	// Simulate auto-generating an ID (in real app, use DB auto-increment or UUID)

	// question := map[string]interface{}{
	// 	"id":      nil,
	// 	"question":   req.Title,
	// 	"options": req.Options,
	// }

	// Write to db
	err := models.DBDashboard.
		Table("question_list").
		Create(&question).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert question"})
		return
	}
	// create a table to record answers
	// add a column to components
	// add a column to components chart
	// add a column to query chart
	// update chart title simutaneously
	err = models.DBManager.
		Table("components").
		Where("index = 'question_statistic'").
		Update("name", question.Title).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chart title"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{
		"message": "Question created successfully",
		"data":    question,
	})
}
func GetQuestionByID(c *gin.Context) {
	// Example of getting a question by ID
	// In a real application, you would fetch the question from a database
	id := c.Param("id")
	
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	qid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var question Question
	err = models.DBDashboard.
		Table("question_list").
		Where("id = ?", qid).
		First(&question).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question retrieved successfully",
		"data":    question,
	})
}
func GetAllQuestions(c *gin.Context) {
	// Example of getting all questions
	// In a real application, you would fetch all questions from a database

	// Simulate fetching questions
	var questions []Question

	// Read from db
	err := models.DBDashboard.
		Table("question_list").
		Raw(`SELECT id, question, options FROM question_list`).
		Scan(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
}
func VoteQuestionByID(c *gin.Context) {
	// 會傳題目 id 和選擇的選項回來
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req struct {
		Choice int `json:"choice"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// Simulate voting
	fmt.Printf("Vote recorded for question ID: %s, choice: %s\n", id, req.Choice)
	answer := Answer{
		ID:     id,
		Answer: req.Choice,
	}
	// Write to db
	err = models.DBDashboard.
		Table("answer_record").
		Create(&answer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote recorded successfully",
		"data":    answer,
	})
}

func GetVoteResultByID(c *gin.Context) {
	id := c.Param("id")


	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// 模擬從資料庫取得題目與投票結果
	question := map[string]interface{}{
		"id":      id,
		"title":   "Sample Question",
		"options": []string{"Option 1", "Option 2", "Option 3"},
	}

	// 模擬選項投票數
	optionCounts := map[string]int{
		"Option 1": 10,
		"Option 2": 5,
		"Option 3": 2,
	}

	// 若有時間區間，可根據 startTime, endTime 過濾資料
	// 這裡僅為示範，實際應查詢資料庫

	result := gin.H{
		"id":         question["id"],
		"title":      question["title"],
		"options": []gin.H{
			{"name": "Option 1", "count": optionCounts["Option 1"]},
			{"name": "Option 2", "count": optionCounts["Option 2"]},
			{"name": "Option 3", "count": optionCounts["Option 3"]},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote result retrieved successfully",
		"data":    result,
	})
}