package helpers

import (
	"reflect"
	"testing"
	"time"
)

func TestStrToTime(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Moscow")
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test",
			args: args{
				str: "17 окт 2022",
			},
			want: time.Date(2022, time.October, 17, 12, 0, 0, 0, location),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToTime(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
