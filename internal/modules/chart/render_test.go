package chart

import (
	"go.uber.org/zap"
	"image/color"
	"reflect"
	"testing"
)

func Test_parseGroup(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr bool
	}{
		{
			name: "one",
			args: args{
				s: "01",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "255",
			args: args{
				s: "FF",
			},
			want:    255,
			wantErr: false,
		},
		{
			name: "short",
			args: args{
				s: "F",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "long",
			args: args{
				s: "FFF",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "wrong",
			args: args{
				s: "FJ",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGroup(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseColor6(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    color.RGBA
		wantErr bool
	}{
		{
			name: "one",
			args: args{
				s: "00FF00",
			},
			want: color.RGBA{
				R: 0,
				G: 255,
				B: 0,
				A: 255,
			},
			wantErr: false,
		},
		{
			name: "two",
			args: args{
				s: "010203",
			},
			want: color.RGBA{
				R: 1,
				G: 2,
				B: 3,
				A: 255,
			},
			wantErr: false,
		},
		{
			name: "three",
			args: args{
				s: "F1F2F3",
			},
			want: color.RGBA{
				R: 241,
				G: 242,
				B: 243,
				A: 255,
			},
			wantErr: false,
		},
		{
			name: "short",
			args: args{
				s: "F1F2F",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor6(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor6() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseColor8(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    color.RGBA
		wantErr bool
	}{
		{
			name: "one",
			args: args{
				s: "F1F2F3F4",
			},
			want: color.RGBA{
				R: 241,
				G: 242,
				B: 243,
				A: 244,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor8(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor8() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChart_parseColor(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    color.RGBA
		wantErr bool
	}{
		{
			name: "six",
			args: args{
				s: "#010203",
			},
			want:    color.RGBA{R: 1, G: 2, B: 3, A: 255},
			wantErr: false,
		},
		{
			name: "eight",
			args: args{
				s: "#01020304",
			},
			want:    color.RGBA{R: 1, G: 2, B: 3, A: 4},
			wantErr: false,
		},
		{
			name: "wrong - long",
			args: args{
				s: "#010203046",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name: "wrong - short",
			args: args{
				s: "#01020",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name: "wrong - wrong sym",
			args: args{
				s: "#01020R",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name: "blue",
			args: args{
				s: "blue",
			},
			want:    color.RGBA{B: 255, A: 255},
			wantErr: false,
		},
		{
			name: "red",
			args: args{
				s: "red",
			},
			want:    color.RGBA{R: 255, A: 255},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &Chart{
				logger: zap.NewNop(),
			}
			got, err := ch.parseColor(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
