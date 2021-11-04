package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func loadCSV(p string) (*csv.Reader, error) {
	r, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open csv file failed: %v", err)
	}
	return csv.NewReader(r), nil
}

func str2Int(s string) int {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		err = fmt.Errorf("parse string(%s) to int failed: %v", s, err)
		panic(err)
	}

	return v
}

func levelToInt(lv string) int {
	lv = strings.ToUpper(strings.TrimSpace(lv))
	switch lv {
	case "L":
		return 1
	case "M":
		return 2
	case "Q":
		return 3
	case "H":
		return 4
	default:
		panic("invalid EC Level:" + lv)
	}
}

func parseVersionsFromCSV(r *csv.Reader) ([]version, error) {
	const (
		verIdx             = 0
		lvIdx              = 1
		capNumericIdx      = 2
		capAlphaNumericIdx = 3
		capByteIdx         = 4
		capJPIdx           = 5
	)

	const (
		groupECBlockwordsPerBlock = 6

		group1NumBlocks        = 7
		group1NumDataCodewords = 8
		group2NumBlocks        = 9
		group2NumDataCodewords = 10
	)

	versions := make([]version, 0, 160)
	// skip first row
	row, err := r.Read()
	for true {
		row, err = r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(row)
		ver := str2Int(row[verIdx])
		v := version{
			Ver:           ver,
			ECLevel:       levelToInt(row[lvIdx]),
			RemainderBits: reminderBitsCache[ver],
			Capacity: capacity{
				Numeric:      str2Int(row[capNumericIdx]),
				AlphaNumeric: str2Int(row[capAlphaNumericIdx]),
				Byte:         str2Int(row[capByteIdx]),
				JP:           str2Int(row[capJPIdx]),
			},
		}

		ecblockwordsPerBlock := str2Int(row[groupECBlockwordsPerBlock])
		// DONE(@yeqown): determine how many group would be need by this version-EC
		groups := make([]group, 0, 2)
		group1 := group{
			NumBlocks:            str2Int(row[group1NumBlocks]),
			ECBlockwordsPerBlock: ecblockwordsPerBlock,
			NumDataCodewords:     str2Int(row[group1NumDataCodewords]),
		}
		// group1 can not be empty
		groups = append(groups, group1)
		group2 := group{
			NumBlocks:            str2Int(row[group2NumBlocks]),
			ECBlockwordsPerBlock: ecblockwordsPerBlock,
			NumDataCodewords:     str2Int(row[group2NumDataCodewords]),
		}

		if group2.NumBlocks != 0 &&
			group2.ECBlockwordsPerBlock != 0 &&
			group2.NumDataCodewords != 0 {
			groups = append(groups, group2)
		}

		v.Groups = groups
		versions = append(versions, v)
	}

	return versions, nil
}

type capacity struct {
	Numeric      int
	AlphaNumeric int
	Byte         int
	JP           int
}

type group struct {
	NumBlocks            int
	ECBlockwordsPerBlock int
	NumDataCodewords     int
}

type version struct {
	Ver           int
	ECLevel       int
	RemainderBits int
	Capacity      capacity
	Groups        []group
}

// fillVersionIntoConfigDotGo get tbl into template and save into version_cfg.go
//
// tbl is expected as "" from https://www.thonky.com/qr-code-tutorial/error-correction-table
func fillVersionIntoConfigDotGo(versions []version) []byte {
	w := bytes.NewBuffer(nil)
	t := template.Must(template.New("tplVersion").Parse(_tplVersions))

	err := t.Execute(w, versions)
	if err != nil {
		panic(err)
	}

	return w.Bytes()
}

var _tplVersions = `var versions = []version{
	{{range .}}
		{
			Ver:     {{.Ver}},
			ECLevel: {{.ECLevel}},
			Cap: capacity{
				Numeric:      {{.Capacity.Numeric}},
				AlphaNumeric: {{.Capacity.AlphaNumeric}},
				Byte:         {{.Capacity.Byte}},
				JP:           {{.Capacity.JP}},
			},
			RemainderBits: {{.RemainderBits}},
			Groups: []group{
				{{range $i, $v := .Groups}}
					{
						NumBlocks:            {{$v.NumBlocks}},
						NumDataCodewords:     {{$v.NumDataCodewords}},
						ECBlockwordsPerBlock: {{$v.ECBlockwordsPerBlock}},
					},
				{{end}}
			},
		},
	{{end}}
}`

func main() {
	r, err := loadCSV("./qr_code_version_configs.csv")
	if err != nil {
		panic(err)
	}

	versions, err := parseVersionsFromCSV(r)
	if err != nil {
		panic(err)
	}

	//byt, _ := json.MarshalIndent(versions, "", "\t")
	//fmt.Printf("%s", byt)
	//
	//w, err := os.Open("./tmp.json")
	//if err != nil {
	//	panic(err)
	//}

	byts := fillVersionIntoConfigDotGo(versions)
	if err = ioutil.WriteFile("./tmp.go", byts, 0644); err != nil {
		panic(err)
	}
}

// https://www.thonky.com/qr-code-tutorial/structure-final-message
var reminderBitsCache = map[int]int{
	1:  0,
	2:  7,
	3:  7,
	4:  7,
	5:  7,
	6:  7,
	7:  0,
	8:  0,
	9:  0,
	10: 0,
	11: 0,
	12: 0,
	13: 0,
	14: 3,
	15: 3,
	16: 3,
	17: 3,
	18: 3,
	19: 3,
	20: 3,
	21: 4,
	22: 4,
	23: 4,
	24: 4,
	25: 4,
	26: 4,
	27: 4,
	28: 3,
	29: 3,
	30: 3,
	31: 3,
	32: 3,
	33: 3,
	34: 3,
	35: 0,
	36: 0,
	37: 0,
	38: 0,
	39: 0,
	40: 0,
}
