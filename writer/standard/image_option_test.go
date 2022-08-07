package standard

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/yeqown/go-qrcode/v2"
)

func Test_stateRGBA(t *testing.T) {
	oo := defaultOutputImageOption()

	type args struct {
		v qrcode.QRValue
	}
	tests := []struct {
		name string
		args args
		want color.Color
	}{
		{
			name: "case 1",
			args: args{v: qrcode.QRValue_DATA_V0},
			want: oo.bgColor,
		},
		{
			name: "case 2",
			args: args{v: qrcode.QRValue_INIT_V0},
			want: oo.bgColor,
		},
		{
			name: "case 3",
			args: args{v: qrcode.QRValue_DATA_V1},
			want: oo.qrColor,
		},
		{
			name: "case 4",
			args: args{v: qrcode.QRValue_FORMAT_V1},
			want: oo.qrColor,
		},
		{
			name: "case 5",
			args: args{v: qrcode.QRValue_VERSION_V1},
			want: oo.qrColor,
		},
		{
			name: "case 6",
			args: args{v: qrcode.QRValue(0x0f)},
			want: oo.qrColor,
		},
		{
			name: "case 7",
			args: args{v: qrcode.QRValue_FINDER_V1},
			want: oo.qrColor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oo.translateQrColor(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("translateQrColor() = %v, want %v", got, tt.want)
			}
		})
	}

}
