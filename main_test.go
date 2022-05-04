package main

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	base := "0124"
	{
		input := bytes.Repeat([]byte(base), 80/len(base))
		copy(input[1:3], []byte("01"))
		got, err := parse(input)
		if err != nil {
			log.Fatalf("error : %+v", err)
		}
		want := &info{
			ID:     "0",
			Type:   "01",
			Amount: "4012401",
			Detail: &userDetail{
				Title:    "2401240124",
				UserName: "01240124012401240124",
				Padding:  "01240124012401240124",
			},
			ReportTitle: "01240124012401240124",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("parse() = %v, want %v", got, want)
		}
	}
	{
		input := bytes.Repeat([]byte(base), 80/len(base))
		copy(input[1:3], []byte("02"))
		got, err := parse(input)
		if err != nil {
			log.Fatalf("error : %+v", err)
		}
		want := &info{
			ID:     "0",
			Type:   "02",
			Amount: "4012401",
			Detail: &itemDetail{
				Title:    "2401240124",
				ItemName: "0124012401",
				Padding:  "240124012401240124012401240124",
			},
			ReportTitle: "01240124012401240124",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("parse() = %v, want %v", got, want)
		}
	}
}
