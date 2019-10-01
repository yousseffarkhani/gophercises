package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const tableName = "phone_numbers"

// Phone represents the phone_numbers table in the DB
type Phone struct {
	ID     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var phone string
	err := db.QueryRow(fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, tableName), id).Scan(&id, &phone)
	if err != nil {
		return "", err
	}
	return phone, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := fmt.Sprintf(`INSERT INTO %s(value) VALUES($1) RETURNING id`, tableName)
	var id int
	err := db.QueryRow(statement, phone).Scan(&id) // Cette manière d'insérer avec $1 permet d'éviter les injections SQL
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) GetAllPhones() ([]Phone, error) {
	rows, err := db.db.Query(fmt.Sprintf(`SELECT id, value FROM %s`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []Phone
	for rows.Next() { // Retourne les résultats row par row. S'il y a une erreur va renvoyer false et quitter la boucle. Il faut donc vérifier que nous n'avons pas d'erreurs après la boucle.
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	err := db.db.QueryRow(fmt.Sprintf(`SELECT * FROM %s WHERE value=$1`, tableName), number).Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (db *DB) UpdatePhone(p *Phone) error {
	_, err := db.db.Exec(fmt.Sprintf(`UPDATE %s SET VALUE=$2 WHERE id=$1`, tableName), p.ID, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	_, err := db.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, tableName), id)
	return err
}

func Migrate(driverName, dataSource, tableName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db, tableName)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB, tableName string) error {
	statement := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id SERIAL,
	value VARCHAR(255)
)`, tableName)
	_, err := db.Exec(statement)
	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}
