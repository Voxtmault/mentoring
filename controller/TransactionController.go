package controller

import (
	// "tutor_rest/internal/model"
	// "tutor_rest/internal/class"
	"log"
	"net/http"
	"strconv"
	model "tutor_rest/internal/model"

	"github.com/labstack/echo/v4"
)

func CreateTransaction(c echo.Context) error {
	var requestBody struct {
		UserId  int   `json:"nama"`
		Total   int   `json:"total"`
		Status  int   `json:"status"`
		StockID []int `json:"stock_id"`
		Jumlah  []int `json:"jumlah"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received UserID: %d", requestBody.UserId)
	log.Printf("Received Total: %d", requestBody.Total)
	log.Printf("Received Status: %d", requestBody.Status)
	log.Printf("Received StockID: %d", requestBody.StockID)
	log.Printf("Received Jumlah: %d", requestBody.Jumlah)

	if requestBody.UserId <= 0 || requestBody.Total <= 0 || requestBody.Status <= 0 || len(requestBody.StockID) == 0 || len(requestBody.Jumlah) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "UserId, Total, Status, Stock ID and Jumlah are required"})
	}

	addResult, err := model.AddTransaction(requestBody.UserId, requestBody.Total, requestBody.Status, requestBody.StockID, requestBody.Jumlah)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Transaction added successfully"})
}

func GetTransactionById(c echo.Context) error {
	transactionIDParam := c.Param("id")
	// c.QueryParam("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid stock ID"})
	}

	dtTransaction, err := model.GetTransactionById(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve stock"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Stock not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func GetTransactionDetailedById(c echo.Context) error {
	transactionIDParam := c.Param("id")
	// c.QueryParam("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction ID"})
	}

	dtTransaction, err := model.GetTransactionDetailedById(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transaction"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Transaction not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func GetAllTransaction(c echo.Context) error {
	dtTransaction, err := model.GetAllTransaction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transaction"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Transaction not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func UpdateTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction ID"})
	}

	var requestBody struct {
		Total   int   `json:"total"`
		Status  int   `json:"status"`
		StockID []int `json:"stock_id"`
		Jumlah  []int `json:"jumlah"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditTransaction(id, requestBody.Total, requestBody.Status, requestBody.StockID, requestBody.Jumlah)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Transaction updated successfully"})
}

func DeleteTransactionById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction ID"})
	}

	deleteResult, err := model.DeleteTransactionById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Transaction deleted successfully"})
}

func CreateTransactionStatus(c echo.Context) error {
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

	addResult, err := model.AddTransactionStatus(requestBody.Status, requestBody.Detail)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Transaction status added successfully"})
}

func GetTransactionStatusById(c echo.Context) error {
	transactionIDParam := c.Param("id")
	// c.QueryParam("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction status ID"})
	}

	dtTransaction, err := model.GetTransactionStatusById(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transaction"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Transaction status not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func GetAllTransactionStatus(c echo.Context) error {
	dtTransaction, err := model.GetAllTransactionStatus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve Transaction status"})
	}
	if dtTransaction.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Transaction status not found"})
	}

	return c.JSON(http.StatusOK, dtTransaction)
}

func UpdateTransactionStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction status ID"})
	}

	var requestBody struct {
		Status string `json:"status"`
		Detail string `json:"detail"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditTransactionStatus(id, requestBody.Status, requestBody.Detail)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Transaction status updated successfully"})
}

func DeleteTransactionStatusById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction status ID"})
	}

	deleteResult, err := model.DeleteTransactionStatusById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Transaction status deleted successfully"})
}
