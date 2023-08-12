package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/rest_api_demo")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Select(table string, condition map[string]interface{}) (*sql.Rows, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(db)

	whereValues := make([]string, 0)
	whereParams := make([]interface{}, 0)

	for column, value := range condition {
		whereValues = append(whereValues, fmt.Sprintf("%s = ?", column))
		whereParams = append(whereParams, value)
	}

	whereClause := strings.Join(whereValues, " AND ")
	if condition != nil {
		whereClause = fmt.Sprintf("WHERE %s", whereClause)
	}
	query := fmt.Sprintf("SELECT * FROM %s %s", table, whereClause)

	rows, err := db.Query(query, whereParams...)

	return rows, err
}

func JoinTables(arrTable []string, joinConditions []string) (*sql.Rows, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	tables := strings.Join(arrTable, " JOIN ")
	joinString := strings.Join(joinConditions, " AND ")

	query := fmt.Sprintf("SELECT * FROM %s ON %s", tables, joinString)
	rows, err := db.Query(query)

	return rows, nil
}

func Insert(table string, data map[string]interface{}) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	columns := make([]string, 0)
	values := make([]string, 0)
	params := make([]interface{}, 0)

	for column, value := range data {
		columns = append(columns, column)
		values = append(values, "?")
		params = append(params, value)
	}

	columnsStr := fmt.Sprintf("`%s`", strings.Join(columns, "`, `"))
	valuesStr := strings.Join(values, ", ")

	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, columnsStr, valuesStr)
	_, err = db.Exec(sql, params...)
	if err != nil {
		return err
	}

	return nil
}

func Update(table string, data map[string]interface{}, condition map[string]interface{}) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	setValues := make([]string, 0)
	setParams := make([]interface{}, 0)

	for column, value := range data {
		setValues = append(setValues, fmt.Sprintf("%s = ?", column))
		setParams = append(setParams, value)
	}

	whereValues := make([]string, 0)
	whereParams := make([]interface{}, 0)

	for column, value := range condition {
		whereValues = append(whereValues, fmt.Sprintf("%s = ?", column))
		whereParams = append(whereParams, value)
	}

	setClause := strings.Join(setValues, ", ")
	whereClause := strings.Join(whereValues, " AND ")

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setClause, whereClause)
	params := append(setParams, whereParams...)

	_, err = db.Exec(sql, params...)
	if err != nil {
		return err
	}

	return nil
}

func Delete(table string, condition map[string]interface{}) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	whereValues := make([]string, 0)
	whereParams := make([]interface{}, 0)

	for column, value := range condition {
		whereValues = append(whereValues, fmt.Sprintf("%s = ?", column))
		whereParams = append(whereParams, value)
	}

	whereClause := strings.Join(whereValues, " AND ")

	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause)

	_, err = db.Exec(sql, whereParams...)
	if err != nil {
		return err
	}

	return nil
}
