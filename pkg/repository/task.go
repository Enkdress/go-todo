package repository

import (
	"database/sql"
	"errors"

	. "github.com/enkdress/go-todo/pkg/model"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
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
				return nil, ErrDuplicate
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
	res, err := tr.DB.Exec("UPDATE tasks SET name = ?, description = ?, is_finished = ? WHERE uuid = ?", updatedTask.Name, updatedTask.Description, updatedTask.IsFinished, updatedTask.UUID)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updatedTask, nil
}

func (tr *TaskRepository) Delete(deletedTask Task) (bool, error) {
	res, err := tr.DB.Exec("DELETE FROM tasks WHERE uuid=?", deletedTask.UUID)

	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, ErrDeleteFailed
	}

	return true, nil
}
