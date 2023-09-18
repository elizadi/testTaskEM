package delivery

import (
	"effective_mobile/internal/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	
)

type User struct {
	ID         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        uint8
	Gender     string
	Country    string
}

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
		id := ctx.Query("id")
		idInt, _ := strconv.Atoi(id)
		name := ctx.Query("name")
		surname := ctx.Query("surname")
		patronymic := ctx.Query("patronymic")
		age := ctx.Query("age")
		ageInt, _ := strconv.Atoi(age)
		gender := ctx.Query("gender")
		country := ctx.Query("country")
		page := ctx.Query("page")
		pageInt, _ := strconv.Atoi(page)
		perPage := ctx.Query("perPage")
		perPageInt, _ := strconv.Atoi(perPage)

		req := domain.GetUsersReq{
			ID:         uint64(idInt),
			Name:       name,
			Surname:    surname,
			Patronymic: patronymic,
			Age:        uint8(ageInt),
			Gender:     gender,
			Country:    country,
			Pag:        domain.Pagination{
				Page:    pageInt,
				PerPage: perPageInt,
			},
		}

		users, err := uc.GetUsersWithPagination(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, users)
	})

	router.DELETE("/user", func(ctx *gin.Context) {
		id := ctx.Query("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		err = uc.DeleteUser(uint64(idInt))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, "Success")
	})

	router.PUT("/user", func(ctx *gin.Context) {
		var user User 

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedUser, err := uc.UpdateUser(user.ID, domain.User{
			ID:         user.ID,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Age:        user.Age,
			Gender:     user.Gender,
			Country:    user.Country,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "")
			return
		}
		ctx.JSON(http.StatusOK, updatedUser)
	})
}
