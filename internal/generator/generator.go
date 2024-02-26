package generator

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"sql-trainer-server/internal/database"
)

type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	TableName   string    `json:"table_name"`
	Difficulty  string    `json:"difficulty"`
	CreatedAt   time.Time `json:"created_at"`
	Solution    string    `json:"-"` // Эталонное решение, не отправляется клиенту
}

type TaskType struct {
	Description     string
	SolutionTpl     string
	Difficulty      string
	RequiredColumns []string
	TableType       database.TableType
}

// Предопределенные типы заданий
var taskTypes = []TaskType{
	{
		Description:     "Найдите всех сотрудников с зарплатой выше средней",
		SolutionTpl:     `SELECT * FROM %s WHERE salary > (SELECT AVG(salary) FROM %s)`,
		Difficulty:      "medium",
		RequiredColumns: []string{"salary"},
		TableType:       database.Employees,
	},
	{
		Description:     "Выведите список из %d сотрудников с самой высокой зарплатой",
		SolutionTpl:     `SELECT * FROM %s ORDER BY salary DESC LIMIT %d`,
		Difficulty:      "easy",
		RequiredColumns: []string{"salary"},
		TableType:       database.Employees,
	},
	{
		Description:     "Посчитайте количество сотрудников в каждом департаменте",
		SolutionTpl:     `SELECT department, COUNT(*) as count FROM %s GROUP BY department ORDER BY count DESC`,
		Difficulty:      "easy",
		RequiredColumns: []string{"department"},
		TableType:       database.Employees,
	},
	{
		Description:     "Найдите среднюю зарплату по каждому департаменту",
		SolutionTpl:     `SELECT department, AVG(salary) as avg_salary FROM %s GROUP BY department ORDER BY avg_salary DESC`,
		Difficulty:      "medium",
		RequiredColumns: []string{"department", "salary"},
		TableType:       database.Employees,
	},
	{
		Description:     "Найдите продукты с рейтингом выше %0.1f",
		SolutionTpl:     `SELECT * FROM %s WHERE rating > %0.1f ORDER BY rating DESC`,
		Difficulty:      "easy",
		RequiredColumns: []string{"rating"},
		TableType:       database.Products,
	},
	{
		Description:     "Посчитайте общую стоимость товаров на складе в каждой категории",
		SolutionTpl:     `SELECT category, SUM(price * stock) as total_value FROM %s GROUP BY category ORDER BY total_value DESC`,
		Difficulty:      "hard",
		RequiredColumns: []string{"category", "price", "stock"},
		TableType:       database.Products,
	},
}

type TaskGenerator struct {
	db *database.Database
}

func NewTaskGenerator(db *database.Database) *TaskGenerator {
	return &TaskGenerator{db: db}
}

func (g *TaskGenerator) GenerateTask() (*Task, error) {
	// Выбираем случайный тип задания
	taskType := taskTypes[rand.Intn(len(taskTypes))]

	// Генерируем уникальное имя таблицы
	tableName := fmt.Sprintf("practice_%d", time.Now().UnixNano())

	// Создаем таблицу с случайными данными
	if err := g.db.CreateRandomTable(tableName, 50); err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	// Генерируем параметры для задания
	var description string
	var solution string

	switch taskType.TableType {
	case database.Employees:
		switch {
		case strings.Contains(taskType.Description, "%d"):
			limit := rand.Intn(5) + 3 // от 3 до 7
			description = fmt.Sprintf(taskType.Description, limit)
			solution = fmt.Sprintf(taskType.SolutionTpl, tableName, limit)
		default:
			description = taskType.Description
			solution = fmt.Sprintf(taskType.SolutionTpl, tableName)
		}
	case database.Products:
		switch {
		case strings.Contains(taskType.Description, "%0.1f"):
			rating := 3.0 + rand.Float64()*1.5 // от 3.0 до 4.5
			description = fmt.Sprintf(taskType.Description, rating)
			solution = fmt.Sprintf(taskType.SolutionTpl, tableName, rating)
		default:
			description = taskType.Description
			solution = fmt.Sprintf(taskType.SolutionTpl, tableName)
		}
	}

	task := &Task{
		ID:          fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Description: description,
		TableName:   tableName,
		Difficulty:  taskType.Difficulty,
		CreatedAt:   time.Now(),
		Solution:    solution,
	}

	return task, nil
}

// ValidateSolution проверяет решение пользователя
func (g *TaskGenerator) ValidateSolution(taskID string, userQuery string) (bool, error) {
	// TODO: Реализовать проверку решения
	// 1. Получить задание по ID
	// 2. Выполнить запрос пользователя
	// 3. Выполнить эталонный запрос
	// 4. Сравнить результаты
	return false, fmt.Errorf("not implemented")
}

// GetTableInfo возвращает информацию о структуре таблицы
func (g *TaskGenerator) GetTableInfo(tableName string) ([]database.ColumnDefinition, error) {
	return g.db.GetTableSchema(tableName)
}
