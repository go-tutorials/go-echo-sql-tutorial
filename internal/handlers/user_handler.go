package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	sv "github.com/core-go/service"
	"github.com/labstack/echo/v4"

	. "go-service/internal/models"
	. "go-service/internal/services"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c echo.Context) error {
	result, err := h.service.GetAll(context.Background())
	if err != nil {
		c.Error(err)
		return nil
	}
	return c.JSON(http.StatusOK, result)
}

func (h *UserHandler) Load(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		log.Fatalf("user id cannot be empty")
	}

	result, err := h.service.Load(context.Background(), id)
	if err != nil {
		c.Error(err)
		return nil
	}
	return c.JSON(http.StatusOK, result)
}

func (h *UserHandler) Insert(c echo.Context) error {
	var user User
	er1 := c.Bind(&user)
	defer c.Request().Body.Close()
	if er1 != nil {
		c.Error(er1)
		return nil
	}

	_, er2 := h.service.Insert(context.Background(), &user)
	if er2 != nil {
		c.Error(er2)
		return nil
	}

	msg := fmt.Sprintf("new user with id '%s' has been created", user.Id)
	return c.JSON(http.StatusCreated, msg)
}

func (h *UserHandler) Update(c echo.Context) error {
	var user User
	er1 := c.Bind(&user)
	defer c.Request().Body.Close()
	if er1 != nil {
		c.Error(er1)
		return nil
	}
	id := c.Param("id")
	if len(id) == 0 {
		log.Fatalf("user id cannot be empty")
		return nil
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		log.Fatalf("user id is not match")
		return nil
	}

	_, er2 := h.service.Update(context.Background(), &user)
	if er2 != nil {
		c.Error(er2)
		return nil
	}

	msg := fmt.Sprintf("user with id '%s' has been updated", id)
	return c.JSON(http.StatusOK, msg)
}

func (h *UserHandler) Patch(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		log.Fatalf("user id cannot be empty")
		return nil
	}

	ids := []string{"id"}

	var user User
	userType := reflect.TypeOf(user)
	_, jsonMap := sv.BuildMapField(userType)
	body, _ := sv.BuildMapAndStruct(c.Request(), &user)
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		log.Fatalf("user id is not match")
		return nil
	}
	json, er1 := sv.BodyToJson(c.Request(), user, body, ids, jsonMap, nil)
	if er1 != nil {
		c.Error(er1)
		return nil
	}

	_, er2 := h.service.Patch(context.Background(), json)
	if er2 != nil {
		c.Error(er2)
		return nil
	}

	msg := fmt.Sprintf("user with id '%s' has been patched", id)
	return c.JSON(http.StatusOK, msg)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		log.Fatalf("user id cannot be empty")
		return nil
	}
	_, err := h.service.Delete(context.Background(), id)
	if err != nil {
		c.Error(err)
		return nil
	}

	msg := fmt.Sprintf("user with id '%s' has been removed", id)
	return c.JSON(http.StatusOK, msg)
}

// func respond(w http.ResponseWriter, result interface{}) {
// 	response, _ := json.Marshal(result)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }
