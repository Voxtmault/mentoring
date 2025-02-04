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

func CreateSupplier(c echo.Context) error {
	var requestBody struct {
		Nama   string `json:"nama"`
		Kontak string `json:"kontak"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}
	log.Printf("Received Nama: %s", requestBody.Nama)
	log.Printf("Received Kontak: %d", requestBody.Kontak)

	if requestBody.Nama == "" || requestBody.Kontak == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Nama and kontak are required"})
	}

	addResult, err := model.AddSupplier(requestBody.Nama, requestBody.Kontak)
	if err != nil {
		return c.JSON(addResult.Status, map[string]string{"message": addResult.Message})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Supplier added successfully"})
}

func GetSupplierById(c echo.Context) error {
	supplierIDParam := c.Param("id")
	// c.QueryParam("id")
	supplierID, err := strconv.Atoi(supplierIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid supplier ID"})
	}

	dtSupplier, err := model.GetSupplierById(supplierID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve supplier"})
	}
	if dtSupplier.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Supplier not found"})
	}

	return c.JSON(http.StatusOK, dtSupplier)
}

func GetAllSupplier(c echo.Context) error {
	dtSupplier, err := model.GetAllSupplier()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve supplier"})
	}
	if dtSupplier.Data == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Supplier not found"})
	}

	return c.JSON(http.StatusOK, dtSupplier)
}

func UpdateSupplier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid supplier ID"})
	}

	var requestBody struct {
		Nama   string `json:"nama"`
		Kontak string `json:"kontak"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updateResult, err := model.EditSupplier(id, requestBody.Nama, requestBody.Kontak)
	if err != nil {
		return c.JSON(updateResult.Status, map[string]string{"message": updateResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplier updated successfully"})
}

func DeleteSupplier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid supplier ID"})
	}

	deleteResult, err := model.DeleteSupplierById(id)
	if err != nil {
		return c.JSON(deleteResult.Status, map[string]string{"message": deleteResult.Message})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplier deleted successfully"})
}
