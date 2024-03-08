package repository

import (
	"database/sql"
	"errors"

	. "github.com/enkdress/go-todo/pkg/model"
	"github.com/mattn/go-sqlite3"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		DB: db,
	}
}

func (tr *TaskRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT,
		is_finished INTEGER DEFAULT 0,
		due_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		owner_id INTEGER DEFAULT 1,
		assignee_id INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := tr.DB.Exec(query)
	return err
}

func (tr *TaskRepository) Create(task Task) (*Task, error) {
	res, err := tr.DB.Exec("INSERT INTO tasks (uuid, name, description) VALUES (?,?,?)", task.UUID, task.Name, task.Description)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, errors.New("record already exists")
			}
		}

		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	task.ID = int(id)
	return &task, nil
}

func (tr *TaskRepository) All() ([]Task, error) {
	rows, err := tr.DB.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.UUID,
			&task.Name,
			&task.Description,
			&task.IsFinished,
			&task.DueDate,
			&task.OwnerId,
			&task.AssigneeId,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		all = append(all, task)
	}

	return all, nil
}
