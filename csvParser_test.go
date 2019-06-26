package parser

import (
	"bytes"
	"fmt"
	"io"
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
	Type int	`json:"type"`
	Subtype int	`json:"subtype"`
	Amount int `json:"amount"`
}
type User struct {
	Id int32 `csv:"id"`
	Name string `csv:"name"`
	Vip int32 `csv:"vip"`
	Cols []int32 `csv:"cols"`
	Ext []Extra `csv:"ext"`
}

type sLoader struct {

}

func (s *sLoader)Reader() (io.Reader, error) {
	return bytes.NewBufferString(`id;name;vip;ext;cols
1;tang;4;"[{""type"":1, ""amount"":2},{""type"":3, ""subtype"":10001,""amount"":4}]";"[1,2,3]"`), nil
}
func TestMain(m *testing.M) {
	csvParser = CsvParser{
		//CsvReader:NewFileLoader("example_files/example.csv"),
		CsvReader:&sLoader{},
		CsvSeparator: ';',
		BindObject:  User{},
		Setter:      LoadUser,
	}

	var ret interface{}
	ret, parseErr1 = csvParser.Parse()
	x := *(ret.(*[]*User))

	fmt.Printf("%v", *x[0])


	//run all the tests
	os.Exit(m.Run())
}


func LoadUser(val reflect.Value, colName, raw string) bool {
	if colName == "vip" {
		i, _ := toInt(raw)
		val.SetInt(i*2)
		return true
	}

	return false
}
