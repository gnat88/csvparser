package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
)

//Parser parse a csv file and returns an array of pointers of the type specified
type Parser interface {
	Parse() (interface{}, error)
}

//CsvParser parses a csv file and returns an array of pointers the type specified
type CsvParser struct {
	CsvFile      string
	CsvSeparator rune
	BindObject   interface{}
	Setter       func(field reflect.Value, colName string, raw string) bool
}

//Parse creates the array of the given type from the csv file
func (parser CsvParser) Parse() (interface{}, error) {

	csvFile, err := os.Open(parser.CsvFile)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	var csvReader = csv.NewReader(csvFile)
	csvReader.Comma = parser.CsvSeparator
	csvRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var resultType = GetMetaType(parser.BindObject)
	if checkType(resultType) {
		return nil, errors.New(fmt.Sprintf("type %v not supported", resultType.Name()))
	}
	results := reflect.New(reflect.SliceOf(reflect.PtrTo(resultType)))

	headers := csvRows[0]
	body := csvRows[1:]
	var csvField = make(map[string]int)
	for _, col := range headers {
		for j := 0; j < resultType.NumField(); j+=1 {
			field := resultType.Field(j)
			tag := field.Tag.Get("csv")
			if col == tag {
				csvField[col] = j
			}
		}
	}


	for _, csvRow := range body {
		obj := reflect.New(resultType)
		for j, csvCol := range csvRow {
			colName := headers[j]
			idx, ok := csvField[colName]
			if !ok {
				continue
			}
			currentField := obj.Elem().Field(idx)
			if parser.Setter != nil && parser.Setter(currentField, colName, csvCol) {
				continue
			}else {
				setField(currentField, csvCol, true)
			}
		}
		ele := reflect.Append(results.Elem(), obj)
		results.Elem().Set(ele)
	}
	return results.Interface(), err
}


// 获取obj的反射类型, 如果obj是指针，则返回指向的类型
func GetMetaType(obj interface{}) reflect.Type {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return reflect.ValueOf(obj).Elem().Type()
	}
	return reflect.TypeOf(obj)
}

func checkType(p reflect.Type) bool {
	if p.Kind() != reflect.Struct {
		return false
	}
	return true
}
