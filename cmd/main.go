package main

import (
	"database/sql"
	"log"
	"os"

	. "github.com/enkdress/go-todo/pkg/handler"
	. "github.com/enkdress/go-todo/pkg/repository"
	"github.com/enkdress/go-todo/pkg/utils"
	"github.com/labstack/echo/v4"
)

const dbFileName = "todo.db"

func main() {
	os.Remove(dbFileName)

	db, err := sql.Open("sqlite3", dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	taskRepository := NewTaskRepository(db)
	taskHandler := TaskHandler{Repository: taskRepository}
	server := echo.New()
	server.GET(utils.CreateURI("tasks"), taskHandler.GetAll)
	server.POST(utils.CreateURI("tasks"), taskHandler.Create)

	server.Logger.Fatal(server.Start(":3000"))
}
