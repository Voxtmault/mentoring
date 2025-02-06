package internal

import (
	"database/sql"
	"log"
	"net/http"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
)

func AddUserRole(role, detail string) (class.Response, error) {
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
	INSERT INTO user_role (role, detail) 
	VALUES (?,?)
	`
	_, err = tx.Exec(query, role, detail)
	if err != nil {
		log.Printf("Failed to insert user role: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert user role", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "User Role added successfully", Data: nil}, nil
}

func EditUserRole(id int, role, detail string) (class.Response, error) {
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
	UPDATE user_role
	SET role = ?, detail = ? 
	WHERE id = ?
	`
	_, err = tx.Exec(query, role, detail, id)
	if err != nil {
		log.Printf("Failed to insert user role: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update user role", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User Role updated successfully", Data: nil}, nil
}

func GetUserRoleById(id int) (class.Response, error) {
	var obj class.UserRole

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM user_role
	WHERE id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.Role, &obj.Detail, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user role found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "User role not found", Data: nil}, nil
		}
		log.Printf("Failed to scan user role: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user role", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User role retrieved successfully", Data: obj}, nil
}

func GetAllUserRole() (class.Response, error) {
	var roles []class.UserRole

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM user_role
	`
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user roles", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var role class.UserRole
		err := rows.Scan(&role.ID, &role.Role, &role.Detail, &role.Created_at, &role.Updated_at)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to process user roles", Data: nil}, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve user roles", Data: nil}, err
	}

	if len(roles) == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "No user roles found", Data: nil}, nil
	}

	return class.Response{Status: http.StatusOK, Message: "User role retrieved successfully", Data: roles}, nil
}

func DeleteUserRoleById(id int) (class.Response, error) {
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

	query := `DELETE FROM user_role WHERE id = ?`

	result, err := tx.Exec(query, id)
	if err != nil {
		log.Printf("Failed to execute delete query: %v\n", err)
		tx.Rollback()
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete user role", Data: nil}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows", Data: nil}, err
	}

	if rowsAffected == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "User role not found", Data: nil}, nil
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User role deleted successfully"}, nil
}
