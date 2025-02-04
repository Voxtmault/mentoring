package internal

import (
	"database/sql"
	"log"
	"net/http"
	"tutor_rest/db"
	class "tutor_rest/internal/class"

	"github.com/google/uuid"
)

func AddTransaction(userid, total, status int, stockid, jumlah []int) (class.Response, error) {
	if len(stockid) != len(jumlah) {
		log.Printf("Stock ID and quantity length mismatch: %d vs %d", len(stockid), len(jumlah))
		return class.Response{Status: http.StatusBadRequest, Message: "Stock ID and quantity length mismatch", Data: nil}, nil
	}

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

	newUUID := uuid.New()

	query := `
	INSERT INTO transactions (transaction_uuid, transaction_status, user_id, transaction_date, total_amount) 
	VALUES (?,?,?,NOW(),?)
	`
	res, err := tx.Exec(query, newUUID, status, userid, total)
	if err != nil {
		log.Printf("Failed to insert user role: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert user role", Data: nil}, err
	}

	transactionId, err := res.LastInsertId()
	if err != nil {
		log.Printf("Failed to retrieve last insert ID: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve last insert ID", Data: nil}, err
	}

	// ambil harga tiap stok
	stockPrice := make([]float32, len(stockid))
	for i, sid := range stockid {
		queryStockPrice := `
		SELECT price_per_unit 
		FROM stock 
		WHERE stock_id = ?
		`
		err = tx.QueryRow(queryStockPrice, sid).Scan(&stockPrice[i])
		if err != nil {
			log.Printf("Failed to retrieve stock price for stock_id %d: %v\n", sid, err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve stock price", Data: nil}, err
		}
	}

	// masukkan ke transaction detail
	queryDetail := `
	INSERT INTO transactiondetails (transaction_id, stock_id, quantity, price_at_time_of_sale) 
	VALUES (?,?,?,?)
	`
	for i := 0; i < len(stockid); i++ {
		_, err = tx.Exec(queryDetail, transactionId, stockid[i], jumlah[i], stockPrice[i])
		if err != nil {
			log.Printf("Failed to insert transaction detail for stock_id %d: %v\n", stockid[i], err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert transaction detail", Data: nil}, err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "User Role added successfully", Data: map[string]interface{}{
		"transaction_id": transactionId,
		"uuid":           newUUID,
	}}, nil
}

func EditTransaction(id, total, status int, stockid, jumlah []int) (class.Response, error) {
	if len(stockid) != len(jumlah) {
		log.Printf("Stock ID and quantity length mismatch: %d vs %d", len(stockid), len(jumlah))
		return class.Response{Status: http.StatusBadRequest, Message: "Stock ID and quantity length mismatch", Data: nil}, nil
	}

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
	UPDATE transactions
	SET transaction_status = ?, total_amount = ? 
	WHERE id = ?
	`
	_, err = tx.Exec(query, status, total)
	if err != nil {
		log.Printf("Failed to insert user role: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update transaction", Data: nil}, err
	}

	for i := range stockid {
		queryCheckStock := `
		SELECT id 
		FROM transaction_details
		WHERE transaction_id = ? AND stock_id = ?
		`
		var detailID int
		err = tx.QueryRow(queryCheckStock, id, stockid[i]).Scan(&detailID)

		if err == sql.ErrNoRows {
			queryInsertDetail := `
			INSERT INTO transaction_details (transaction_id, stock_id, quantity, price_at_time_of_sale)
			VALUES (?, ?, ?, (SELECT price_per_unit FROM stock WHERE stock_id = ?))
			`
			_, err = tx.Exec(queryInsertDetail, id, stockid[i], jumlah[i], stockid[i])
			if err != nil {
				log.Printf("Failed to insert transaction detail for stock_id %d: %v\n", stockid[i], err)
				return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert transaction detail", Data: nil}, err
			}
		} else if err != nil {
			log.Printf("Failed to check existing transaction detail: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check transaction detail", Data: nil}, err
		} else {
			queryUpdateDetail := `
			UPDATE transaction_details
			SET quantity = ?
			WHERE id = ?
			`
			_, err = tx.Exec(queryUpdateDetail, jumlah[i], detailID)
			if err != nil {
				log.Printf("Failed to update transaction detail for stock_id %d: %v\n", stockid[i], err)
				return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update transaction detail", Data: nil}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User Role updated successfully", Data: nil}, nil
}

func GetTransactionById(id int) (class.Response, error) {
	var obj class.Transaction

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM transactions
	WHERE id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.UUID, &obj.Status, &obj.UserID, &obj.TransactionDate, &obj.TotalAmount, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No transaction found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "Transaction not found", Data: nil}, nil
		}
		log.Printf("Failed to scan transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Transaction retrieved successfully", Data: obj}, nil
}

func GetTransactionByUserId(id int) (class.Response, error) {
	var obj class.Transaction

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM transactions
	WHERE user_id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.UUID, &obj.Status, &obj.UserID, &obj.TransactionDate, &obj.TotalAmount, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No transaction found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "Transaction not found", Data: nil}, nil
		}
		log.Printf("Failed to scan transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Transaction retrieved successfully", Data: obj}, nil
}

func GetAllTransaction() (class.Response, error) {
	var allTransaction []class.Transaction

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM transaction
	`
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction class.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UUID, &transaction.Status, &transaction.UserID, &transaction.TransactionDate,
			&transaction.TotalAmount, &transaction.Created_at, &transaction.Updated_at)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to process transaction", Data: nil}, err
		}
		allTransaction = append(allTransaction, transaction)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction", Data: nil}, err
	}

	if len(allTransaction) == 0 {
		return class.Response{Status: http.StatusNotFound, Message: "No transaction found", Data: nil}, nil
	}

	return class.Response{Status: http.StatusOK, Message: "Transaction retrieved successfully", Data: allTransaction}, nil
}

func GetTransactionDetailedById(id int) (class.Response, error) {
	var obj class.TransactionDetailed
	var objDetail []class.TransactionDetail

	// database
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := `
	SELECT *
	FROM transactions
	WHERE id = ?
	`
	row := con.QueryRow(query, id)
	err = row.Scan(&obj.ID, &obj.UUID, &obj.Status, &obj.UserID, &obj.TransactionDate, &obj.TotalAmount, &obj.Created_at, &obj.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No transaction found with ID %d\n", id)
			return class.Response{Status: http.StatusNotFound, Message: "Transaction not found", Data: nil}, nil
		}
		log.Printf("Failed to scan transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction", Data: nil}, err
	}

	queryDetails := `
	SELECT id, transaction_id, stock_id, quantity, price_at_time_of_sale
	FROM transaction_details
	WHERE transaction_id = ?
	`
	rows, err := con.Query(queryDetails, id)
	if err != nil {
		log.Printf("Failed to query transaction details: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction details", Data: nil}, err
	}
	defer rows.Close()

	// Parse each row of transaction details
	for rows.Next() {
		var detail class.TransactionDetail
		err := rows.Scan(
			&detail.ID,
			&detail.TransactionID,
			&detail.StockID,
			&detail.Quantity,
			&detail.PriceAtTimeOfSale,
		)
		if err != nil {
			log.Printf("Failed to scan transaction detail: %v\n", err)
			return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve transaction detail", Data: nil}, err
		}
		objDetail = append(objDetail, detail)
	}
	obj.Detail = objDetail

	return class.Response{Status: http.StatusOK, Message: "Transaction retrieved successfully", Data: obj}, nil
}

func DeleteTransactionById(id int) (class.Response, error) {
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

	// transaction detail
	query := `DELETE FROM transaction_details WHERE transaction_id = ?`

	result, err := tx.Exec(query, id)
	if err != nil {
		log.Printf("Failed to execute delete query: %v\n", err)
		tx.Rollback()
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete transaction detail", Data: nil}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows in transaction_details: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows in transaction details", Data: nil}, err
	}
	if rowsAffected == 0 {
		log.Printf("No transaction details found for transaction_id %d\n", id)
		return class.Response{Status: http.StatusNotFound, Message: "Transaction details not found", Data: nil}, nil
	}

	// transaction
	queryTransaction := `DELETE FROM transactions WHERE id = ?`
	result, err = tx.Exec(queryTransaction, id)
	if err != nil {
		log.Printf("Failed to delete from transactions: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete transaction", Data: nil}, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check affected rows in transactions: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to check affected rows in transactions", Data: nil}, err
	}
	if rowsAffected == 0 {
		log.Printf("No transaction found for id %d\n", id)
		return class.Response{Status: http.StatusNotFound, Message: "Transaction not found", Data: nil}, nil
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User role deleted successfully"}, nil
}
