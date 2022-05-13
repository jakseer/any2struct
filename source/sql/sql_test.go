package sql

import (
	"github.com/jakseer/any2struct/convert"
	"reflect"
	"testing"
)

func TestSource_Convert(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name    string
		args    args
		want    *convert.Struct
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{sql: ""},
			want:    nil,
			wantErr: true,
		}, {
			name: "case1",
			args: args{sql: "CREATE TABLE `t` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'primary key') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='comment';"},
			want: &convert.Struct{
				Name:    "t",
				Comment: "",
				Fields: []convert.StructField{
					{
						Key: "Id",
						Typ: convert.FieldTyp{
							Ptr: nil,
							Typ: convert.Int,
						},
						Tags:    nil,
						Comment: "primary key",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Source{}
			got, err := s.Convert(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSQLType(t *testing.T) {
	type args struct {
		sqlType string
	}
	tests := []struct {
		name string
		args args
		want convert.FieldTyp
	}{
		{
			name: "int",
			args: args{sqlType: "int"},
			want: convert.FieldTyp{
				Ptr: nil,
				Typ: convert.Int,
			},
		}, {
			name: "varchar",
			args: args{sqlType: "varchar"},
			want: convert.FieldTyp{
				Ptr: nil,
				Typ: convert.String,
			},
		}, {
			name: "unknown",
			args: args{sqlType: "xxx"},
			want: convert.FieldTyp{
				Ptr: nil,
				Typ: convert.Unknown,
			},
		}, {
			name: "empty",
			args: args{sqlType: ""},
			want: convert.FieldTyp{
				Ptr: nil,
				Typ: convert.Unknown,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSQLType(tt.args.sqlType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSQLType() = %v, want %v", got, tt.want)
			}
		})
	}
}
