package handler

import (
	"log"
	"net/http"

	. "github.com/enkdress/go-todo/pkg/model"
	. "github.com/enkdress/go-todo/pkg/repository"
	"github.com/enkdress/go-todo/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	Repository *TaskRepository
}

func (ht TaskHandler) GetAll(c echo.Context) error {
	repository := ht.Repository
	if err := repository.Migrate(); err != nil {
		log.Fatal(err)
		return err
	}
	allTasks, err := repository.All()
	res := utils.CreateReturnObject[Task](allTasks)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
	return nil
}

func (ht TaskHandler) Create(c echo.Context) error {
	repository := ht.Repository
	var task Task
	err := c.Bind(&task)
	task.UUID = uuid.NewString()

	if err != nil {
		log.Fatal(err)
		return err
	}

	createdTask, err := repository.Create(task)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.JSON(http.StatusOK, createdTask)
	return nil
}

func (ht TaskHandler) Update(c echo.Context) error {
	repository := ht.Repository
	var task Task

	err := c.Bind(&task)
	if err != nil {
		log.Fatal(err)
		return err
	}

	uTask, err := repository.Update(task)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.JSON(http.StatusOK, uTask)

	return nil
}

func (ht TaskHandler) Delete(c echo.Context) error {
	repository := ht.Repository
	var task Task

	err := c.Bind(&task)
	if err != nil {
		log.Fatal(err)
		return err
	}

	isDeleted, err := repository.Delete(task)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.JSON(http.StatusOK, utils.CreateReturnMessage(isDeleted))

	return nil
}
