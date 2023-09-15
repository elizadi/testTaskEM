package delivery

import (
	"effective_mobile/internal/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(uc domain.UseCase, router *gin.Engine) {
	router.POST("/user", func(ctx *gin.Context) {
		name := ctx.Query("name")
		if name == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("empty parameter"))
			return
		}
		surname := ctx.Query("surname")
		if surname == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("empty parameter"))
			return
		}
		patronymic := ctx.Query("patronymic")
		if patronymic == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("empty parameter"))
			return
		}
		user, err := uc.CreateUser(name, surname, patronymic)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, user)
	})

	router.GET("/user", func(ctx *gin.Context) {
		users, err := uc.GetUsers()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, users)
	})
	router.GET("/pagination", func(ctx *gin.Context) {
		page := ctx.Query("page")
		if page == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("empty parameter"))
			return
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		perPage := ctx.Query("perPage")
		if perPage == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("empty parameter"))
			return
		}
		perPageInt, err := strconv.Atoi(perPage)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		users, err := uc.GetUsersWithPagination(uint(pageInt), uint(perPageInt))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, users)
	})
}
