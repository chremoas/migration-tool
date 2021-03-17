package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/chremoas/auth-srv/model"
	"github.com/doug-martin/goqu/v9"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
)

const (
	tagName     = "db"
	schemaQuery = `select kcu.table_name,
       kcu.column_name as key_column
from information_schema.table_constraints tco
         join information_schema.key_column_usage kcu
              on kcu.constraint_name = tco.constraint_name
                  and kcu.constraint_schema = tco.constraint_schema
                  and kcu.constraint_name = tco.constraint_name
where tco.constraint_type = 'PRIMARY KEY'
order by kcu.table_schema,
         kcu.table_name`
)

type tables struct {
	name  string
	model interface{}
}

func main() {
	var primaryKey = make(map[string]string)
	var tableName string
	var keyColumn string

	tableList := []tables{
		{name: "alliances", model: model.Alliance{}},
		{name: "authentication_codes", model: model.AuthenticationCode{}},
		{name: "characters", model: model.Character{}},
		{name: "corporations", model: model.Corporation{}},
		{name: "roles", model: model.Role{}},
		{name: "users", model: model.User{}},
	}

	// Read config file
	// Open source SQL DB
	sourceConnectionString := "postgres://root:@10.42.2.27:26257/chremoas-aba?sslmode=disable"
	sourceDB := sqlx.MustConnect("postgres", sourceConnectionString)
	defer sourceDB.Close()

	rows, err := sourceDB.Query(schemaQuery)

	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&tableName, &keyColumn)
		primaryKey[tableName] = keyColumn
	}

	// Open destination SQL DB
	destConnectionString := "postgres://postgres:@localhost:5432/chremoas_aba?sslmode=disable"
	destDB := sqlx.MustConnect("postgres", destConnectionString)
	defer destDB.Close()

	// Open source Redis DB
	// Open destination Redis DBo

	for _, table := range tableList {
		//tempData := table.model
		selectSql, _, _ := goqu.From(table.name).ToSQL()
		rows, _ := sourceDB.Queryx(selectSql)
		for rows.Next() {
			rows.StructScan(table.model)
			fmt.Println(table.model)
		}
		//	for _, alliance := range tempData {
		//		insert := goqu.Insert("alliances").Rows(
		//			alliance,
		//		)
		//		insertSQL, _, _ := insert.ToSQL()
		//		_, err := destDB.Exec(insertSQL)
		//		if err != nil {
		//			columnName, columnValue, err := getPrimaryKey(primaryKey["alliances"], alliance)
		//			if err != nil {
		//				fmt.Println(err.Error())
		//			}
		//			update := goqu.Update("alliances").Set(
		//				alliance,
		//			).Where(goqu.Ex{
		//				columnName: goqu.Op{"eq": columnValue},
		//			})
		//			updateSQL, _, _ := update.ToSQL()
		//			_, err = destDB.Exec(updateSQL)
		//			if err != nil {
		//				fmt.Println(err.Error())
		//			}
		//		}
		//	}
	}
}

func getPrimaryKey(keyName string, dataStruct interface{}) (string, string, error) {
	val := reflect.ValueOf(dataStruct) // could be any underlying type

	// if its a pointer, resolve its value
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// should double check we now have a struct (could still be anything)
	if val.Kind() != reflect.Struct {
		return "", "", errors.New("unexpected type")
	}

	// now we grab our values as before (note: I assume table name should come from the struct type)
	structType := val.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag

		columnName := tag.Get(tagName)

		if columnName == keyName {
			columnValue := val.Field(i).Interface()
			return columnName, fmt.Sprintf("%v", columnValue), nil
		}
	}

	return "", "", errors.New("fell out the bottom")
}
