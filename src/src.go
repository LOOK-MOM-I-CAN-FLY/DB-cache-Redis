package src

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // Импорт PostgreSQL драйвера
)

const (
	cacheExpiration = time.Minute * 5
)

type App struct {
	DB    *sql.DB
	Cache *redis.Client
	Ctx   context.Context
}

func NewApp(dbConnectionString, redisAddr, redisPassword string) *App {
	ctx := context.Background()

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error to ping the database: %v", err)
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return &App{
		DB:    db,
		Cache: redisClient,
		Ctx:   ctx,
	}
}

func (app *App) QueryWithCache(query string) (string, error) {
	// Проверка данных в кэше
	cachedResult, err := app.Cache.Get(app.Ctx, query).Result()
	if err == nil {
		return formatAsTable(cachedResult), nil
	} else if err != redis.Nil {
		return "", fmt.Errorf("failed to check Redis cache: %v", err)
	}

	// Выполнение SQL-запроса
	rows, err := app.DB.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Получение списка колонок
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to fetch columns: %v", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return "", fmt.Errorf("failed to scan row: %v", err)
		}

		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	// Преобразование результата в JSON для хранения в кэше
	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to cache result: %v", err)
	}

	// Сохранение в кэш
	err = app.Cache.Set(app.Ctx, query, resultJSON, cacheExpiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save to Redis: %v", err)
	}

	// Преобразование результата в текст для вывода
	return formatAsTable(string(resultJSON)), nil
}

func formatAsTable(jsonData string) string {
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Printf("failed to parse JSON: %v", err)
		return "Error formatting data"
	}

	// Форматирование в виде таблицы
	var result string
	if len(data) == 0 {
		return "No results"
	}

	// Получаем имена колонок
	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}

	// Формируем заголовок таблицы
	result += fmt.Sprintf(" %-20s\n", strings.Join(headers, " "))
	result += strings.Repeat("-", 20*len(headers)) + "\n"

	// Добавляем строки
	for _, row := range data {
		var line string
		for _, header := range headers {
			value := fmt.Sprintf("%v", row[header])
			line += fmt.Sprintf(" %-20s", value)
		}
		result += line + "\n"
	}
	return result
}

func (app *App) Close() {
	if err := app.DB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	if err := app.Cache.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}
}
