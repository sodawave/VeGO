package bpe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sodawave/VeGO/src/pkg/bpe"
	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func TestCompareTiktoken(t *testing.T) {
	goSrc := []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hi\")\n}\n")
	vego, err := transpiler.New(nil).Transform(goSrc, transpiler.ToVeGo)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := bpe.Compare(goSrc, vego)
	if err != nil {
		t.Fatal(err)
	}
	if rep.GoTokens == 0 || rep.VeGoTokens == 0 {
		t.Fatalf("zero tokens: %+v", rep)
	}
	if rep.Estimator != "tiktoken/cl100k_base" {
		t.Fatalf("estimator: %s", rep.Estimator)
	}
	t.Logf("go=%d vego=%d save=%.1f%% bytes=%.1f%% ir=%s", rep.GoTokens, rep.VeGoTokens, rep.TokenSaving, rep.ByteSaving, vego)
}

func TestExampleFilesSaveOrParity(t *testing.T) {
	eng := transpiler.New(nil)
	root := filepath.Join("..", "..", "examples")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		matches, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		if len(matches) == 0 {
			continue
		}
		goSrc, err := os.ReadFile(matches[0])
		if err != nil {
			t.Fatal(err)
		}
		vego, err := eng.Transform(goSrc, transpiler.ToVeGo)
		if err != nil {
			t.Fatalf("%s: %v", e.Name(), err)
		}
		rep, err := bpe.Compare(goSrc, vego)
		if err != nil {
			t.Fatal(err)
		}
		// Optimal line: never regress tokens vs Go for Beta examples.
		if rep.VeGoTokens > rep.GoTokens {
			t.Fatalf("%s: token regression go=%d vego=%d ir=%s", e.Name(), rep.GoTokens, rep.VeGoTokens, vego)
		}
		t.Logf("%s: tokens %d→%d (%.1f%%) bytes %.1f%%", e.Name(), rep.GoTokens, rep.VeGoTokens, rep.TokenSaving, rep.ByteSaving)
	}
}
