package user

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DSN string
)

// Create database connection using the MySQL driver
func createDBConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getLoginInfoFromEmail(email string) (*LoginInfo, error) {
	var userInfo LoginInfo
	db, err := createDBConn()
	if err != nil {
		return &userInfo, err
	}
	defer db.Close()
	err = db.QueryRow("SELECT id,email,password FROM user WHERE email=?", email).Scan(&userInfo.UserID, &userInfo.Email, &userInfo.HashedPassword)
	return &userInfo, err
}

func addUserEntry(userData User, id string) error {
	var err error
	db, err := createDBConn()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user(id,firstname,lastname,email, password) VALUES(?, ?, ?, ?, ?)",
		id, userData.FirstName, userData.LastName, userData.Email, userData.Password)
	if err != nil {
		return err
	}
	return err
}

func addToBlackList(token string) error {
	var err error
	db, err := createDBConn()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO blacklist(token) VALUES(?)", token)
	if err != nil {
		return err
	}
	return err
}

func getFromBlackList(token string) (bool, error) {
	var tokenCheck bool = true
	db, err := createDBConn()
	if err != nil {
		return tokenCheck, err
	}
	defer db.Close()
	var id int
	err = db.QueryRow("SELECT id FROM blacklist WHERE token=?", token).Scan(&id)
	if id == 0 && err == sql.ErrNoRows {
		err = nil
		tokenCheck = false
	}
	return tokenCheck, err
}

func getUserDataFromDB(id string) (*User, error) {
	var userInfo User
	db, err := createDBConn()
	if err != nil {
		return &userInfo, err
	}
	defer db.Close()
	err = db.QueryRow("SELECT email,firstname,lastname FROM user WHERE id=?", id).Scan(&userInfo.Email,
		&userInfo.FirstName, &userInfo.LastName)
	return &userInfo, err
}

func updateUserNameDB(id, firstName, lastName string) error {
	db, err := createDBConn()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("UPDATE user SET firstname =?,lastname=? WHERE id=?", id).Scan(firstName, lastName, id)
	return err
}

func deleteUserFromDB(id string) error {
	db, err := createDBConn()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("DELETE FROM user WHERE id=?", id).Scan(id)
	return err
}
