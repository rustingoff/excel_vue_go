package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/excel_vue_go/internal/models"
	"github.com/rustingoff/excel_vue_go/internal/services"
	"github.com/rustingoff/excel_vue_go/packages/token"
	"io/ioutil"
	"log"
	"net/http"
)

type UserController interface {
	Login(c *gin.Context)

	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserById(c *gin.Context)
	GetUserByEmail(c *gin.Context)
}

type userController struct {
	service services.UserService
}

func GetUserController(service services.UserService) UserController {
	return &userController{service: service}
}

func (controller *userController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("[ERR]: failed binding json to structure, ", err.Error())
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "invalid json structure")
		return
	}

	err := controller.service.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create")
		return
	}

	log.Println("[INF]: user was created successfully")
	c.JSON(http.StatusCreated, "OK")
}

func (controller *userController) DeleteUser(c *gin.Context) {
	var userID = c.Param("id")

	err := controller.service.DeleteUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete")
		return
	}

	log.Println("[INF]: user was deleted successfully")
	c.JSON(http.StatusNoContent, "deleted")
}

func (controller *userController) GetUserById(c *gin.Context) {
	userID := c.Param("id")

	user, err := controller.service.GetUserById(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get one")
		return
	}

	log.Println("[INF]: successfully got user by id")
	c.JSON(http.StatusOK, user)
}

func (controller *userController) GetUserByEmail(c *gin.Context) {
	userEmail := c.Query("email")

	user, err := controller.service.GetUserByEmail(userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get one")
		return
	}

	log.Println("[INF]: successfully got user by email")
	c.JSON(http.StatusOK, user)
}

func (controller *userController) Login(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)

	var input models.SingIn

	if err := json.Unmarshal(jsonData, &input); err != nil {
		log.Println("[ERR]: failed unmarshal json to structure")
		c.AbortWithStatusJSON(http.StatusBadRequest, "provided invalid credentials")
		return
	}

	user, err := controller.service.GetUserByEmail(input.Email)
	if user.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "user not found")
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "user not found")
		return
	}

	ok := token.CheckPasswordHash(input.Password, user.Password)
	if !ok {
		log.Println("[ERR]: provided invalid password")
		c.AbortWithStatusJSON(http.StatusBadRequest, "incorrect password")
		return
	}

	userToken, tokenGenerateError := token.GenerateToken(user.ID, user.Email, user.Active)

	if tokenGenerateError != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed create auth token")
		return
	}

	err = controller.service.Login(input.Email, userToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to identify user")
		return
	}

	c.JSON(http.StatusOK, userToken)
}
