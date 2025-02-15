package handler

import (
	"avito/internal/service"
	"avito/internal/types"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInfo(c *gin.Context) {
	param, ok := c.Get("username")
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	username, ok := param.(string)
	if !ok || username == "" {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}
	output, err := h.services.Coins.GetInfo(username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			NewErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) SendCoins(c *gin.Context) {
	param, ok := c.Get("username")
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	username, ok := param.(string)
	if !ok || username == "" {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	var input types.SendCoinRequest

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.services.Send(username, input)
	if err != nil {
		if errors.Is(err, service.ErrRecieverNotFound) {
			NewErrorResponse(c, http.StatusBadRequest, "user not found")
			return
		} else if errors.Is(err, service.ErrNotEnoughtCoins) {
			NewErrorResponse(c, http.StatusBadRequest, "not enought coins")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (h *Handler) BuyItem(c *gin.Context) {
	param, ok := c.Get("username")
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	username, ok := param.(string)
	if !ok || username == "" {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	fmt.Println(username)

	item := c.Param("item")
	if item == "" {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	input := new(types.BuyRequest)

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	err = h.services.BuyItem(username, item, input)
	if err != nil {
		if errors.Is(err, service.ErrNotEnoughtCoins) {
			NewErrorResponse(c, http.StatusBadRequest, "not enought coins")
			return
		} else if errors.Is(err, service.ErrItemNotFound) {
			NewErrorResponse(c, http.StatusBadRequest, "item not found")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
