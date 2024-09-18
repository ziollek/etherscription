package model

import (
	"reflect"
	"testing"
)

func TestShouldConvertHexStringToProperInt(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "empty string", args: args{hex: ""}, want: 0},
		{name: "not a number", args: args{hex: "test"}, want: 0},
		{name: "incorrect format", args: args{hex: "10"}, want: 0},
		{name: "incorrect format with proper prefix", args: args{hex: "0xx10"}, want: 0},
		{name: "proper like decimal value", args: args{hex: "0x10"}, want: 16},
		{name: "proper value containing letters", args: args{hex: "0x1f"}, want: 31},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertHexToInt(tt.args.hex); got != tt.want {
				t.Errorf("ConvertHexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldConvertRawToSimplifiedTransaction(t *testing.T) {
	type fields struct {
		From             string
		To               string
		Input            string
		Value            string
		TransactionIndex string
	}
	tests := []struct {
		name   string
		fields fields
		want   Transaction
	}{
		{name: "empty", fields: fields{}, want: Transaction{}},
		{name: "simple", fields: fields{From: "0x1", To: "0x2", Value: "0x10"}, want: Transaction{From: "0x1", To: "0x2", Value: 16}},
		{name: "complex", fields: fields{From: "0x1", To: "0x2", Value: "0x1f"}, want: Transaction{From: "0x1", To: "0x2", Value: 31}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			t := &RawTransaction{
				From:             tt.fields.From,
				To:               tt.fields.To,
				Input:            tt.fields.Input,
				Value:            tt.fields.Value,
				TransactionIndex: tt.fields.TransactionIndex,
			}
			if got := t.ToTransaction(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("ToTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
