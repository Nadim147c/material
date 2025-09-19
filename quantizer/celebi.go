package quantizer

import "context"

func QuantizeCelebi(input pixels, maxColor int) QuantizedMap {
	qm, _ := QuantizeCelebiWithContext(context.Background(), input, maxColor)
	return qm
}

func QuantizeCelebiWithContext(ctx context.Context, input pixels, maxColor int) (QuantizedMap, error) {
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
