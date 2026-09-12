// Package bpe provides token-count helpers for VeGo compression goals (v0.1β).
// Beta measures with tiktoken (cl100k_base). Alpha did not gate on these metrics.
package bpe

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// DefaultEncoding is the Beta primary tokenizer.
const DefaultEncoding = "cl100k_base"

// Report is a JSON-serializable comparison of Go vs .vego token counts.
type Report struct {
	GoBytes     int     `json:"go_bytes"`
	VeGoBytes   int     `json:"vego_bytes"`
	GoTokens    int     `json:"go_tokens"`
	VeGoTokens  int     `json:"vego_tokens"`
	ByteRatio   float64 `json:"byte_ratio"`
	TokenRatio  float64 `json:"token_ratio"`
	TokenSaving float64 `json:"token_saving_pct"`
	ByteSaving  float64 `json:"byte_saving_pct"`
	Estimator   string  `json:"estimator"`
	Note        string  `json:"note"`
}

// Count returns tiktoken token count for src using encoding name.
func Count(src []byte, encoding string) (int, error) {
	if encoding == "" {
		encoding = DefaultEncoding
	}
	enc, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, fmt.Errorf("bpe: encoding %q: %w", encoding, err)
	}
	return len(enc.Encode(string(src), nil, nil)), nil
}

// Compare estimates Go vs VeGo load with tiktoken cl100k_base.
func Compare(goSrc, vegoSrc []byte) (Report, error) {
	gt, err := Count(goSrc, DefaultEncoding)
	if err != nil {
		return Report{}, err
	}
	vt, err := Count(vegoSrc, DefaultEncoding)
	if err != nil {
		return Report{}, err
	}
	br, tr, tSave, bSave := 0.0, 0.0, 0.0, 0.0
	if len(goSrc) > 0 {
		br = float64(len(vegoSrc)) / float64(len(goSrc))
		bSave = (1 - br) * 100
	}
	if gt > 0 {
		tr = float64(vt) / float64(gt)
		tSave = (1 - tr) * 100
	}
	return Report{
		GoBytes:     len(goSrc),
		VeGoBytes:   len(vegoSrc),
		GoTokens:    gt,
		VeGoTokens:  vt,
		ByteRatio:   br,
		TokenRatio:  tr,
		TokenSaving: tSave,
		ByteSaving:  bSave,
		Estimator:   "tiktoken/" + DefaultEncoding,
		Note:        "v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.",
	}, nil
}
