package main

import (
	"strings"
	"os"
	"time"
	"fmt"
	"encoding/base64"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Database *sql.DB
	FoundTables []string = []string{}
)

///////////////////////////////////////////
// DB Helpers
///////////////////////////////////////////

func StringToB64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func B64ToString(str string) (string, error) {
	f, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(f), nil
}

func DBBoolean(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func EstablishDB() {
	var (
		err error
		dsn string = os.Getenv("DB_DSN")
	)

	if dsn == "" {
		panic("EstablishDB(): Environment variable DB_DSN required with suitable DSN string for connecting to the database! Exiting...")
		os.Exit(1)
	}

	Database, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	Database.SetConnMaxLifetime(time.Minute * 2)
	Database.SetMaxOpenConns(5)
	Database.SetMaxIdleConns(5)

	err = Database.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping database with error: %s", err.Error()))
	}
}

func FetchKeyDB(key string, apiKey string) (res *Task, err error){
	var (
		t Task
		textHolder string
		childrenHolder string
	)

	err = CheckTableExistsDB(apiKey)
	if err != nil {
		return &t, err
	}

	q := fmt.Sprintf("SELECT id, text, checked, children FROM %s WHERE id = '%s';", apiKey, key)
	Log("FetchKeyDB(): Query => %s", q)
	rows, err := Database.Query(q)
	if err != nil {
		Log("FetchKeyDB(): Failed to fetch key %s from database with error - %s", key, err.Error())
		return &t, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &textHolder, &t.Checked, &childrenHolder)
		if err != nil {
			return &t, err
		}
	}

	err = rows.Err()
	if err != nil {
		return &t, err
	}

	t.Children = strings.Split(",", childrenHolder)
	t.Text, err = B64ToString(textHolder)
	if err != nil {
		return &t, err
	}
	return &t, nil
}

func DeleteKeyDB(key string, apiKey string) (err error){
	err = CheckTableExistsDB(apiKey)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", apiKey, key)
	Log("DeleteKeyDB(): Query => %s", q)
	_, err = Database.Exec(q)
	return err
}

func AddKeyDB(t *Task, apiKey string) (err error){
	err = CheckTableExistsDB(apiKey)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("INSERT INTO %s (id, text, checked, children) VALUE ('%s', '%s', '%s', '%s');", apiKey, t.ID, StringToB64(t.Text), DBBoolean(t.Checked), strings.Join(t.Children, ","))
	Log("AddKeyDB(): Query => %s", q)
	_, err = Database.Exec(q)
	return err
}

func CheckTableExistsDB(tableName string) error {
	if SliceContainsString(FoundTables, tableName) {
		return nil
	}
	
	var (
		name string
	)

	checkTablesStmt, err := Database.Query("SHOW TABLES;")
	if err != nil {
		Log("CheckTableExistsDB(): Could not get tables in the database. Exiting with Error:" + err.Error())
		return err
	}
	defer checkTablesStmt.Close()

	for checkTablesStmt.Next() {
		err = checkTablesStmt.Scan(&name)
		if err != nil {
			Log("CheckTableExistsDB(): Error when looping through DB Tables. Exiting with Error:" + err.Error())
			return err
		}
		if name == tableName {
			FoundTables = append(FoundTables, name)
			return nil
		}
	}

	err = CreateTableDB(tableName)
	if err != nil {
		Log("CheckTableExistsDB(): Failed to create a new database " + tableName + " with error:" + err.Error())
		return err
	}
	return nil
}

func CreateTableDB(tableName string) error {
	q := fmt.Sprintf("CREATE TABLE %s (id VARCHAR(255), text TEXT, checked BOOLEAN, children MEDIUMTEXT);", tableName)
	Log("CreateTableDB(): Creating new table %s with query '%s'", tableName, q)
	_, err := Database.Exec(q)
	return err
}