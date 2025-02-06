package controller

import (
	// "tutor_rest/internal/model"
	// "tutor_rest/internal/class"

	"mime/multipart"
	"net/http"
	model "tutor_rest/internal/model"

	"github.com/labstack/echo/v4"
)

func GetFile(c echo.Context) error {
	path := c.FormValue("path")
	return c.File(path)
}

func UploadFile(c echo.Context) error {
	id := c.FormValue("id")
	kolom_id := c.FormValue("kolom_id")
	folder := c.FormValue("folder")
	kolom := c.FormValue("kolom")
	file1, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UploadFile(file1, id, kolom_id, folder, kolom)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Berhasil upload file"})
}

func UploadMultipleFile(c echo.Context) error {
	id := c.FormValue("id")
	kolom_id := c.FormValue("kolom_id")
	folder := c.FormValue("folder")
	kolom := c.FormValue("kolom")
	var files []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to parse multipart form: " + err.Error()})
	}

	uploadedFiles := form.File["file"]
	if len(uploadedFiles) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No files uploaded"})
	}

	files = append(files, uploadedFiles...)

	result, err := model.UploadMultipleFile(files, id, kolom_id, folder, kolom)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Multiple file added successfully"})
}
