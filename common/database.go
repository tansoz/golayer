package common

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
)
import _ "github.com/go-sql-driver/mysql"

type Database interface {
	Connection()
	Query(string, map[string]interface{}) *sql.Rows
	FetchAll(*sql.Rows, interface{}) []interface{}
	FetchObject(*sql.Rows, interface{}) interface{}
}

type DatabaseImpl struct {
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	Conn             *sql.DB
	Charset          string
	DataBaseTemplate TemplateDataType
}

func NewDatabase(host string, port int, user string, password string, DBName string, charset string) Database {
	db := &DatabaseImpl{
		Host:             host,
		Port:             port,
		User:             user,
		Password:         password,
		DBName:           DBName,
		Charset:          charset,
		DataBaseTemplate: NewTemplateDataType(),
	}

	db.DataBaseTemplate.AddDataType("MYSQL_INT", "[+-]?\\d")
	db.DataBaseTemplate.AddDataType("MYSQL_UINT", "\\d")
	db.DataBaseTemplate.AddDataType("MYSQL_DATETIME", "\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}-\\d{1,2}-\\d{1,2}")
	db.DataBaseTemplate.AddDataType("MYSQL_VARCHAR", ".")
	db.DataBaseTemplate.AddDataType("MYSQL_TEXT", ".{0,3000}")

	return db
}

func (this *DatabaseImpl) Connection() {

	dbhost := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", this.User, this.Password, this.Host, this.Port, this.DBName, this.Charset)

	if conn, err := sql.Open("mysql", dbhost); err == nil {
		if err := conn.Ping(); err != nil {
			panic(err)
		} else {
			this.Conn = conn
		}
	} else {
		panic(err)
	}
}

func (this *DatabaseImpl) Query(tsql string, data map[string]interface{}) *sql.Rows {

	sqlTemplate := this.DataBaseTemplate.NewTemplate(tsql)

	sqlStr := sqlTemplate.ToString(data, "\\'\n\r")

	if sqlStr == "" {
		return nil
	}
	return this.query(sqlStr)
}
func (this *DatabaseImpl) query(sql string) *sql.Rows {

	if rows, err := this.Conn.Query(sql); err == nil {
		return rows
	}
	return nil
}
func (this *DatabaseImpl) FetchAll(rows *sql.Rows, object interface{}) []interface{} {

	objectList := []interface{}{}

	for rows != nil {
		tmpObject := this.FetchObject(rows, object)
		if tmpObject != nil {
			objectList = append(objectList, tmpObject)
		} else {
			break
		}
	}

	return objectList
}
func (this *DatabaseImpl) FetchObject(rows *sql.Rows, object interface{}) interface{} {

	if rows.Next() {

		objectValue := reflect.ValueOf(object)
		objectInstance := reflect.New(objectValue.Type())

		objectFields := make(map[string]int)

		for i := 0; i < objectValue.NumField(); i++ {
			name := strings.ToLower(objectValue.Type().Field(i).Tag.Get("json"))
			objectFields[name] = i + 1
		}

		fieldsType, _ := rows.ColumnTypes()
		tmp := make([]interface{}, len(fieldsType))
		tmpKey := make([]string, len(fieldsType))

		for k, i := range fieldsType {

			name := strings.ToLower(i.Name())
			if objectFields[name] > 0 && !strings.Contains(i.ScanType().Name(), "Null") && !strings.Contains(i.ScanType().Name(), "RawBytes") {
				structfield := objectInstance.Elem().Field(objectFields[name] - 1)
				tmp[k] = reflect.New(structfield.Type()).Interface()
			} else {
				tmp[k] = reflect.New(i.ScanType()).Interface()
			}
			tmpKey[k] = name
		}

		if err := rows.Scan(tmp...); err == nil {

			for k, i := range tmp {
				if objectFields[tmpKey[k]] > 0 {
					structfield := objectInstance.Elem().Field(objectFields[tmpKey[k]] - 1)
					structfield.Set(convertType(i, objectValue.Type().Field(objectFields[tmpKey[k]]-1)))
				}
			}
			return objectInstance.Elem().Interface()
		} else {
			fmt.Println(err)
		}
	}

	return nil
}

func convertType(i interface{}, target reflect.StructField) reflect.Value {
	switch t := i.(type) {

	case *mysql.NullTime:
		return convertTime(t, target)
	case *sql.NullString:
		return convertString(t, target)
	case *sql.RawBytes:
		return convertStringFromBytes(t, target)
	case *sql.NullInt64:
		return convertInt64(t, target)
	default:
		return reflect.ValueOf(t).Elem()

	}
}
func convertTime(t *mysql.NullTime, target reflect.StructField) reflect.Value {

	switch target.Type.Name() {

	default:
		fallthrough
	case "string":
		if t.Valid {
			return reflect.ValueOf(t.Time.Format("2006-01-02 15:04:05"))
		} else if target.Tag.Get("default") != "" {
			return reflect.ValueOf(target.Tag.Get("default"))
		} else {
			return reflect.ValueOf("0000-00-00 00:00:00")
		}
	case "int":
		if t.Valid {
			return reflect.ValueOf(int(t.Time.Unix()))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 64); err == nil {
				return reflect.ValueOf(int(num))
			}
		}
		return reflect.ValueOf(0)
	case "int8":
		if t.Valid {
			return reflect.ValueOf(int16(t.Time.Unix()))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 8); err == nil {
				return reflect.ValueOf(int8(num))
			}
		}
		return reflect.ValueOf(int8(0))
	case "int16":
		if t.Valid {
			return reflect.ValueOf(int16(t.Time.Unix()))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 16); err == nil {
				return reflect.ValueOf(int16(num))
			}
		}
		return reflect.ValueOf(int16(0))
	case "int32":
		if t.Valid {
			return reflect.ValueOf(int32(t.Time.Unix()))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 32); err == nil {
				return reflect.ValueOf(int32(num))
			}
		}
		return reflect.ValueOf(int32(0))
	case "int64":
		if t.Valid {
			return reflect.ValueOf(int64(t.Time.Unix()))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 64); err == nil {
				return reflect.ValueOf(int64(num))
			}
		}
		return reflect.ValueOf(int64(0))
	case "mysql.NullTime":
		return reflect.ValueOf(*t)

	}
}
func convertString(t *sql.NullString, target reflect.StructField) reflect.Value {

	switch target.Type.Name() {

	default:
		fallthrough
	case "string":
		if t.Valid {
			return reflect.ValueOf(t.String)
		} else if target.Tag.Get("default") != "" {
			return reflect.ValueOf(target.Tag.Get("default"))
		} else {
			return reflect.ValueOf("")
		}
	case "sql.NullString":
		return reflect.ValueOf(*t)
	}

}
func convertStringFromBytes(t *sql.RawBytes, target reflect.StructField) reflect.Value {

	switch target.Type.Name() {

	default:
		fallthrough
	case "string":
		if len(*t) > 0 {
			return reflect.ValueOf(string(*t))
		} else if target.Tag.Get("default") != "" {
			return reflect.ValueOf(target.Tag.Get("default"))
		} else {
			return reflect.ValueOf("")
		}
	case "sql.RawBytes":
		return reflect.ValueOf(*t)
	}
}
func convertInt64(t *sql.NullInt64, target reflect.StructField) reflect.Value {
	switch target.Type.Name() {

	default:
		fallthrough
	case "string":
		if t.Valid {
			return reflect.ValueOf(strconv.FormatInt(t.Int64, 10))
		} else if target.Tag.Get("default") != "" {
			return reflect.ValueOf(target.Tag.Get("default"))
		} else {
			return reflect.ValueOf("")
		}
	case "int":
		if t.Valid {
			return reflect.ValueOf(int(t.Int64))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseInt(target.Tag.Get("default"), 10, 64); err == nil {
				return reflect.ValueOf(int(num))
			}
			return reflect.ValueOf(int(0))
		} else {
			return reflect.ValueOf(int(0))
		}
	case "uint":
		if t.Valid {
			return reflect.ValueOf(uint(t.Int64))
		} else if target.Tag.Get("default") != "" {
			if num, err := strconv.ParseUint(target.Tag.Get("default"), 10, 64); err == nil {
				return reflect.ValueOf(uint(num))
			}
			return reflect.ValueOf(uint(0))
		} else {
			return reflect.ValueOf(uint(0))
		}
	}
}
