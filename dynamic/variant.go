package dynamic

type Variant int

const (
	Monochrome Variant = iota
	Neutral
	TonalSpot
	Vibrant
	Expressive
	Fidelity
	Content
	Rainbow
	FruitSalad
)

func (v Variant) String() string {
	switch v {
	case Monochrome:
		return "Monochrome"
	case Neutral:
		return "Neutral"
	case TonalSpot:
		return "TonalSpot"
	case Vibrant:
		return "Vibrant"
	case Expressive:
		return "Expressive"
	case Fidelity:
		return "Fidelity"
	case Content:
		return "Content"
	case Rainbow:
		return "Rainbow"
	case FruitSalad:
		return "FruitSalad"
	default:
		return "Unknown"
	}
}
