package bpe_test

import (
	"testing"

	"github.com/sodawave/VeGO/src/pkg/bpe"
	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func TestEstimateAndCompare(t *testing.T) {
	goSrc := []byte("package main\n\nfunc main() {\n\tprintln(1)\n}\n")
	if bpe.EstimateTokens(goSrc) == 0 {
		t.Fatal("expected tokens")
	}
	vego, err := transpiler.New(nil).Transform(goSrc, transpiler.ToVeGo)
	if err != nil {
		t.Fatal(err)
	}
	rep := bpe.Compare(goSrc, vego)
	if rep.GoTokens <= 0 || rep.VeGoTokens <= 0 {
		t.Fatalf("bad report: %+v", rep)
	}
	if rep.Estimator == "" {
		t.Fatal("missing estimator")
	}
}
