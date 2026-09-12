package transpiler_test

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

var multiNL = regexp.MustCompile(`\n{2,}`)

func normWS(src []byte) []byte {
	return multiNL.ReplaceAll(src, []byte("\n"))
}

func TestRoundTripFixtures(t *testing.T) {
	eng := transpiler.New(nil)
	files := []string{"hello_http.go", "struct_range.go"}
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			goSrc, err := os.ReadFile(filepath.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}
			canon, err := format.Source(goSrc)
			if err != nil {
				t.Fatal(err)
			}
			vego, err := eng.Transform(canon, transpiler.ToVeGo)
			if err != nil {
				t.Fatalf("ToVeGo: %v", err)
			}
			if bytes.Equal(vego, canon) {
				t.Fatal("expected .vego to differ from Go (keywords should be glyphs)")
			}
			if !bytes.Contains(vego, []byte("ƒ")) && !bytes.Contains(vego, []byte("ð")) {
				t.Fatalf("expected glyph keywords in .vego, got: %s", vego)
			}
			back, err := eng.Transform(vego, transpiler.ToGo)
			if err != nil {
				t.Fatalf("ToGo: %v\nvego=%q", err, vego)
			}
			// Semantic round-trip: ignore blank-line-only gofmt differences.
			if !bytes.Equal(normWS(back), normWS(canon)) {
				t.Fatalf("round-trip mismatch\n--- go ---\n%s\n--- back ---\n%s\n--- vego ---\n%q", canon, back, vego)
			}
		})
	}
}

func TestInvalidVeGoFailsClosed(t *testing.T) {
	eng := transpiler.New(nil)
	_, err := eng.Transform([]byte("ð not valid {{{"), transpiler.ToGo)
	if err == nil {
		t.Fatal("expected error for invalid .vego")
	}
}

func TestInvalidGoFailsClosed(t *testing.T) {
	eng := transpiler.New(nil)
	_, err := eng.Transform([]byte("package !!!"), transpiler.ToVeGo)
	if err == nil {
		t.Fatal("expected error for invalid Go")
	}
}
