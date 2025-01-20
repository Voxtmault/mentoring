package controller

import (
	model "tutor_rest/internal/model"
	// "tutor_rest/internal/class"
	"net/http"
	"github.com/labstack/echo/v4"
    "log"
    "strconv"
)

func AddUser(c echo.Context) error {
    var requestBody struct {
		Username   string `json:"username"`
		Password string `json:"password"`
        Role     string `json:"role"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Username: %s", requestBody.Username)
	log.Printf("Received Password: %s", requestBody.Password) // Log the qr_log_id
	log.Printf("Received Role: %s", requestBody.Role) // Log the qr_log_id

    if requestBody.Username == "" || requestBody.Password == "" || requestBody.Role == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username, password and role are required"})
    }

    addResult, err := model.AddUser(requestBody.Username, requestBody.Password, requestBody.Role)
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
        "status":    http.StatusOK,
        "message":   "Login successful",
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
        Role     string `json:"role"`
    }
    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
    }

    updateResult, err := model.UpdateUser(id, requestBody.Username, requestBody.Role)
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
    return nil
}
func ListSuppliers(c echo.Context) error {
    return nil
}
func ViewStock(c echo.Context) error {
    return nil
}
func UpdateStock(c echo.Context) error {
    return nil
}
func AddTransaction(c echo.Context) error {
    return nil
}
func AddPurchaseFromSupplier(c echo.Context) error {
    return nil
}