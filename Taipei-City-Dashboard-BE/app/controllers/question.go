package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type Question struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

func CreateQuestion(c *gin.Context) {
	// Example of creating a question with auto-generated ID and options
	var req struct {
		Title   string   `json:"title"`
		Options []string `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate auto-generating an ID (in real app, use DB auto-increment or UUID)
	newID := "generated-id-123" // Replace with actual ID generation

	question := map[string]interface{}{
		"id":      newID,
		"title":   req.Title,
		"options": req.Options,
	}

	// Write to db

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

	// Simulate fetching a question
	question := map[string]interface{}{
		"id":      id,
		"title":   "Sample Question",
		"options": []string{"Option 1", "Option 2", "Option 3"},
	}

	// Read from db

	c.JSON(http.StatusOK, gin.H{
		"message": "Question retrieved successfully",
		"data":    question,
	})
}
func GetAllQuestions(c *gin.Context) {
	// Example of getting all questions
	// In a real application, you would fetch all questions from a database

	// Simulate fetching questions
	questions := []map[string]interface{}{
		{
			"id":      "1",
			"title":   "Question 1",
			"options": []string{"Option A", "Option B"},
		},
		{
			"id":      "2",
			"title":   "Question 2",
			"options": []string{"Option X", "Option Y"},
		},
	}

	// Read from db

	c.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
func VoteQuestionByID(c *gin.Context) {
	// 會傳題目 id 和選擇的選項回來
	var req struct {
		ID     string `json:"id"`
		Choice string `json:"choice"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ID == "" || req.Choice == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and choice are required"})
		return
	}

	// Simulate voting
	fmt.Printf("Vote recorded for question ID: %s, choice: %s\n", req.ID, req.Choice)

	// Write to db

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote recorded successfully",
		"data":    map[string]string{"id": req.ID, "choice": req.Choice},
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