package db_utils

import (
	"database/sql"
	"fmt"
	"log"
	"main_service_core/models"
	salt_utils "main_service_core/salt_utils"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func StartUpDB() error {
	postgresConnectionLine := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", postgresConnectionLine)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CreateNewUser(creds models.Credentials) (id uint32, err error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var cnt_such_login int
	err = tx.QueryRow(`
		SELECT count(*)
		FROM users
		WHERE login = $1`,
		creds.Login,
	).Scan(&cnt_such_login)
	if err != nil {
		return 0, err
	}

	if cnt_such_login > 0 {
		return 0, nil
	}

	salt, err := salt_utils.GenerateSalt()
	if err != nil {
		return 0, err
	}
	password := salt_utils.HashPassword(creds.Password, salt)

	err = tx.QueryRow(`
		INSERT INTO users (login, salt, password)
		VALUES ($1, $2, $3)
		RETURNING id`,
		creds.Login,
		salt,
		password,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CheckPassword(creds models.Credentials) (bool, error) {
	var salt []byte
	var existing_password string
	err := DB.QueryRow(`
		SELECT salt, password
		FROM users
		WHERE login = $1`,
		creds.Login,
	).Scan(&salt, &existing_password)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return salt_utils.DoPasswordsMatch(creds.Password, salt, existing_password), nil
}

func GetId(login string) (id uint32, err error) {
	err = DB.QueryRow(`
		SELECT id
		FROM users
		WHERE login = $1`,
		login,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdatePersonal(id uint32, personal models.PersonalData) error {
	log.Println("UpdatePersonal id=", id)
	_, err := DB.Exec(`
		UPDATE users
		SET (name, surname, birthdate, email, phone) = ($2, $3, $4, $5, $6)
		WHERE id = $1`,
		id,
		personal.Name,
		personal.Surname,
		personal.Birthdate,
		personal.Email,
		personal.Phone,
	)
	return err
}
