package handler

import (
	"avito/internal/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInfo(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok || username == "" {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	usrname := username.(string)
	output, err := h.services.Coins.GetInfo(usrname)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) SendCoins(c *gin.Context) {
	username, ok := c.Get("username")
	usrname := username.(string)
	if !ok || username == "" {
		fmt.Println(username)
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	var input types.SendCoinRequest

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid request")
		return
	}

	err = h.services.Send(usrname, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (h *Handler) BuyItem(c *gin.Context) {
	username, ok := c.Get("username")
	usrname := username.(string)
	if !ok || username == "" {
		fmt.Println(username)
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	item := c.Param("item")
	if item == "" {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	err := h.services.BuyItem(usrname, item)

}
