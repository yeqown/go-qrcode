package qrcode

import (
	"reflect"
	"testing"

	"github.com/yeqown/reedsolomon/binary"
)

// func Test_load(t *testing.T) {
// 	type args struct {
// 		pathToCfg string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "case 0",
// 			args:    args{pathToCfg: defaultVersionCfg},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := load(tt.args.pathToCfg); (err != nil) != tt.wantErr {
// 				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestVersion_Dimension(t *testing.T) {
	type fields struct {
		Ver           int
		ECLevel       ECLevel
		Cap           capacity
		RemainderBits int
		Groups        []group
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Ver:           tt.fields.Ver,
				ECLevel:       tt.fields.ECLevel,
				Cap:           tt.fields.Cap,
				RemainderBits: tt.fields.RemainderBits,
				Groups:        tt.fields.Groups,
			}
			if got := v.Dimension(); got != tt.want {
				t.Errorf("Version.Dimension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_NumTotalCodewrods(t *testing.T) {
	type fields struct {
		Ver           int
		ECLevel       ECLevel
		Cap           capacity
		RemainderBits int
		Groups        []group
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Ver:           tt.fields.Ver,
				ECLevel:       tt.fields.ECLevel,
				Cap:           tt.fields.Cap,
				RemainderBits: tt.fields.RemainderBits,
				Groups:        tt.fields.Groups,
			}
			if got := v.NumTotalCodewrods(); got != tt.want {
				t.Errorf("Version.NumTotalCodewrods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_NumGroups(t *testing.T) {
	type fields struct {
		Ver           int
		ECLevel       ECLevel
		Cap           capacity
		RemainderBits int
		Groups        []group
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Ver:           tt.fields.Ver,
				ECLevel:       tt.fields.ECLevel,
				Cap:           tt.fields.Cap,
				RemainderBits: tt.fields.RemainderBits,
				Groups:        tt.fields.Groups,
			}
			if got := v.NumGroups(); got != tt.want {
				t.Errorf("Version.NumGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_verInfo(t *testing.T) {
	type fields struct {
		Ver           int
		ECLevel       ECLevel
		Cap           capacity
		RemainderBits int
		Groups        []group
	}
	tests := []struct {
		name   string
		fields fields
		want   *binary.Binary
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Ver:           tt.fields.Ver,
				ECLevel:       tt.fields.ECLevel,
				Cap:           tt.fields.Cap,
				RemainderBits: tt.fields.RemainderBits,
				Groups:        tt.fields.Groups,
			}
			if got := v.verInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.verInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_formatInfo(t *testing.T) {
	type fields struct {
		Ver           int
		ECLevel       ECLevel
		Cap           capacity
		RemainderBits int
		Groups        []group
	}
	type args struct {
		maskPattern int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *binary.Binary
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Ver:           tt.fields.Ver,
				ECLevel:       tt.fields.ECLevel,
				Cap:           tt.fields.Cap,
				RemainderBits: tt.fields.RemainderBits,
				Groups:        tt.fields.Groups,
			}
			if got := v.formatInfo(tt.args.maskPattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.formatInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadVersion(t *testing.T) {
	// load(defaultVersionCfg)
	type args struct {
		lv         int
		recoveryLv ECLevel
	}
	tests := []struct {
		name string
		args args
		want Version
	}{
		{
			name: "case 0",
			args: args{
				lv:         1,
				recoveryLv: Highest,
			},
			want: Version{
				Ver:     1,
				ECLevel: Highest,
				Cap: capacity{
					Numeric:      17,
					AlphaNumeric: 10,
					Byte:         7,
					JP:           4,
				},
				RemainderBits: 0,
				Groups: []group{
					group{
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
	v := loadVersion(1, Meddium)
	type args struct {
		raw   []byte
		ecLv  ECLevel
		eMode EncMode
	}
	tests := []struct {
		name    string
		args    args
		want    *Version
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				raw:   []byte("TEXT"),
				ecLv:  Meddium,
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
