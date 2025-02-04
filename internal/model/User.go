package internal

import (
	// "fmt"
	"database/sql"
	"net/http"
	"os"
	"tutor_rest/db"
	class "tutor_rest/internal/class"
	"tutor_rest/pkg"

	// "tutor_rest/internal/class"
	"log"
)

func AddUser(username, password string, role, status int) (class.Response, error) {
	var obj class.User

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

	query := "SELECT user_id FROM users WHERE username = ?"
	err = tx.QueryRow(query, username).Scan(&obj.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying user: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Database query error", Data: nil}, err
	}
	if obj.ID != 0 {
		return class.Response{Status: http.StatusConflict, Message: "Username already exists", Data: nil}, nil
	}

	hashedPassword, err := pkg.HashBcrypt(password)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to hash password", Data: nil}, err
	}

	insertQuery := "INSERT INTO users (username, password, role, status) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(insertQuery, username, hashedPassword, role, status)
	if err != nil {
		log.Printf("Failed to insert user: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to insert user", Data: nil}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction", Data: nil}, err
	}

	return class.Response{Status: http.StatusCreated, Message: "User added successfully", Data: nil}, nil
}

func LoginUser(username, password string) (class.Response, error) {
	var obj class.User

	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := "SELECT user_id, username, role, password FROM users WHERE username = ?"

	err = con.QueryRow(query, username).
		Scan(&obj.ID, &obj.Username, &obj.Role, &obj.Password)

	if err != nil {
		log.Printf("Error querying user: %v\n", err)
		return class.Response{Status: http.StatusUnauthorized, Message: "Invalid Username", Data: nil}, err
	}

	err = pkg.CompareHashBcrypt(obj.Password, password)
	if err != nil {
		log.Printf("Invalid password: %v\n", err)
		return class.Response{Status: http.StatusUnauthorized, Message: "Invalid password", Data: nil}, err
	}

	data := map[string]interface{}{
		"token":   os.Getenv("API_KEY"),
		"user_id": obj.ID,
	}
	log.Printf("Login Successful: %v\n", data)

	return class.Response{Status: http.StatusOK, Message: "Login Successful", Data: data}, nil

}

func GetAllUsers() (class.Response, error) {
	users := []class.User{}

	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}, err
	}
	defer db.DbClose(con)

	query := "SELECT user_id, username, role, password FROM users"
	rows, err := con.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve users", Data: nil}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user class.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.Password)
		if err != nil {
			log.Printf("Failed to scan user: %v\n", err)
			continue // Optionally handle the error differently
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v\n", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Error processing user data", Data: nil}, err
	}

	return class.Response{Status: http.StatusOK, Message: "Retrieved all users successfully", Data: users}, nil
}

func GetUserByID(userID int) (*class.User, error) {
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return nil, err
	}
	defer db.DbClose(con)

	var user class.User
	query := "SELECT id, username, role FROM users WHERE id = ?"
	err = con.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil user if no user is found, handle this case in the controller
		}
		log.Printf("Error querying user by ID: %v\n", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user's details in the database.
func UpdateUser(id int, username string, role, status int) (class.Response, error) {
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error()}, err
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

	// Prepare the update statement
	stmt, err := tx.Prepare("UPDATE users SET username= ?, role= ?, status = ? WHERE id=?")
	if err != nil {
		log.Printf("Failed to prepare update statement: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to prepare update statement"}, err
	}
	defer stmt.Close()

	// Execute the update statement
	_, err = stmt.Exec(username, role, status, id)
	if err != nil {
		log.Printf("Failed to execute update: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to update user"}, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to commit transaction"}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User updated successfully"}, nil
}

// DeleteUser removes a user from the database.
func DeleteUser(id int) (class.Response, error) {
	con, err := db.DbConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: err.Error()}, err
	}
	defer db.DbClose(con)

	_, err = con.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return class.Response{Status: http.StatusInternalServerError, Message: "Failed to delete user"}, err
	}

	return class.Response{Status: http.StatusOK, Message: "User deleted successfully"}, nil
}
