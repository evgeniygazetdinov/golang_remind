package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// TableType определяет тип генерируемой таблицы
type TableType int

const (
	Employees TableType = iota
	Products
	Orders
	Customers
)

// Database представляет подключение к базе данных
type Database struct {
	db *sql.DB
}

// Config содержит параметры подключения к базе данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// ColumnDefinition описывает структуру колонки
type ColumnDefinition struct {
	Name     string
	Type     string
	Nullable bool
}

// TableDefinition описывает структуру таблицы
type TableDefinition struct {
	Type    TableType
	Columns []ColumnDefinition
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(config Config) (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}

	return &Database{db: db}, nil
}

// Предопределенные структуры таблиц
var tableDefinitions = map[TableType]TableDefinition{
	Employees: {
		Type: Employees,
		Columns: []ColumnDefinition{
			{Name: "id", Type: "SERIAL PRIMARY KEY", Nullable: false},
			{Name: "first_name", Type: "VARCHAR(50)", Nullable: false},
			{Name: "last_name", Type: "VARCHAR(50)", Nullable: false},
			{Name: "email", Type: "VARCHAR(100)", Nullable: false},
			{Name: "salary", Type: "DECIMAL(10,2)", Nullable: false},
			{Name: "department", Type: "VARCHAR(50)", Nullable: false},
			{Name: "hire_date", Type: "DATE", Nullable: false},
		},
	},
	Products: {
		Type: Products,
		Columns: []ColumnDefinition{
			{Name: "id", Type: "SERIAL PRIMARY KEY", Nullable: false},
			{Name: "name", Type: "VARCHAR(100)", Nullable: false},
			{Name: "category", Type: "VARCHAR(50)", Nullable: false},
			{Name: "price", Type: "DECIMAL(10,2)", Nullable: false},
			{Name: "stock", Type: "INTEGER", Nullable: false},
			{Name: "rating", Type: "DECIMAL(3,2)", Nullable: true},
			{Name: "created_at", Type: "TIMESTAMP", Nullable: false},
		},
	},
}

// Данные для генерации
var (
	firstNames        = []string{"John", "Alice", "Bob", "Emma", "Michael", "Sarah", "David", "Lisa"}
	lastNames         = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller"}
	departments       = []string{"IT", "HR", "Sales", "Marketing", "Finance", "Operations"}
	productCategories = []string{"Electronics", "Books", "Clothing", "Food", "Sports", "Home"}
	productNames      = []string{"Laptop", "Smartphone", "Headphones", "Camera", "Tablet", "Watch"}
)

func (d *Database) CreateRandomTable(tableName string, rowCount int) error {
	// Выбираем случайный тип таблицы
	tableTypes := []TableType{Employees, Products}
	tableType := tableTypes[rand.Intn(len(tableTypes))]
	tableDef := tableDefinitions[tableType]

	// Создаем таблицу
	columns := make([]string, len(tableDef.Columns))
	for i, col := range tableDef.Columns {
		if col.Nullable {
			columns[i] = fmt.Sprintf("%s %s", col.Name, col.Type)
		} else {
			columns[i] = fmt.Sprintf("%s %s NOT NULL", col.Name, col.Type)
		}
	}

	createTableSQL := fmt.Sprintf("CREATE TABLE %s (%s)",
		tableName,
		strings.Join(columns, ", "),
	)

	if _, err := d.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Генерируем данные
	for i := 0; i < rowCount; i++ {
		var err error
		switch tableType {
		case Employees:
			err = d.insertRandomEmployee(tableName)
		case Products:
			err = d.insertRandomProduct(tableName)
		}
		if err != nil {
			return fmt.Errorf("error inserting data: %v", err)
		}
	}

	return nil
}

func (d *Database) insertRandomEmployee(tableName string) error {
	sql := fmt.Sprintf(`
        INSERT INTO %s (first_name, last_name, email, salary, department, hire_date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, tableName)

	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	email := strings.ToLower(fmt.Sprintf("%s.%s@example.com", firstName, lastName))
	salary := 30000 + rand.Float64()*120000
	department := departments[rand.Intn(len(departments))]
	hireDate := time.Now().AddDate(-rand.Intn(5), -rand.Intn(12), -rand.Intn(28))

	_, err := d.db.Exec(sql,
		firstName,
		lastName,
		email,
		salary,
		department,
		hireDate.Format("2006-01-02"),
	)
	return err
}

func (d *Database) insertRandomProduct(tableName string) error {
	sql := fmt.Sprintf(`
        INSERT INTO %s (name, category, price, stock, rating, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, tableName)

	name := productNames[rand.Intn(len(productNames))]
	category := productCategories[rand.Intn(len(productCategories))]
	price := 10.0 + rand.Float64()*990.0
	stock := rand.Intn(1000)
	rating := 1.0 + rand.Float64()*4.0
	createdAt := time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(28))

	_, err := d.db.Exec(sql,
		name,
		category,
		price,
		stock,
		rating,
		createdAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

// GetTableSchema возвращает схему таблицы
func (d *Database) GetTableSchema(tableName string) ([]ColumnDefinition, error) {
	query := `
        SELECT column_name, data_type, is_nullable
        FROM information_schema.columns
        WHERE table_name = $1
        ORDER BY ordinal_position
    `

	rows, err := d.db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnDefinition
	for rows.Next() {
		var col ColumnDefinition
		var isNullable string
		if err := rows.Scan(&col.Name, &col.Type, &isNullable); err != nil {
			return nil, err
		}
		col.Nullable = isNullable == "YES"
		columns = append(columns, col)
	}

	return columns, nil
}
