package dynamic

// ToneDeltaPair documents a constraint between two DynamicColors,
// in which their tones must have a certain distance from each other.
type ToneDeltaPair struct {
	RoleA, RoleB *DynamicColor
	Delta        float64
	Polarity     TonePolarity
	StayTogether bool
	Constraint   DeltaConstraint
}

// NewToneDeltaPair creates a new ToneDeltaPair with default constraint = "exact".
func NewToneDeltaPair(
	roleA, roleB *DynamicColor,
	delta float64,
	polarity TonePolarity,
	stayTogether bool,
	constraint ...DeltaConstraint,
) *ToneDeltaPair {
	c := ConstraintExact
	if len(constraint) > 0 {
		c = constraint[0]
	}
	return &ToneDeltaPair{
		RoleA:        roleA,
		RoleB:        roleB,
		Delta:        delta,
		Polarity:     polarity,
		StayTogether: stayTogether,
		Constraint:   c,
	}
}
