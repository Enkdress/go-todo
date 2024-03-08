package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	. "github.com/enkdress/go-todo/pkg/handler"
	. "github.com/enkdress/go-todo/pkg/repository"
	"github.com/labstack/echo/v4"
)

func buildResourceURI(resource string) string {
	const V1URI = "/v1"
	return fmt.Sprintf("%s/%s", V1URI, resource)
}

const dbFileName = "tasks.db"

func main() {
	os.Remove(dbFileName)

	db, err := sql.Open("sqlite3", dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	taskRepository := NewTaskRepository(db)
	taskHandler := TaskHandler{Repository: taskRepository}

	server := echo.New()
	server.GET(buildResourceURI("tasks"), taskHandler.GetAll)
	server.POST(buildResourceURI("tasks"), taskHandler.Create)

	server.Logger.Fatal(server.Start(":3000"))
}
