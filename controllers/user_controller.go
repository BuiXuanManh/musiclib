package controllers

import (
	"musiclib/helper"
	"musiclib/models"
	"musiclib/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserService services.UserService
}

type ResetPassword struct {
	Username    string `form:"username" json:"username" binding:"required"`
	OldPassword string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required"`
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// CreateUser 	godoc
// @Summary      CreateUser
// @Description  create a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body     dto.UserDto  true  "User data to create"
// @Router       /user/create [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if user.Username == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to hash the password"})
	}
	user.Password = hashPassword

	if err := uc.UserService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

// GetUser 	godoc
// @Summary      GetUser
// @Description  Get a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Find by User ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}   models.User
// @Router       /user/get/{id} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := uc.UserService.GetUser(&userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser 	godoc
// @Summary      UpdateUser
// @Description  Update a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body     models.User  true  "User data to update"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /user/update [patch]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if user.UserId == "" || user.Username == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong input structure"})
		return
	}

	if err := uc.UserService.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

// ChangePass 	godoc
// @Summary      ChangePass
// @Description  change pass a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        changePass   body     controllers.ResetPassword  true  "User data to change password"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /user/change_password [patch]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
	var resetPassword ResetPassword

	if err := ctx.ShouldBindJSON(&resetPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := uc.UserService.ChangePassword(&resetPassword.Username, &resetPassword.OldPassword, &resetPassword.NewPassword); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})

}

// DeleteUser 	godoc
// @Summary      DeleteUser
// @Description  delete a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Delete by User ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router       /user/delete/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	err = uc.UserService.DeleteUser(&userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (uc *UserController) RegisterUserRoute(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user")
	// The URI must be diffent structure from each other !
	userRoute.POST("/create", uc.CreateUser)

	userRoute.GET("/get/:id", uc.GetUser)

	userRoute.PATCH("/update", uc.UpdateUser)

	userRoute.PATCH("/change_password", uc.ChangePassword)

	userRoute.DELETE("/delete/:id", uc.DeleteUser)
}
