package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"video-agent-go/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected successfully")
}

func SaveTask(taskID, input, output string) error {
	query := `INSERT INTO video_tasks (task_id, input, output, created_at) VALUES (?, ?, ?, NOW())`
	_, err := DB.Exec(query, taskID, input, output)
	return err
}

func GetTask(taskID string) (*VideoTask, error) {
	query := `SELECT id, task_id, input, output, created_at FROM video_tasks WHERE task_id = ?`
	row := DB.QueryRow(query, taskID)

	var task VideoTask
	err := row.Scan(&task.ID, &task.TaskID, &task.Input, &task.Output, &task.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func GetAllTasks() ([]VideoTask, error) {
	query := `SELECT id, task_id, input, output, created_at FROM video_tasks ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []VideoTask
	for rows.Next() {
		var task VideoTask
		err := rows.Scan(&task.ID, &task.TaskID, &task.Input, &task.Output, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func UpdateTaskOutput(taskID string, output interface{}) error {
	outputJSON, _ := json.Marshal(output)
	query := `UPDATE video_tasks SET output = ? WHERE task_id = ?`
	_, err := DB.Exec(query, string(outputJSON), taskID)
	return err
}
