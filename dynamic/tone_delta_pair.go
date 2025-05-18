package dynamic

type TonePolarity int

const (
	Darker TonePolarity = iota
	Lighter
	Nearer
	Farther
)

func (p TonePolarity) String() string {
	switch p {
	case Darker:
		return "Darker"
	case Lighter:
		return "Lighter"
	case Nearer:
		return "Nearer"
	case Farther:
		return "Farther"
	default:
		return "Unknown"
	}
}

// ToneDeltaPair defines a tone difference constraint between two DynamicColors.
type ToneDeltaPair struct {
	RoleA        DynamicColor
	RoleB        DynamicColor
	Delta        float64
	Polarity     TonePolarity
	StayTogether bool
}

// NewToneDeltaPair constructs a ToneDeltaPair with the given parameters.
func NewToneDeltaPair(roleA, roleB DynamicColor, delta float64, polarity TonePolarity, stayTogether bool) ToneDeltaPair {
	return ToneDeltaPair{
		RoleA:        roleA,
		RoleB:        roleB,
		Delta:        delta,
		Polarity:     polarity,
		StayTogether: stayTogether,
	}
}
