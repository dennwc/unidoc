package fonts

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	"github.com/unidoc/unidoc/pdf/model/textencoding"
)

const fontDir = `../../creator/testdata`

var casesTTFParse = []struct {
	path         string
	name         string
	bold         bool
	italicAngle  float32
	underlinePos int16
	underlineTh  int16
	isFixed      bool
	bbox         [4]int
	runes        map[rune]uint16
	widths       map[rune]int
}{
	{
		path:         "FreeSans.ttf",
		name:         "FreeSans",
		underlinePos: -151,
		underlineTh:  50,
		bbox:         [4]int{-631, 1632, -462, 1230},
		runes: map[rune]uint16{
			'x': 0x5d,
			'ё': 0x32a,
		},
		widths: map[rune]int{
			'x': 500,
			'ё': 556,
		},
	},
	{
		path:         "roboto/Roboto-Bold.ttf",
		name:         "Roboto-Bold",
		bold:         true,
		underlinePos: -150,
		underlineTh:  100,
		bbox:         [4]int{-1488, 2439, -555, 2163},
		runes: map[rune]uint16{
			'x': 0x5c,
			'ё': 0x3cb,
		},
		widths: map[rune]int{
			'x': 1042,
			'ё': 1107,
		},
	},
	{
		path:         "roboto/Roboto-BoldItalic.ttf",
		name:         "Roboto-BoldItalic",
		bold:         true,
		italicAngle:  -12,
		underlinePos: -150,
		underlineTh:  100,
		bbox:         [4]int{-1459, 2467, -555, 2163},
		runes: map[rune]uint16{
			'x': 0x5c,
			'ё': 0x3cb,
		},
		widths: map[rune]int{
			'x': 1021,
			'ё': 1084,
		},
	},
}

var testRunes = []rune{'x', 'ё'}

func TestTTFParse(t *testing.T) {
	for _, c := range casesTTFParse {
		t.Run(c.path, func(t *testing.T) {
			path := filepath.Join(fontDir, c.path)
			t.Run("unidoc", func(t *testing.T) {
				ft, err := TtfParse(path)
				if err != nil {
					t.Fatal(err)
				}
				if ft.Bold != c.bold {
					t.Error(ft.Bold, c.bold)
				}
				if float32(ft.ItalicAngle) != c.italicAngle {
					t.Error(ft.ItalicAngle, c.italicAngle)
				}
				if ft.UnderlinePosition != c.underlinePos {
					t.Error(ft.UnderlinePosition, c.underlinePos)
				}
				if ft.UnderlineThickness != c.underlineTh {
					t.Error(ft.UnderlineThickness, c.underlineTh)
				}
				if ft.IsFixedPitch != c.isFixed {
					t.Error(ft.IsFixedPitch, c.isFixed)
				}
				if b := [4]int{int(ft.Xmin), int(ft.Xmax), int(ft.Ymin), int(ft.Ymax)}; b != c.bbox {
					t.Error(b, c.bbox)
				}
				if ft.PostScriptName != c.name {
					t.Errorf("%q %q", ft.PostScriptName, c.name)
				}
				enc := textencoding.NewTrueTypeFontEncoder(ft.Chars)

				for _, r := range testRunes {
					t.Run(string(r), func(t *testing.T) {
						ind, ok := enc.RuneToCharcode(r)
						if !ok {
							t.Fatal("no char")
						} else if ind != c.runes[r] {
							t.Fatalf("%x != %x", ind, c.runes[r])
						}
						w := ft.Widths[ft.Chars[uint16(r)]]
						if int(w) != c.widths[r] {
							t.Errorf("%d != %d", int(w), c.widths[r])
						}
					})
				}
			})
			t.Run("x", func(t *testing.T) {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}

				ft, err := sfnt.Parse(data)
				if err != nil {
					t.Fatal(err)
				}
				if ft.Selection().Bold() != c.bold {
					t.Error(ft.Selection().Bold(), c.bold)
				}
				post := ft.Post()
				epost := sfnt.PostInfo{
					ItalicAngle:        c.italicAngle,
					UnderlinePosition:  c.underlinePos,
					UnderlineThickness: c.underlineTh,
					IsFixedPitch:       c.isFixed,
				}
				if post != epost {
					t.Error(post, epost)
				}
				ppem := fixed.Int26_6(ft.UnitsPerEm())
				hint := font.HintingNone

				var buf sfnt.Buffer
				r, err := ft.Bounds(&buf, ppem, hint)
				if err != nil {
					t.Fatal(err)
				}
				// TODO: which standard to use?
				r.Min.Y, r.Max.Y = -r.Max.Y, -r.Min.Y
				if b := [4]int{int(r.Min.X), int(r.Max.X), int(r.Min.Y), int(r.Max.Y)}; b != c.bbox {
					t.Error(b, c.bbox)
				}
				name, err := ft.Name(&buf, sfnt.NameIDPostScript)
				if err != nil {
					t.Fatal(err)
				}
				if name != c.name {
					t.Errorf("%q %q", name, c.name)
				}

				for _, r := range testRunes {
					t.Run(string(r), func(t *testing.T) {
						ind, err := ft.GlyphIndex(&buf, r)
						if err != nil || ind == 0 {
							t.Fatal(err)
						} else if uint16(ind) != c.runes[r] {
							t.Errorf("%x != %x", ind, c.runes[r])
						}
						w, err := ft.GlyphAdvance(&buf, ind, ppem, hint)
						if err != nil {
							t.Fatal(err)
						} else if int(w) != c.widths[r] {
							t.Errorf("%d != %d", int(w), c.widths[r])
						}
					})
				}
			})
		})
	}
}
