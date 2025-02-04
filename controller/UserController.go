package controller

import (
	model "tutor_rest/internal/model"
	// "tutor_rest/internal/class"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func AddUser(c echo.Context) error {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     int    `json:"role"`
		Status   int    `json:"status"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Username: %s", requestBody.Username)
	log.Printf("Received Password: %s", requestBody.Password)
	log.Printf("Received Role: %d", requestBody.Role)
	log.Printf("Received Status: %d", requestBody.Status)

	if requestBody.Username == "" || requestBody.Password == "" || requestBody.Role == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username, password and role are required"})
	}

	addResult, err := model.AddUser(requestBody.Username, requestBody.Password, requestBody.Role, requestBody.Status)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User added successfully"})
}

func LoginUser(c echo.Context) error {

	username := c.FormValue("username") // This can be either user_name or email
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username/Email and password are required"})
	}

	loginResult, err := model.LoginUser(username, password)
	if err != nil {
		return c.JSON(loginResult.Status, map[string]string{"message": loginResult.Message})
	}

	data, ok := loginResult.Data.(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to extract data"})
	}
	apiToken, ok := data["token"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve API token"})
	}

	userID, ok := data["user_id"]
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve User ID"})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   http.StatusOK,
		"message":  "Login successful",
		"apiToken": apiToken,
		"userID":   userID,
	})
}

func GetAllUsers(c echo.Context) error {
	usersResult, err := model.GetAllUsers() // Assuming this model function exists and fetches all user data
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve users"})
	}

	if usersResult.Status != http.StatusOK {
		return c.JSON(usersResult.Status, map[string]string{"message": usersResult.Message})
	}

	return c.JSON(http.StatusOK, usersResult.Data)
}

func GetUserByID(c echo.Context) error {
	userIDParam := c.Param("id")
	// c.QueryParam("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	user, err := model.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve user"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser handles the HTTP request to update a user's details.
func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var requestBody struct {
		Username string `json:"username"`
		Role     int    `json:"role"`
		Status   int    `json:"status"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.UpdateUser(id, requestBody.Username, requestBody.Role, requestBody.Status)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles the HTTP request to delete a user.
func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	deleteResult, err := model.DeleteUser(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func GetUserTransactions(c echo.Context) error {
	userIDParam := c.Param("id")
	// c.QueryParam("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	user, err := model.GetTransactionByUserId(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve user"})
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUserStatus(c echo.Context) error {
	var requestBody struct {
		Status string `json:"status"`
		Detail string `json:"detail"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Status: %s", requestBody.Status)
	log.Printf("Received Detail: %s", requestBody.Detail)

	if requestBody.Status == "" || requestBody.Detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Status and Detail are required"})
	}

	addResult, err := model.AddUserStatus(requestBody.Status, requestBody.Detail)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User status added successfully"})
}

func GetUserStatusById(c echo.Context) error {
	transactionIDParam := c.Param("id")
	// c.QueryParam("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user status ID"})
	}

	dtTransaction, err := model.GetUserStatusById(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transaction"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User status not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func GetAllUserStatus(c echo.Context) error {
	dtTransaction, err := model.GetAllUserStatus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve user status"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User status not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func UpdateUserStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user status ID"})
	}

	var requestBody struct {
		Status string `json:"status"`
		Detail string `json:"detail"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditUserStatus(id, requestBody.Status, requestBody.Detail)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User status updated successfully"})
}

func DeleteUserStatusById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user status ID"})
	}

	deleteResult, err := model.DeleteUserStatusById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User status deleted successfully"})
}

func CreateUserRole(c echo.Context) error {
	var requestBody struct {
		Role   string `json:"role"`
		Detail string `json:"detail"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Status: %s", requestBody.Role)
	log.Printf("Received Detail: %s", requestBody.Detail)

	if requestBody.Role == "" || requestBody.Detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Role and Detail are required"})
	}

	addResult, err := model.AddUserRole(requestBody.Role, requestBody.Detail)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User role added successfully"})
}

func GetUserRoleById(c echo.Context) error {
	transactionIDParam := c.Param("id")
	// c.QueryParam("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user role ID"})
	}

	dtTransaction, err := model.GetUserRoleById(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transaction"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User role not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func GetAllUserRole(c echo.Context) error {
	dtTransaction, err := model.GetAllUserRole()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve user role"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User role not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func UpdateUserRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user role ID"})
	}

	var requestBody struct {
		Role   string `json:"role"`
		Detail string `json:"detail"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditUserRole(id, requestBody.Role, requestBody.Detail)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User role updated successfully"})
}

func DeleteUserRoleById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user role ID"})
	}

	deleteResult, err := model.DeleteUserRoleById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User role deleted successfully"})
}
