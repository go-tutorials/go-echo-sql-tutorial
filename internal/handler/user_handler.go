package handler

import (
	sv "github.com/core-go/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"

	. "go-service/internal/model"
	. "go-service/internal/service"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c echo.Context) error {
	res, err := h.service.All(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Load(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	res, err := h.service.Load(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Insert(c echo.Context) error {
	var user User
	er1 := c.Bind(&user)

	defer c.Request().Body.Close()
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	res, er2 := h.service.Insert(c.Request().Context(), &user)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Update(c echo.Context) error {
	var user User
	er1 := c.Bind(&user)
	defer c.Request().Body.Close()

	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}

	res, er2 := h.service.Update(c.Request().Context(), &user)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Patch(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	r := c.Request()
	var user User
	userType := reflect.TypeOf(user)
	_, jsonMap, _ := sv.BuildMapField(userType)
	body, er0 := sv.BuildMapAndStruct(r, &user)
	if er0 != nil {
		return c.String(http.StatusInternalServerError, er0.Error())
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}
	json, er1 := sv.BodyToJsonMap(r, user, body, []string{"id"}, jsonMap)
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	res, er2 := h.service.Patch(r.Context(), json)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	res, err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
