package internal

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
)

func UploadFile(file1 *multipart.FileHeader, id, kolom_id, folder, kolom string) (class.Response, error) {
	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	log.Println("Upload File")
	nId, _ := strconv.Atoi(id)
	// file.Filename =
	pathFile := "uploads/user/" + file1.Filename
	//source
	src, err := file1.Open()
	if err != nil {
		log.Println(err.Error())
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("uploads/" + folder + "/" + file1.Filename)
	if err != nil {
		log.Println("2")
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err.Error())
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	dst.Close()

	// update link ke db
	// tabel, kolom, link file, nama kolom id, id mu
	err = UpdateDataPath(folder, kolom, pathFile, kolom_id, nId)
	if err != nil {
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "File added successfully", Data: nil}, nil
}

func UploadMultipleFile(file1 []*multipart.FileHeader, id, kolom_id, folder, kolom string) (class.Response, error) {
	log.Println("Upload Multiple File")
	log.Println(file1)
	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	filePaths := []string{}
	for _, fileHeader := range file1 {
		filePath, err := saveAssetGambar(folder, id, fileHeader)
		if err != nil {
			log.Printf("Failed to store file: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
		}
		filePaths = append(filePaths, filePath)
	}

	nId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Failed to convert id from string to int: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}

	linkFilePath := strings.Join(filePaths, ",")

	err = UpdateDataPath(folder, kolom, linkFilePath, kolom_id, nId)
	if err != nil {
		log.Printf("Failed to update file link on database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "Multiple File added successfully", Data: nil}, nil
}

func UpdateDataPath(tabel string, kolom string, path string, kolom_id string, id int) error {
	// Open DB connection
	con, err := db.DbConnection()
	if err != nil {
		log.Println("error: " + err.Error())
		return err
	}
	defer db.DbClose(con)

	query := fmt.Sprintf("UPDATE %s SET %s='%s' WHERE %s = %d", tabel, kolom, path, kolom_id, id)

	_, err = con.Exec(query)
	if err != nil {
		log.Println("error executing query: " + err.Error())
		return err
	}

	return nil
}

func saveAssetGambar(folder, id string, fileHeader *multipart.FileHeader) (string, error) {

	//source
	src, err := fileHeader.Open()
	if err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("error: " + err.Error()), err
	}
	defer src.Close()

	// Destination
	filePath := fmt.Sprintf("uploads/" + folder + "/" + id + "_" + fileHeader.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("2")
		return fmt.Sprintf("error: " + err.Error()), err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("error: " + err.Error()), err
	}
	dst.Close()

	return filePath, nil
}
