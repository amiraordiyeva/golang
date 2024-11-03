package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
    "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ключ для подписи токенов
var jwtKey = []byte("my_secret_key")

// cтруктура для пользователя
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// структура для заданий
type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// структура для хранения клеймов
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// списки задач и пользователей
var tasks []Task
var users = []User{
	{Username: "admin", Password: "123", Role: "admin"},
	{Username: "user", Password: "123", Role: "user"},
}

// для уникального ID задач
var currentID = 1

func main() {
	router := gin.Default()
	// конфигурация CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// обрабатка входа
	router.POST("/login", Login)

	// создание защищённых маршрутов
	protected := router.Group("/")
	protected.Use(Authorization())
	{
		protected.GET("/tasks", GetTasks)
		protected.POST("/tasks", CreateTask)
		protected.PUT("/tasks/:id", UpdateTask)
		protected.DELETE("/tasks/:id", DeleteTask)
	}

	router.Run(":8080") // запуск сервера
}

// функция для входа
func Login(c *gin.Context) {
	var credentials User
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &Claims{
				Username: user.Username,
				Role:     user.Role,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			// создание токена
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": tokenString}) // отправка токен клиенту
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}) // ошибка если неверные данные
}

// функция для получения задач
func GetTasks(c *gin.Context) {
	username := c.MustGet("username").(string)
	userTasks := []Task{}
	for _, task := range tasks {
		if task.Username == username {
			userTasks = append(userTasks, task) // добавление задачи пользователя в список
		}
	}
	c.JSON(http.StatusOK, userTasks) // отправка задачи пользователю
}

// для создания задачи
func CreateTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = currentID
	currentID++
	task.Username = c.MustGet("username").(string) // сохранение имя пользователя
	tasks = append(tasks, task)                    // добавление задачи в список
	c.JSON(http.StatusOK, task)                    // отправка созданной задачи
}

// для обновления задачи
func UpdateTask(c *gin.Context) {
	var updatedTask Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	for i, task := range tasks {
		if fmt.Sprint(task.ID) == id {
			if task.Username != c.MustGet("username").(string) {
				c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
				return
			}
			tasks[i].Name = updatedTask.Name // обновление имени задачи
			c.JSON(http.StatusOK, tasks[i])  // отправка обновлённой задачи
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"}) // если не найдена
}

// для авторизации
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		// удаление bearer из токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username) // cохранение имя пользователя в контексте
		c.Set("role", claims.Role)         // сохранение роль пользователя в контексте
		c.Next()
	}
}

// для удаления задачи
func DeleteTask(c *gin.Context) {
	id := c.Param("id") // получение ID таска из параметров
	for i, task := range tasks {
		if fmt.Sprint(task.ID) == id {
			if task.Username != c.MustGet("username").(string) {
				c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
				return
			}
			tasks = append(tasks[:i], tasks[i+1:]...)               // удаление задачи
			c.JSON(http.StatusOK, gin.H{"message": "task deleted"}) // отправка сообщение об успешном удалении
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"}) // если задача не найдена
}
