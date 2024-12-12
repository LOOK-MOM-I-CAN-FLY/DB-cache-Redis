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
	_ "github.com/lib/pq"
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


	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error to ping the database: %v", err)
	}


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

	cachedResult, err := app.Cache.Get(app.Ctx, query).Result()
	if err == nil {
		return formatAsTable(cachedResult), nil
	} else if err != redis.Nil {
		return "", fmt.Errorf("failed to check Redis cache: %v", err)
	}


	rows, err := app.DB.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()


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


	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to cache result: %v", err)
	}


	err = app.Cache.Set(app.Ctx, query, resultJSON, cacheExpiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save to Redis: %v", err)
	}


	return formatAsTable(string(resultJSON)), nil
}

func formatAsTable(jsonData string) string {
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Printf("failed to parse JSON: %v", err)
		return "Error formatting data"
	}


	var result string
	if len(data) == 0 {
		return "No results"
	}


	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}


	result += fmt.Sprintf(" %-20s\n", strings.Join(headers, " "))
	result += strings.Repeat("-", 20*len(headers)) + "\n"


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
