package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// todo: tangkap input dari user
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Register account failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse(
			"Register account failed",
			http.StatusBadRequest,
			"error",
			nil,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "tokentokentokentoken")

	response := helper.APIResponse(
		"Account has been registered",
		http.StatusOK,
		"success",
		formatter,
	)

	c.JSON(http.StatusOK, response)

	// * map input dari user ke struct RegsiterUserInput
	// * struct di atas kita passing sebagai parameter service

}

func (h *userHandler) Login(c *gin.Context) {
	/*
	* 1. User memasukkan input (email & password)
	* 2. Input ditangkap handler
	* 3. Mapping dari input user ke input struct
	* 4. Input struct passing service
	* 5. Di service mencari dengan bantuan repository user dengan email X
	* 6. Mencocokkan password
	 */

	// * [1]
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Login failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			"Login failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentokenloginlogin")

	response := helper.APIResponse(
		"Succesfully loggedin",
		http.StatusOK,
		"success",
		formatter,
	)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	/*
	* 1. Ada input email dari user
	* 2. input email di-mapping ke struct input
	* 3. struct input di-passing ke service
	* 4. service akan memanggil repository - email sudah digunakan atau belum
	* 5. repository - db
	 */

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Email checking failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {

		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse(
			"Email checking failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(
		metaMessage,
		http.StatusOK,
		"success",
		data,
	)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	/*
		1. Input dari user
		2. Simpan gambarnya di folder "images/"
		3. di service kita panggil repo
		4. JWT -> (Sementara hardcore, seakan2 user yg login ID = 1)
		5. Repo ambil data user yg ID = 1
		6. repo update data user simpan lokasi file
	*/

	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := 1 // note : harusnya dapet dari JWT

	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/avatar-img-%d.png", userID)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			"error",
			data,
		)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}

	response := helper.APIResponse(
		"Avatar image successfully uploaded",
		http.StatusOK,
		"success",
		data,
	)

	c.JSON(http.StatusOK, response)

}
