package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "sql-trainer-server/docs" // Этот импорт будет создан автоматически
	"sql-trainer-server/internal/database"
	"sql-trainer-server/internal/generator"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	db            *database.Database
	taskGenerator *generator.TaskGenerator
)

func init() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Инициализация базы данных
	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	var err error
	db, err = database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	taskGenerator = generator.NewTaskGenerator(db)
}

// @title SQL Trainer API
// @version 1.0
// @description API для тренировки SQL запросов
// @host localhost:8080
// @BasePath /api
func main() {
	r := mux.NewRouter()

	// // Применяем CORS middleware
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Определяем маршруты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/generate-task", generateTaskHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("/check-solution", checkSolutionHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on http://0.0.0.0:%s", port)
	log.Printf("Swagger UI available at http://0.0.0.0:%s/swagger/index.html", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// @Summary Генерация нового задания
// @Description Создает новую случайную таблицу и генерирует задание для неё
// @Tags tasks
// @Produce json
// @Success 200 {object} generator.Task
// @Router /generate-task [get]
func generateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := taskGenerator.GenerateTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// @Summary Проверка решения
// @Description Проверяет правильность SQL запроса для заданного задания
// @Tags tasks
// @Accept json
// @Produce json
// @Param solution body CheckSolutionRequest true "Решение задания"
// @Success 200 {object} CheckSolutionResponse
// @Router /check-solution [post]
func checkSolutionHandler(w http.ResponseWriter, r *http.Request) {
	var solution CheckSolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&solution); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Проверка решения
	w.WriteHeader(http.StatusOK)
}

// @Summary Проверка здоровья сервиса
// @Description Возвращает статус работоспособности сервиса
// @Tags system
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(HealthCheckResponse{Status: "ok"})
}

// Структуры для документации API
type CheckSolutionRequest struct {
	TaskID string `json:"task_id" example:"task_1234567890"`
	Query  string `json:"query" example:"SELECT * FROM employees WHERE salary > 50000"`
}

type CheckSolutionResponse struct {
	Correct bool   `json:"correct" example:"true"`
	Message string `json:"message,omitempty" example:"Решение верное!"`
}

type HealthCheckResponse struct {
	Status string `json:"status" example:"ok"`
}
