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

func CreateStockItem(c echo.Context) error {
	var requestBody struct {
		Nama   string  `json:"nama"`
		Jumlah int     `json:"jumlah"`
		Harga  float32 `json:"harga"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Nama: %s", requestBody.Nama)
	log.Printf("Received Jumlah: %d", requestBody.Jumlah)
	log.Printf("Received Harga: %f", requestBody.Harga)

	if requestBody.Nama == "" || requestBody.Jumlah <= 0 || requestBody.Harga <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Nama, jumlah and harga are required"})
	}

	addResult, err := model.AddStock(requestBody.Nama, requestBody.Jumlah, requestBody.Harga)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Stock added successfully"})
}

func GetStockItem(c echo.Context) error {
	stockIDParam := c.Param("id")
	// c.QueryParam("id")
	stockID, err := strconv.Atoi(stockIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid stock ID"})
	}

	dtStock, err := model.GetStockById(stockID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve stock"})
	}
	if dtStock.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Stock not found"})
	}

	return c.JSON(http.StatusOK, dtStock)
}

func GetAllStockItem(c echo.Context) error {
	dtStock, err := model.GetAllStock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve stock"})
	}
	if dtStock.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Stock not found"})
	}

	return c.JSON(http.StatusOK, dtStock)
}

func UpdateStockItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid stock ID"})
	}

	var requestBody struct {
		Nama   string  `json:"nama"`
		Jumlah int     `json:"jumlah"`
		Harga  float32 `json:"harga"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditStock(id, requestBody.Nama, requestBody.Jumlah, requestBody.Harga)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Stock updated successfully"})
}

func DeleteStockItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid stock ID"})
	}

	deleteResult, err := model.DeleteStockById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Stock deleted successfully"})
}
