package models

import (
	"database/sql"
	"strconv"
	
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Person struct {
	Id			int		`json:"id"`
	FirstName 		string		`json:"first_name"`
	LastName		string		`json:"last_name"`
	Email 			string		`json:"email"`
	IpAddress			string		`json:"ip_address"`
}

func DeletePerson(personId int) (bool, error) {
	tx, err := DB.Begin()

	if err != nil {
		return false, nil
	}

	stmt, err := DB.Prepare("DELETE from people WHERE id = ?")

	if err != nil {
		return false, nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(personId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdatePerson(ourPerson Person, id int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE id = ?")
	
	if err != nil {
		return false, nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, ourPerson.IpAddress, id)

	if err != nil {
		return false, err
	}
	tx.Commit()

	return true, nil

}

func AddPerson(newPerson Person) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)

	if err != nil {
		return false, nil
	}

	tx.Commit()

	return true, nil
}

func GetPersonById(id string) (Person, error){

	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email, ip_address FROM people WHERE id = ?")
	if err != nil {
		return Person{}, err
	}
	person := Person{}

	sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email, &person.IpAddress)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}

	return person, nil
}

func GetPersons(count int) ([]Person, error) {
	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address FROM people LIMIT " + strconv.Itoa(count))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IpAddress)
		if err != nil {
			return nil, err
		}
		people = append(people, singlePerson)
	}
	err = rows.Err()
	if err != nil {
		return nil, err	
	}
	return people, err
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./names.db")
	if err != nil {
		return err	
	}
	DB = db
	return nil
}
