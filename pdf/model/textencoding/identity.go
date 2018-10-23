/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package textencoding

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/unidoc/unidoc/pdf/core"
)

// IdentityEncoder represents an 2-byte identity encoding
type IdentityEncoder struct {
	baseName string
}

// NewIdentityTextEncoder returns a new IdentityEncoder based on predefined
// encoding `baseName` and difference map `differences`.
func NewIdentityTextEncoder(baseName string) IdentityEncoder {
	return IdentityEncoder{baseName}
}

// String returns a string that describes `enc`.
func (enc IdentityEncoder) String() string {
	return enc.baseName
}

// Encode converts the Go unicode string `raw` to a PDF encoded string.
func (enc IdentityEncoder) Encode(raw string) []byte {
	return doEncode2(enc, raw)
}

// CharcodeToGlyph returns the glyph name matching character code `code`.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) CharcodeToGlyph(code uint16) (string, bool) {
	r, found := enc.CharcodeToRune(code)
	if found && r == 0x20 {
		return "space", true
	}

	// Returns "uniXXXX" format where XXXX is the code in hex format.
	glyph := fmt.Sprintf("uni%.4X", code)
	return glyph, true
}

// GlyphToCharcode returns the character code matching glyph `glyph`.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) GlyphToCharcode(glyph string) (uint16, bool) {
	r, ok := enc.GlyphToRune(glyph)
	if !ok {
		return 0, false
	}
	return enc.RuneToCharcode(r)
}

// RuneToCharcode converts rune `r` to a PDF character code.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) RuneToCharcode(r rune) (uint16, bool) {
	return uint16(r), true
}

// CharcodeToRune converts PDF character code `code` to a rune.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) CharcodeToRune(code uint16) (rune, bool) {
	return rune(code), true
}

// RuneToGlyph returns the glyph name for rune `r`.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) RuneToGlyph(r rune) (string, bool) {
	if r == ' ' {
		return "space", true
	}
	glyph := fmt.Sprintf("uni%.4X", r)
	return glyph, true
}

// GlyphToRune returns the rune corresponding to glyph name `glyph`.
// The bool return flag is true if there was a match, and false otherwise.
func (enc IdentityEncoder) GlyphToRune(glyph string) (rune, bool) {
	// String with "uniXXXX" format where XXXX is the hexcode.
	if glyph == "space" {
		return ' ', true
	} else if !strings.HasPrefix(glyph, "uni") || len(glyph) != 7 {
		return 0, false
	}
	r, err := strconv.ParseUint(glyph[3:], 16, 16)
	if err != nil {
		return 0, false
	}
	return rune(r), true
}

// ToPdfObject returns a nil as it is not truly a PDF object and should not be attempted to store in file.
func (enc IdentityEncoder) ToPdfObject() core.PdfObject {
	if enc.baseName != "" {
		return core.MakeName(enc.baseName)
	}
	return core.MakeNull()
}
