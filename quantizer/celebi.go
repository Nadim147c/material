package quantizer

import "context"

// QuantizeCelebi is an image quantizer that improves on the quality of a
// standard K-Means algorithm by setting the K-Means initial state to the output
// of a Wu quantizer, instead of random centroids. Improves on speed by several
// optimizations, as implemented in Wsmeans, or Weighted Square Means, K-Means
// with those optimizations.
//
// This algorithm was designed by M. Emre Celebi, and was found in their 2011
// paper, Improving the Performance of K-Means for Color Quantization.
// https://arxiv.org/abs/1101.0395
func QuantizeCelebi(input pixels, maxColor int) QuantizedMap {
	qm, _ := QuantizeCelebiWithContext(context.Background(), input, maxColor)
	return qm
}

// QuantizeCelebiWithContext is QuantizeCelebi with context.Context support.
//
// Deprecated: Use QuantizeCelebiContext
func QuantizeCelebiWithContext(ctx context.Context, input pixels, maxColor int) (QuantizedMap, error) {
	return QuantizeCelebiContext(ctx, input, maxColor)
}

// QuantizeCelebiContext is QuantizeCelebi with context.Context support. Returns
// ctx.Err() if context is Done.
func QuantizeCelebiContext(ctx context.Context, input pixels, maxColor int) (QuantizedMap, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	wu, err := QuantizeWuWithContext(ctx, input, maxColor*5)
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

	qm, err := QuantizeWsMeansWithContext(ctx, input, colors, maxColor)
	if err != nil {
		return nil, err
	}

	return qm, nil
}
