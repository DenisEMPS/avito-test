package handler

import (
	"avito/internal/service"
	"avito/internal/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	var input types.UserCreate

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	id, err := h.services.RegisterNewUser(input)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			NewErrorResponse(c, http.StatusConflict, "user already exists")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusCreated, map[string]int64{
		"id": id,
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	var input types.UserLoginDTO

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	token, err := h.services.Authorization.LoginUser(input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			NewErrorResponse(c, http.StatusUnauthorized, "invalid request params")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
