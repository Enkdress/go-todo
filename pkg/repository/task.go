package repository

import (
	"database/sql"
	"errors"
	"fmt"

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

func (tr *TaskRepository) Update(updatedTask Task) (*Task, error) {
	res := tr.DB.QueryRow(fmt.Sprintf("SELECT uuid, name, description, is_finished, owner_id, due_date, created_at FROM tasks WHERE uuid = '%s'", updatedTask.UUID))
	var task Task

	err := res.Scan(
		&task.UUID,
		&task.Name,
		&task.Description,
		&task.IsFinished,
		&task.OwnerId,
		&task.DueDate,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if task.Name != updatedTask.Name {
		_, err = tr.DB.Exec(fmt.Sprintf("UPDATE tasks SET name = '%s' WHERE uuid = '%s'", updatedTask.Name, updatedTask.UUID))
		if err != nil {
			return nil, err
		}
	}

	if task.Description != updatedTask.Description {
		_, err = tr.DB.Exec(fmt.Sprintf("UPDATE tasks SET description = '%s' WHERE uuid = '%s'", updatedTask.Description, updatedTask.UUID))
		if err != nil {
			return nil, err
		}
	}

	if task.IsFinished != updatedTask.IsFinished {
		_, err = tr.DB.Exec(fmt.Sprintf("UPDATE tasks SET is_finished = %d WHERE uuid = '%s'", updatedTask.IsFinished, updatedTask.UUID))
		if err != nil {
			return nil, err
		}
	}

	return &updatedTask, nil
}

func (tr *TaskRepository) Delete(deletedTask Task) (bool, error) {
	_, err := tr.DB.Exec("DELETE FROM tasks WHERE uuid=?", deletedTask.UUID)

	if err != nil {
		return false, err
	}

	return true, nil
}
