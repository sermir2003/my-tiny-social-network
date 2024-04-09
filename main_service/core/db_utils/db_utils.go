package db_utils

import (
	"database/sql"
	"fmt"
	"main_service_core/models"
	salt_utils "main_service_core/salt_utils"
	"os"

	uuid "github.com/google/uuid"
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

func CreateNewUser(creds models.Credentials) (uuid.UUID, error) {
	tx, err := DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var cnt_such_login int
	err = tx.QueryRow(`
		SELECT count(*)
		FROM credentials
		WHERE login = $1`,
		creds.Login,
	).Scan(&cnt_such_login)
	if err != nil {
		return uuid.Nil, err
	}

	if cnt_such_login > 0 {
		return uuid.Nil, nil
	}

	id := uuid.New()
	salt, err := salt_utils.GenerateSalt()
	if err != nil {
		return uuid.Nil, err
	}
	password := salt_utils.HashPassword(creds.Password, salt)

	_, err = tx.Exec(`
		INSERT INTO credentials (id, login, salt, password)
		VALUES ($1, $2, $3, $4)`,
		id,
		creds.Login,
		salt,
		password,
	)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO personal_data (id, name, surname, birthdate, email, phone)
		VALUES ($1, NULL, NULL, NULL, NULL, NULL)`,
		id,
	)
	if err != nil {
		return uuid.Nil, err
	}

	tx.Commit()

	return id, nil
}

func CheckPassword(creds models.Credentials) (bool, error) {
	var salt []byte
	var existing_password string
	err := DB.QueryRow(`
		SELECT salt, password
		FROM credentials
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

func GetId(login string) (uuid.UUID, error) {
	var id uuid.UUID
	err := DB.QueryRow(`
		SELECT id
		FROM credentials
		WHERE login = $1`,
		login,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func UpdatePersonal(id uuid.UUID, data models.PersonalData) error {
	_, err := DB.Exec(`
		UPDATE personal_data
		SET (name, surname, birthdate, email, phone) = ($2, $3, $4, $5, $6)
		WHERE id = $1`,
		id,
		data.Name,
		data.Surname,
		data.Birthdate,
		data.Email,
		data.Phone,
	)
	return err
}
