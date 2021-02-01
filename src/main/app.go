package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Message message type
type Message struct {
	ID         int64
	Content    string
	CreateDate time.Time
}

func main() {

	// TODO: read from config
	HOSTNAME := "localhost"
	PORT := "3306"
	USERNAME := "root"
	PASSWORD := "password_here"
	DATABASE := "acme"

	connection := USERNAME + ":" + PASSWORD + "@tcp(" + HOSTNAME + ":" + PORT + ")/" + DATABASE + "?parseTime=true"
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	id, err := create(db, "Inserted row")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted row id %d\n", id)

	msg, err := readOne(db, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read row id %d\n", id)
	fmt.Println(msg)

	err = update(db, id, "Updated row")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated row id %d\n", id)

	msgs, err := readAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read all rows\n")
	for _, msg := range msgs {
		fmt.Printf("  %+v\n", msg)
	}

	err = delete(db, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted row id %d\n", id)

}

func create(db *sql.DB, content string) (int64, error) {
	var (
		id int64
	)
	stmt, err := db.Prepare("CALL messages_create(?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(content).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func readOne(db *sql.DB, _id int64) (*Message, error) {
	var (
		id         int64
		content    string
		createdate time.Time
	)
	stmt, err := db.Prepare("CALL messages_read_by_id(?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(_id).Scan(&id, &content, &createdate)
	if err != nil {
		return nil, err
	}
	msg := Message{ID: id, Content: content, CreateDate: createdate}
	return &msg, nil
}

func readAll(db *sql.DB) ([]Message, error) {
	var (
		id         int64
		content    string
		createdate time.Time
	)
	stmt, err := db.Prepare("CALL messages_read_all()")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgs := make([]Message, 0, 100) // TODO: get length
	for rows.Next() {
		if err := rows.Scan(&id, &content, &createdate); err != nil {
			return nil, err
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, Message{ID: id, Content: content, CreateDate: createdate})
	}
	return msgs, nil
}

func update(db *sql.DB, id int64, content string) error {
	stmt, err := db.Prepare("CALL messages_update(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, content)
	return err
}

func delete(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("CALL messages_delete(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}
