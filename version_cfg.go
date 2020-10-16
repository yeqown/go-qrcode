package qrcode

var (
	versions = []version{
		{
			Ver:     1,
			ECLevel: 0,
			Cap: capacity{
				Numeric:      41,
				AlphaNumeric: 25,
				Byte:         17,
				JP:           10,
			},
			RemainderBits: 0,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     19,
					ECBlockwordsPerBlock: 7,
				},
			},
		},
		{
			Ver:     1,
			ECLevel: 1,
			Cap: capacity{
				Numeric:      34,
				AlphaNumeric: 20,
				Byte:         14,
				JP:           8,
			},
			RemainderBits: 0,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     16,
					ECBlockwordsPerBlock: 10,
				},
			},
		},
		{
			Ver:     1,
			ECLevel: 2,
			Cap: capacity{
				Numeric:      27,
				AlphaNumeric: 16,
				Byte:         11,
				JP:           7,
			},
			RemainderBits: 0,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     13,
					ECBlockwordsPerBlock: 13,
				},
			},
		},
		{
			Ver:     1,
			ECLevel: 3,
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
		{
			Ver:     2,
			ECLevel: 0,
			Cap: capacity{
				Numeric:      77,
				AlphaNumeric: 47,
				Byte:         32,
				JP:           20,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     34,
					ECBlockwordsPerBlock: 10,
				},
			},
		},
		{
			Ver:     2,
			ECLevel: 1,
			Cap: capacity{
				Numeric:      63,
				AlphaNumeric: 38,
				Byte:         26,
				JP:           16,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     34,
					ECBlockwordsPerBlock: 16,
				},
			},
		},
		{
			Ver:     2,
			ECLevel: 2,
			Cap: capacity{
				Numeric:      48,
				AlphaNumeric: 29,
				Byte:         20,
				JP:           12,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     22,
					ECBlockwordsPerBlock: 22,
				},
			},
		},
		{
			Ver:     2,
			ECLevel: 3,
			Cap: capacity{
				Numeric:      34,
				AlphaNumeric: 20,
				Byte:         14,
				JP:           8,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     16,
					ECBlockwordsPerBlock: 28,
				},
			},
		},
		{
			Ver:     3,
			ECLevel: 0,
			Cap: capacity{
				Numeric:      127,
				AlphaNumeric: 77,
				Byte:         53,
				JP:           32,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     55,
					ECBlockwordsPerBlock: 15,
				},
			},
		},
		{
			Ver:     3,
			ECLevel: 1,
			Cap: capacity{
				Numeric:      101,
				AlphaNumeric: 61,
				Byte:         42,
				JP:           26,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     44,
					ECBlockwordsPerBlock: 26,
				},
			},
		},
		{
			Ver:     3,
			ECLevel: 2,
			Cap: capacity{
				Numeric:      77,
				AlphaNumeric: 47,
				Byte:         32,
				JP:           20,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     17,
					ECBlockwordsPerBlock: 18,
				},
			},
		},
		{
			Ver:     3,
			ECLevel: 3,
			Cap: capacity{
				Numeric:      58,
				AlphaNumeric: 35,
				Byte:         24,
				JP:           15,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     13,
					ECBlockwordsPerBlock: 22,
				},
			},
		},
		{
			Ver:     4,
			ECLevel: 0,
			Cap: capacity{
				Numeric:      187,
				AlphaNumeric: 114,
				Byte:         78,
				JP:           48,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     80,
					ECBlockwordsPerBlock: 20,
				},
			},
		},
		{
			Ver:     4,
			ECLevel: 1,
			Cap: capacity{
				Numeric:      149,
				AlphaNumeric: 90,
				Byte:         62,
				JP:           38,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     32,
					ECBlockwordsPerBlock: 18,
				},
			},
		},
		{
			Ver:     4,
			ECLevel: 2,
			Cap: capacity{
				Numeric:      111,
				AlphaNumeric: 67,
				Byte:         46,
				JP:           28,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     24,
					ECBlockwordsPerBlock: 26,
				},
			},
		},
		{
			Ver:     4,
			ECLevel: 3,
			Cap: capacity{
				Numeric:      82,
				AlphaNumeric: 50,
				Byte:         34,
				JP:           21,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            4,
					NumDataCodewords:     9,
					ECBlockwordsPerBlock: 16,
				},
			},
		},
		{
			Ver:     5,
			ECLevel: 0,
			Cap: capacity{
				Numeric:      255,
				AlphaNumeric: 154,
				Byte:         106,
				JP:           65,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            1,
					NumDataCodewords:     108,
					ECBlockwordsPerBlock: 26,
				},
			},
		},
		{
			Ver:     5,
			ECLevel: 1,
			Cap: capacity{
				Numeric:      202,
				AlphaNumeric: 122,
				Byte:         84,
				JP:           52,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     43,
					ECBlockwordsPerBlock: 24,
				},
			},
		},
		{
			Ver:     5,
			ECLevel: 2,
			Cap: capacity{
				Numeric:      144,
				AlphaNumeric: 87,
				Byte:         60,
				JP:           37,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     15,
					ECBlockwordsPerBlock: 18,
				},
				{
					NumBlocks:            2,
					NumDataCodewords:     16,
					ECBlockwordsPerBlock: 18,
				},
			},
		},
		{
			Ver:     5,
			ECLevel: 3,
			Cap: capacity{
				Numeric:      106,
				AlphaNumeric: 64,
				Byte:         44,
				JP:           27,
			},
			RemainderBits: 7,
			Groups: []group{
				{
					NumBlocks:            2,
					NumDataCodewords:     11,
					ECBlockwordsPerBlock: 22,
				},
				{
					NumBlocks:            2,
					NumDataCodewords:     12,
					ECBlockwordsPerBlock: 22,
				},
			},
		},
	}
)
