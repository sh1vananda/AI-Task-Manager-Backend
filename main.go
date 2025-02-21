package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ai-task-manager-backend/controllers"
	"ai-task-manager-backend/middleware"
	"ai-task-manager-backend/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open("task_manager.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Auto-migrate the database schema
	db.AutoMigrate(&models.User{}, &models.Task{})

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ai-task-manager-frontend-ee9twrfmy-llamatypes-projects.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db))
	r.POST("/tasks", middleware.AuthMiddleware(), controllers.CreateTask(db))
	r.GET("/tasks", middleware.AuthMiddleware(), controllers.GetTasks(db))

	// AI-Powered Task Suggestions
	r.POST("/suggest-task", middleware.AuthMiddleware(), func(c *gin.Context) {
		var input struct {
			Goal string `json:"goal" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call OpenAI API
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("OpenAI API key not found")
		}

		url := "https://api.openai.com/v1/chat/completions"
		payload := map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": fmt.Sprintf("Suggest actionable tasks for the following goal: %s", input.Goal),
				},
			},
			"max_tokens": 100,
		}

		payloadBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error calling OpenAI API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call OpenAI API"})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading OpenAI API response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read OpenAI API response"})
			return
		}

		var openAIResponse map[string]interface{}
		if err := json.Unmarshal(body, &openAIResponse); err != nil {
			log.Printf("Error parsing OpenAI API response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OpenAI API response"})
			return
		}

		choices := openAIResponse["choices"].([]interface{})
		if len(choices) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No suggestions from OpenAI"})
			return
		}

		suggestion := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		c.JSON(http.StatusOK, gin.H{"suggestions": suggestion})
	})

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", message)
			conn.WriteMessage(messageType, message)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local development
	}
	r.Run(":" + port)
}
