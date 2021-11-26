package qrcode

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion_Dimension(t *testing.T) {
	r := rand.Intn(len(versions))
	d := versions[r].Dimension()

	assert.Equal(t, versions[r].Ver*4+17, d)
}

func Test_loadVersion(t *testing.T) {
	// load(defaultVersionCfg)
	type args struct {
		lv         int
		recoveryLv ecLevel
	}
	tests := []struct {
		name string
		args args
		want version
	}{
		{
			name: "case 0",
			args: args{
				lv:         1,
				recoveryLv: ErrorCorrectionHighest,
			},
			want: version{
				Ver:     1,
				ECLevel: ErrorCorrectionHighest,
				Cap: capacity{
					Numeric:      17,
					AlphaNumeric: 10,
					Byte:         7,
					JP:           4,
				},
				RemainderBits: 0,
				Groups: []group{
					{
						NumBlocks:            1,
						NumDataCodewords:     9,
						ECBlockwordsPerBlock: 17,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadVersion(tt.args.lv, tt.args.recoveryLv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_analyzeVersion(t *testing.T) {
	// load(defaultVersionCfg)
	v := loadVersion(1, ErrorCorrectionMedium)
	type args struct {
		raw   []byte
		ecLv  ecLevel
		eMode encMode
	}
	tests := []struct {
		name    string
		args    args
		want    *version
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				raw:   []byte("TEXT"),
				ecLv:  ErrorCorrectionMedium,
				eMode: EncModeAlphanumeric,
			},
			want:    &v,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := analyzeVersion(tt.args.raw, tt.args.ecLv, tt.args.eMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyzeVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("analyzeVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
