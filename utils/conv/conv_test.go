package conv

import (
	"reflect"
	"testing"
)

func TestValueToString(t *testing.T) {
	var str string
	var err error

	str, err = ValueToString(100.2)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[float64](str); err != nil || val != 100.2 {
		t.Error(val, err)
	}

	str, err = ValueToString(200)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[int64](str); err != nil || val != 200 {
		t.Error(val, err)
	}

	str, err = ValueToString(true)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[bool](str); err != nil || val != true {
		t.Error(val, err)
	}

	str, err = ValueToString("this is str")
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[string](str); err != nil || val != "this is str" {
		t.Error(val, err)
	}

	st := testDataStruct{
		Field1: "f1",
		Field2: 1000,
		Field3: 2,
	}
	str, err = ValueToString(st)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[testDataStruct](str); err != nil || !reflect.DeepEqual(val, st) {
		t.Error(val, err)
	}

	if val, err := StringToValue[*testDataStruct](str); err != nil || !reflect.DeepEqual(*val, st) {
		t.Error(val, err)
	}

	mst := map[string]*testDataStruct{"one": &st}
	str, err = ValueToString(mst)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[map[string]*testDataStruct](str); err != nil || !reflect.DeepEqual(val, mst) {
		t.Error(val, err)
	}

	slice := []string{"a", "b", "c"}
	str, err = ValueToString(slice)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[[]string](str); err != nil || !reflect.DeepEqual(val, slice) {
		t.Error(val, err)
	}

	sample := &ConvSample{
		Field1: "192.168.0.1",
		Field2: "192.168.0.100",
		Ts:     12345,
	}
	str, err = ValueToString(sample)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[ConvSample](str); err != nil || val.Field1 != sample.Field1 || val.Field2 != sample.Field2 || val.Ts != sample.Ts {
		t.Error(val, err)
	}
}

type testDataStruct struct {
	Field1 string  `json:"field1" binding:"required"`
	Field2 int     `json:"field2" `
	Field3 float64 `json:"field3" `
}
