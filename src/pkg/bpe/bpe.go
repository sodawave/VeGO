// Package bpe provides mid-term token-estimate helpers for VeGo compression goals.
// Alpha does not gate on these metrics (AD-5 / NFR-2).
package bpe

import (
	"unicode"
	"unicode/utf8"
)

// Report is a JSON-serializable comparison of Go vs .vego token estimates.
type Report struct {
	GoBytes      int     `json:"go_bytes"`
	VeGoBytes    int     `json:"vego_bytes"`
	GoTokens     int     `json:"go_tokens"`
	VeGoTokens   int     `json:"vego_tokens"`
	ByteRatio    float64 `json:"byte_ratio"`
	TokenRatio   float64 `json:"token_ratio"`
	TokenSaving  float64 `json:"token_saving_pct"`
	Estimator    string  `json:"estimator"`
	Note         string  `json:"note"`
}

// EstimateTokens is a BPE-ish stand-in: whitespace-separated runs plus each
// non-ASCII glyph counts as one token. Replace with tiktoken-go for mid-term gates.
func EstimateTokens(src []byte) int {
	if len(src) == 0 {
		return 0
	}
	n := 0
	i := 0
	for i < len(src) {
		r, size := utf8.DecodeRune(src[i:])
		if r == utf8.RuneError && size == 1 {
			n++
			i++
			continue
		}
		if unicode.IsSpace(r) {
			i += size
			continue
		}
		if r > unicode.MaxASCII {
			n++ // glyph / unicode keyword stand-in
			i += size
			continue
		}
		// ASCII word/punct run as one rough token chunk for words; split punct.
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			n++
			i += size
			for i < len(src) {
				r2, s2 := utf8.DecodeRune(src[i:])
				if unicode.IsLetter(r2) || unicode.IsDigit(r2) || r2 == '_' {
					i += s2
					continue
				}
				break
			}
			continue
		}
		n++ // punctuation
		i += size
	}
	return n
}

// Compare estimates Go vs VeGo token load.
func Compare(goSrc, vegoSrc []byte) Report {
	gt := EstimateTokens(goSrc)
	vt := EstimateTokens(vegoSrc)
	br, tr, save := 0.0, 0.0, 0.0
	if len(goSrc) > 0 {
		br = float64(len(vegoSrc)) / float64(len(goSrc))
	}
	if gt > 0 {
		tr = float64(vt) / float64(gt)
		save = (1 - tr) * 100
	}
	return Report{
		GoBytes:     len(goSrc),
		VeGoBytes:   len(vegoSrc),
		GoTokens:    gt,
		VeGoTokens:  vt,
		ByteRatio:   br,
		TokenRatio:  tr,
		TokenSaving: save,
		Estimator:   "vego-rough-v1",
		Note:        "Mid-term heuristic only; Alpha does not gate on token_saving_pct. Swap for tiktoken-go later.",
	}
}
