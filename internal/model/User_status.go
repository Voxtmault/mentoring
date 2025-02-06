package internal

import (
	"database/sql"
	"log"
	"net/http"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
)

func AddUserStatus(status, detail string) (class.Response, error) {
	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	// transaction
	tx, err := con.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to begin transaction", Data: nil}, err
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	query := `
	INSERT INTO user_status (status, detail) 
	VALUES (?,?)
	`
	_, err = tx.Exec(query, status, detail)
	if err != nil {
		log.Printf("Failed to insert user status: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert user status", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "User status added successfully", Data: nil}, nil
}

func EditUserStatus(id int, status, detail string) (class.Response, error) {
	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	// transaction
	tx, err := con.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to begin transaction", Data: nil}, err
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	query := `
	UPDATE user_status
	SET status = ?, detail = ? 
	WHERE id = ?
	`
	_, err = tx.Exec(query, status, detail, id)
	if err != nil {
		log.Printf("Failed to insert user status: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update user status", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User status updated successfully", Data: nil}, nil
}

func GetUserStatusById(id int) (class.Response, error) {
	var obj class.UserStatus

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM user_status
	WHERE id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.Status, &obj.Detail, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user status found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "User status not found", Data: nil}, nil
		}
		log.Printf("Failed to scan user status: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user status", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User status retrieved successfully", Data: obj}, nil
}

func GetAllUserStatus() (class.Response, error) {
	var statusList []class.UserStatus

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM user_status
	`
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user status", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var status class.UserStatus
		err := rows.Scan(&status.ID, &status.Status, &status.Detail, &status.Created_at, &status.Updated_at)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to process user status", Data: nil}, err
		}
		statusList = append(statusList, status)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user status", Data: nil}, err
	}

	if len(statusList) == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "No user status found", Data: nil}, nil
	}

	return class.Response{Status: http.StatusOK, Message: "User status retrieved successfully", Data: statusList}, nil
}

func DeleteUserStatusById(id int) (class.Response, error) {
	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	tx, err := con.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to begin transaction", Data: nil}, err
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			log.Printf("Rolling back transaction due to an error or panic: %v\n", r)
			tx.Rollback()
		}
	}()

	query := `DELETE FROM user_status WHERE id = ?`

	result, err := tx.Exec(query, id)
	if err != nil {
		log.Printf("Failed to execute delete query: %v\n", err)
		tx.Rollback()
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete user status", Data: nil}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows", Data: nil}, err
	}

	if rowsAffected == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "User status not found", Data: nil}, nil
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User status deleted successfully"}, nil
}
