package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type detail interface {
	detail()
}
type userDetail struct {
	Title    string `bytegp:"length:10;offset:0"`
	UserName string `bytegp:"length:20;offset:10"`
	Padding  string `bytegp:"length:20;offset:30"`
}

func (*userDetail) detail() {}

type itemDetail struct {
	Title    string `bytegp:"length:10;offset:0"`
	ItemName string `bytegp:"length:10;offset:10"`
	Padding  string `bytegp:"length:30;offset:20"`
}

func (*itemDetail) detail() {}

type detailType string

const (
	detailTypeUser detailType = "01"
	detailTypeItem detailType = "02"
)

type info struct {
	ID          string     `bytegp:"length:1;offset:0"`
	Type        detailType `bytegp:"length:2;offset:1"`
	Amount      string     `bytegp:"length:7;offset:3"`
	Detail      detail     `bytegp:"length:50;offset:10"`
	ReportTitle string     `bytegp:"length:20;offset:60"`
}

func parseDetail(input []byte, typ detailType) (detail, error) {
	var (
		result detail
		v      reflect.Value
	)

	switch typ {
	case detailTypeUser:
		result = &userDetail{}
		v = reflect.ValueOf(result)
	case detailTypeItem:
		result = &itemDetail{}
		v = reflect.ValueOf(result)
	default:
		return nil, fmt.Errorf("invalid type(%s)", typ)
	}
	v = v.Elem()
	for i := 0; i < v.NumField(); i++ {
		tfield := v.Type().Field(i)
		l, o, err := parseTag(tfield.Tag.Get("bytegp"))
		if err != nil {
			return nil, err
		}
		v.Field(i).SetString(string(input[o : o+l])) // TODO check length
	}
	return result, nil
}

func parse(input []byte) (*info, error) {
	var result info
	v := reflect.ValueOf(&result).Elem()
	for i := 0; i < v.NumField(); i++ {
		tfield := v.Type().Field(i)
		l, o, err := parseTag(tfield.Tag.Get("bytegp"))
		if err != nil {
			return nil, err
		}
		vfield := v.Field(i)
		data := input[o : o+l] // TODO check length
		switch vfield.Kind() {
		case reflect.String:
			vfield.SetString(string(data))
		case reflect.Interface:
			result, err := parseDetail(data, result.Type)
			if err != nil {
				return nil, err
			}
			vfield.Set(reflect.ValueOf(result))
		default:
			fmt.Printf("%s\n", vfield.Kind())
		}
	}
	return &result, nil
}

func parseChildTag(childTag string) (string, int, error) {
	key, value, ok := strings.Cut(childTag, ":")
	if !ok {
		return "", 0, fmt.Errorf("invalid tag in %s", childTag)
	}
	numValue, err := strconv.Atoi(value)
	if err != nil {
		return "", 0, err
	}
	return key, numValue, nil
}

func parseTag(tag string) (int, int, error) {
	tagL, tagR, ok := strings.Cut(tag, ";")
	if !ok {
		return 0, 0, errors.New("invalid tag")
	}
	keyL, l, err := parseChildTag(tagL)
	if err != nil {
		return 0, 0, err
	}
	keyR, r, err := parseChildTag(tagR)
	if err != nil {
		return 0, 0, err
	}
	if keyL == "length" && keyR == "offset" {
		return l, r, nil
	}
	if keyL == "offset" && keyR == "length" {
		return r, l, nil
	}
	return 0, 0, errors.New("invalid")
}

func main() {
	base := "0124"
	input := bytes.Repeat([]byte(base), 80/len(base))
	copy(input[1:3], []byte("01"))
	result, err := parse(input)
	if err != nil {
		log.Fatalf("error : %+v", err)
	}
	fmt.Println(result)
	fmt.Println(result.Detail)
}
