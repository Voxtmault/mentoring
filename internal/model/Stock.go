package internal

import (
	"database/sql"
	"log"
	"net/http"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
)

func AddStock(nama string, jumlah int, harga float32) (class.Response, error) {
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
	INSERT INTO stock (fruit_name, quantity, price_per_unit) 
	VALUES (?,?,?)
	`
	_, err = tx.Exec(query, nama, jumlah, harga)
	if err != nil {
		log.Printf("Failed to insert stock: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert stock", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "Stock added successfully", Data: nil}, nil
}

func EditStock(id int, nama string, jumlah int, harga float32) (class.Response, error) {
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
	UPDATE stock
	SET fruit_name = ?, quantity = ?, price_per_unit = ? 
	WHERE stock_id = ?
	`
	_, err = tx.Exec(query, nama, jumlah, harga, id)
	if err != nil {
		log.Printf("Failed to update stock: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update stock", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Stock updated successfully", Data: nil}, nil
}

func GetStockById(id int) (class.Response, error) {
	var obj class.Stock

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM stock
	WHERE stock_id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.FruitName, &obj.Quantity, &obj.PricePerUnit, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No stock found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "Stock not found", Data: nil}, nil
		}
		log.Printf("Failed to scan stock: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve stock", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Stock retrieved successfully", Data: obj}, nil
}

func GetAllStock() (class.Response, error) {
	var stockList []class.Stock

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM stock
	`
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve stock", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock class.Stock
		err := rows.Scan(&stock.ID, &stock.FruitName, &stock.Quantity, &stock.PricePerUnit, &stock.Created_at, &stock.Updated_at)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to process stock", Data: nil}, err
		}
		stockList = append(stockList, stock)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve stock", Data: nil}, err
	}

	if len(stockList) == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "No stock found", Data: nil}, nil
	}

	return class.Response{Status: http.StatusOK, Message: "Stock retrieved successfully", Data: stockList}, nil
}

func DeleteStockById(id int) (class.Response, error) {
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

	query := `DELETE FROM stock WHERE stock_id = ?`

	result, err := tx.Exec(query, id)
	if err != nil {
		log.Printf("Failed to execute delete query: %v\n", err)
		tx.Rollback()
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete stock", Data: nil}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows", Data: nil}, err
	}

	if rowsAffected == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "Stock not found", Data: nil}, nil
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Stock deleted successfully", Data: nil}, nil
}
