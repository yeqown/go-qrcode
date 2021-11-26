package qrcode

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

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
	v1 := loadVersion(1, ErrorCorrectionMedium)
	v2 := loadVersion(5, ErrorCorrectionMedium)
	v3 := loadVersion(23, ErrorCorrectionMedium)

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
			want:    &v1,
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				raw:   []byte(strings.Repeat("TEXT", 30)),
				ecLv:  ErrorCorrectionMedium,
				eMode: EncModeAlphanumeric,
			},
			want:    &v2,
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				raw:   []byte(strings.Repeat("TEXT", 300)),
				ecLv:  ErrorCorrectionMedium,
				eMode: EncModeAlphanumeric,
			},
			want:    &v3,
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

func Test_binarySearchVersion(t *testing.T) {
	t.Logf("length of versions: %d", len(versions))

	find := func(v int, ec ecLevel) func(cursor version) int {
		return func(cursor version) int {
			if cursor.Ver == v && cursor.ECLevel == ec {
				return 0
			}

			if cursor.Ver > v || cursor.ECLevel > ec {
				return -1
			}

			return 1
		}
	}

	tests := []struct {
		name string
		ecLv ecLevel
		v    int
		want int
	}{
		{
			name: "case 0",
			ecLv: ErrorCorrectionLow,
			v:    1,
			want: 0,
		},
		{
			name: "case 1",
			ecLv: ErrorCorrectionHighest,
			v:    40,
			want: 159,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := binarySearchVersion(tt.v, find(tt.v, tt.ecLv))
			require.True(t, found)
			require.Equal(t, versions[tt.want], got)
		})
	}
}

func Test_binarySearchVersion_all(t *testing.T) {
	for _, v := range versions {
		hit, found := binarySearchVersion(v.Ver, func(cursor version) int {
			if cursor.Ver == v.Ver && cursor.ECLevel == v.ECLevel {
				return 0
			}

			// less
			if cursor.Ver > v.Ver || cursor.ECLevel > v.ECLevel {
				return -1
			}

			return 1
		})

		if !found {
			t.Errorf("binarySearchVersions() failed to find version %d", v.Ver)
		}
		assert.Equal(t, v, hit)
	}
}

// // go test -run=NONE -bench . -count 10 > new/old.txt
func Benchmark_loadVersion_top(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadVersion(2, ErrorCorrectionMedium)
		loadVersion(5, ErrorCorrectionMedium)
	}
}

func Benchmark_loadVersion_waist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadVersion(25, ErrorCorrectionMedium)
		loadVersion(15, ErrorCorrectionMedium)
	}
}

func Benchmark_loadVersion_bottom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadVersion(40, ErrorCorrectionHighest)
		loadVersion(35, ErrorCorrectionHighest)
	}
}

func Benchmark_analyzeVersion_short(b *testing.B) {
	source := []byte("text")

	for i := 0; i < b.N; i++ {
		_, _ = analyzeVersion(source, ErrorCorrectionMedium, EncModeByte)
	}
}

func Benchmark_analyzeVersion_middle(b *testing.B) {
	source := []byte(strings.Repeat("text", 30))

	for i := 0; i < b.N; i++ {
		_, _ = analyzeVersion(source, ErrorCorrectionMedium, EncModeByte)
	}
}

func Benchmark_analyzeVersion_long(b *testing.B) {
	source := []byte(strings.Repeat("text", 300))

	for i := 0; i < b.N; i++ {
		_, _ = analyzeVersion(source, ErrorCorrectionMedium, EncModeByte)
	}
}
