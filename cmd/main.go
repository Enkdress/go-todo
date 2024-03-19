package main

import (
	"database/sql"
	"log"

	. "github.com/enkdress/go-todo/pkg/handler"
	. "github.com/enkdress/go-todo/pkg/repository"
	"github.com/enkdress/go-todo/pkg/utils"
	"github.com/labstack/echo/v4"
)

const dbFileName = "todo.db"

func main() {
	db, err := sql.Open("sqlite3", dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	taskRepository := NewTaskRepository(db)
	taskHandler := TaskHandler{Repository: taskRepository}

	taskRepository.Migrate()

	server := echo.New()
	server.GET(utils.CreateURI("tasks"), taskHandler.GetAll)
	server.POST(utils.CreateURI("tasks"), taskHandler.Create)
	server.PUT(utils.CreateURI("tasks"), taskHandler.Update)
	server.DELETE(utils.CreateURI("tasks"), taskHandler.Delete)

	server.Logger.Fatal(server.Start(":3000"))
}
