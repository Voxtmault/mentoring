package internal

import (
	"database/sql"
	"log"
	"net/http"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
)

func AddSupplier(nama, kontak string) (class.Response, error) {
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
	INSERT INTO suppliers (name, contact_info) 
	VALUES (?,?)
	`
	_, err = tx.Exec(query, nama, kontak)
	if err != nil {
		log.Printf("Failed to insert suppliers: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert suppliers", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "Suppliers added successfully", Data: nil}, nil
}

func EditSupplier(id int, nama, kontak string) (class.Response, error) {
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
	UPDATE suppliers
	SET name = ?, contact_info = ? 
	WHERE supplier_id = ?
	`
	_, err = tx.Exec(query, nama, kontak, id)
	if err != nil {
		log.Printf("Failed to update supplier: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update supplier", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Supplier updated successfully", Data: nil}, nil
}

func GetSupplierById(id int) (class.Response, error) {
	var obj class.Supplier

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM suppliers
	WHERE supplier_id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.Name, &obj.ContactInfo, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No supplier found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "Supplier not found", Data: nil}, nil
		}
		log.Printf("Failed to scan supplier: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve supplier", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Supplier retrieved successfully", Data: obj}, nil
}

func GetAllSupplier() (class.Response, error) {
	var roles []class.Supplier

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM supplier
	`
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve supplier", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var role class.Supplier
		err := rows.Scan(&role.ID, &role.Name, &role.ContactInfo, &role.Created_at, &role.Updated_at)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to process supplier", Data: nil}, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve supplier", Data: nil}, err
	}

	if len(roles) == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "No supplier found", Data: nil}, nil
	}

	return class.Response{Status: http.StatusOK, Message: "Supplier retrieved successfully", Data: roles}, nil
}

func DeleteSupplierById(id int) (class.Response, error) {
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

	query := `DELETE FROM supplier WHERE supplier_id = ?`

	result, err := tx.Exec(query, id)
	if err != nil {
		log.Printf("Failed to execute delete query: %v\n", err)
		tx.Rollback()
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete supplier", Data: nil}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows", Data: nil}, err
	}

	if rowsAffected == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "Supplier not found", Data: nil}, nil
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Supplier deleted successfully"}, nil
}
