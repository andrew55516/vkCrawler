package helpers

import (
	"reflect"
	"testing"
	"time"
)

func TestStrToTime(t *testing.T) {
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
				str: "29 янв 2023 в 12:05",
			},
			want: time.Date(2023, time.January, 29, 18, 5, 0, 0, time.UTC),
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
