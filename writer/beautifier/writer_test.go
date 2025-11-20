package beautifier

import (
	"testing"

	"github.com/yeqown/go-qrcode/v2"
)

func TestWriter_GenerateSVG_Styles(t *testing.T) {
	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		t.Fatalf("failed to create qrcode: %v", err)
	}

	tests := []struct {
		name  string
		style Style
	}{
		{
			name: "test_style_base.svg",
			style: BaseStyle{Color: "#000000"},
		},
		{
			name: "test_style_liquid.svg",
			style: LiquidStyle{Color: "#000000"}, // LiquidStyle now handles Finders as circles
		},
		{
			name: "test_style_geometric.svg",
			style: BaseStyle{Color: "#0000FF"}, // BaseStyle handles geometric shapes by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := New(tt.name, tt.style)
			if err != nil {
				t.Fatalf("failed to create writer: %v", err)
			}
			defer w.Close()

			if err := qrc.Save(w); err != nil {
				t.Fatalf("failed to save qrcode: %v", err)
			}
		})
	}
}
