package kv_test

import (
	"reflect"
	"testing"

	"cos-backend-com/src/common/kv"
	"cos-backend-com/src/common/proto"
)

func TestParser_parseKV(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData proto.Data
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "simple test",
			fields: fields{
				name: "parser"},
			args: args{"key=v"},
			wantData: proto.Data{
				"key": "v",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
		{
			name: "normal test",
			fields: fields{
				name: "parser"},
			args: args{`key=v	a="b"	c=["adsf", "asdf=asdf"]	d={"position":"engineer",salary:15000}`},
			wantData: proto.Data{
				"key": "v",
				"a":   `"b"`,
				"c":   `["adsf", "asdf=asdf"]`,
				"d":   `{"position":"engineer",salary:15000}`,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
		{
			name: "two continous \\t",
			fields: fields{
				name: "parser"},
			args: args{`key=v	a="b"		c="dfafasdf"`},
			wantData: proto.Data{
				"key": "v",
				"a":   `"b"`,
				"c":   `"dfafasdf"`,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
		{
			name: "\\t in value, expect error",
			fields: fields{
				name: "parser"},
			args: args{`key=v	a="b"		c="dfaf	asdf"	d=123`},
			wantData: nil,
			wantErr:  true,
		},
		// TODO: Add test cases.
		{
			name: "null value",
			fields: fields{
				name: "parser"},
			args: args{`key=v	a=	c="c"`},
			wantData: proto.Data{
				"key": "v",
				"a":   "",
				"c":   `"c"`,
			},
			wantErr: false,
		},
		// test escape
		{
			name: "escape",
			fields: fields{
				name: "parser"},
			args: args{`key=v	a=	c=\"c\"`},
			wantData: proto.Data{
				"key": "v",
				"a":   "",
				"c":   `"c"`,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
		{
			name: "long value",
			fields: fields{
				name: "parser"},
			args: args{`agentKey=51af9e27547a295e2d1ae7df38520ae8	name=1_power_usage	plantID=9793245517840385	time=1533582602000000000	value={ "eventType": "out", "time": "2018-08-10T06:00:52.123Z", "plantId": "123456", "regionId": "321", "materialId": "111111", "supplierId": "22222", "operatorId": "33333", "amount": 100 }`},
			wantData: proto.Data{
				"agentKey": "51af9e27547a295e2d1ae7df38520ae8",
				"name":     "1_power_usage",
				"plantID":  `9793245517840385`,
				"time":     `1533582602000000000`,
				"value":    `{ "eventType": "out", "time": "2018-08-10T06:00:52.123Z", "plantId": "123456", "regionId": "321", "materialId": "111111", "supplierId": "22222", "operatorId": "33333", "amount": 100 }`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := kv.Decode([]byte(tt.args.line))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil {
				if !reflect.DeepEqual(gotData, tt.wantData) {
					t.Errorf("Parser.parseLine() = %v, want %v", gotData, tt.wantData)
				}
			}
		})
	}
}
