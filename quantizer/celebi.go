package quantizer

import (
	"context"

	"github.com/Nadim147c/material/color"
)

// QuantizeCelebi is an image quantizer that improves on the quality of a
// standard K-Means algorithm by setting the K-Means initial state to the output
// of a Wu quantizer, instead of random centroids. Improves on speed by several
// optimizations, as implemented in Wsmeans, or Weighted Square Means, K-Means
// with those optimizations.
//
// This algorithm was designed by M. Emre Celebi, and was found in their 2011
// paper, Improving the Performance of K-Means for Color Quantization.
// https://arxiv.org/abs/1101.0395
func QuantizeCelebi(input []color.ARGB, maxColor int) QuantizedMap {
	qm, _ := QuantizeCelebiContext(context.Background(), input, maxColor)
	return qm
}

// QuantizeCelebiContext is QuantizeCelebi with context.Context support. Returns
// ctx.Err() if context is Done.
func QuantizeCelebiContext(
	ctx context.Context,
	input []color.ARGB,
	maxColor int,
) (QuantizedMap, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	wu, err := QuantizeWuContext(ctx, input, maxColor*5)
	if err != nil {
		return nil, err
	}

	colors := make(pixelsLab, len(wu))
	for i, c := range wu {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		colors[i] = c.ToLab()
	}

	qm, err := QuantizeWsMeansContext(ctx, input, colors, maxColor)
	if err != nil {
		return nil, err
	}

	return qm, nil
}
