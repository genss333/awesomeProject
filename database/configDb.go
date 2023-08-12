package database

import (
	"awesomeProject/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

var username = utils.GoDotEnvVariable("DB_USERNAME")
var host = utils.GoDotEnvVariable("DB_HOST")
var port = utils.GoDotEnvVariable("DB_PORT")
var database = utils.GoDotEnvVariable("DB_DATABASE")

func Connect() (*sql.DB, error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s:%s)/%s", username, host, port, database))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CustomExec(query string) error {
	db, err := Connect()
	if err != nil {
		return err
	}

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CustomQuery(query string) (*sql.Rows, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
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

func Count(table string, condition map[string]interface{}) (int, error) {
	db, err := Connect()
	if err != nil {
		return 0, err
	}

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
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, whereClause)

	var count int
	err = db.QueryRow(query, whereParams...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func Sum(table string, column string, condition map[string]interface{}) (int, error) {
	db, err := Connect()
	if err != nil {
		return 0, err
	}

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
	query := fmt.Sprintf("SELECT SUM(%s) FROM %s %s", column, table, whereClause)

	var sum int
	err = db.QueryRow(query, whereParams...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func Max(table string, column string, condition map[string]interface{}) (int, error) {
	db, err := Connect()
	if err != nil {
		return 0, err
	}

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
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s %s", column, table, whereClause)

	var max int
	err = db.QueryRow(query, whereParams...).Scan(&max)
	if err != nil {
		return 0, err
	}

	return max, nil
}

func Min(table string, column string, condition map[string]interface{}) (int, error) {
	db, err := Connect()
	if err != nil {
		return 0, err
	}

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
	query := fmt.Sprintf("SELECT MIN(%s) FROM %s %s", column, table, whereClause)

	var min int
	err = db.QueryRow(query, whereParams...).Scan(&min)
	if err != nil {
		return 0, err
	}

	return min, nil
}

func Avg(table string, column string, condition map[string]interface{}) (int, error) {
	db, err := Connect()
	if err != nil {
		return 0, err
	}

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
	query := fmt.Sprintf("SELECT AVG(%s) FROM %s %s", column, table, whereClause)

	var avg int
	err = db.QueryRow(query, whereParams...).Scan(&avg)
	if err != nil {
		return 0, err
	}

	return avg, nil
}

func CreateTable(table string, model interface{}) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY)", table))
	if err != nil {
		return err
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct")
	}

	numFields := modelType.NumField()
	for i := 0; i < numFields; i++ {
		field := modelType.Field(i)
		_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, field.Name, getColumnType(field.Type)))
		if err != nil {
			return err
		}
	}

	return nil
}

func getColumnType(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INT UNSIGNED"
	case reflect.Float32, reflect.Float64:
		return "FLOAT"
	case reflect.String:
		return "VARCHAR(255)"
	default:
		return "TEXT"
	}
}

func IsTable(table string) (bool, error) {
	db, err := Connect()
	if err != nil {
		return false, err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", table).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
