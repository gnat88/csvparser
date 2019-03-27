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
	X int	`json:"x"`
	Y int	`json:"y"`
}
type User struct {
	Id int32 `csv:"id"`
	Name string `csv:"name"`
	Vip int32 `csv:"vip"`
	Ext Extra `csv:"ext"`
}

type sLoader struct {

}

func (s *sLoader)Reader() (io.Reader, error) {
	return bytes.NewBufferString(`id;name;vip;ext
1;tang;4;"{""x"":1,""y"":2}"`), nil
}
func TestMain(m *testing.M) {
	csvParser = CsvParser{
		//Loader:NewFileLoader("example_files/example.csv"),
		Loader:&sLoader{},
		CsvSeparator: ';',
		BindObject:  User{},
		Setter:      LoadUser,
	}

	var ret interface{}
	ret, parseErr1 = csvParser.Parse()

	fmt.Printf("%v", (*(ret.(*[]*User)))[0])


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
