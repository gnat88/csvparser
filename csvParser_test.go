package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

var contacts1 []interface{}
var contacts2 []interface{}
var contacts3 []interface{}

var parseErr1 error
var parseErr2 error
var parseErr3 error
var csvParser CsvParser

type Extra struct {
	X int	`json:"x"`
	Y int	`json:"y"`
}
type User struct {
	Id int32 `csv:"id"`
	Name string `csv:"name"`
	Vip int32 `csv:"vip"`
	Ext Extra `csv:"ext"`
}

func TestMain(m *testing.M) {

	e1 := &Extra{
		X:1,
		Y:2,
	}
	bts, _ := json.Marshal(e1)
	fmt.Println(string(bts))
	csvParser = CsvParser{
		CsvFile:      "example_files/example.csv",
		CsvSeparator: ';',
		BindObject:   reflect.TypeOf(User{}),
		Handler:      LoadUser,
	}

	var ret interface{}
	ret, parseErr1 = csvParser.Parse(nil)

	fmt.Printf("%v", (*(ret.(*[]*User)))[0])


	//run all the tests
	os.Exit(m.Run())
}


func LoadUser(colName string, val reflect.Value, raw string) bool {
	if colName == "vip" {
		i, _ := toInt(raw)
		val.SetInt(i*2)
		return true
	}

	return false
}
