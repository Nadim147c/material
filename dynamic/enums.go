package dynamic

// Variant indicates type of scheme to generate
//
// ENUM(monochrome, neutral, tonal_spot, vibrant, expressive, fidelity, content, rainbow, fruit_salad)
type Variant uint

// Version indicates the material color specification year
//
// ENUM(2021=2021, 2025=2025)
type Version uint

// Platform indicates target platform for generating colors
//
// ENUM(phone, watch)
type Platform uint

// TonePolarity describes the difference in tone between colors.
//
// ENUM(tone_darker, tone_lighter, tone_nearer, tone_farther, tone_relative_darker, tone_relative_lighter)
type TonePolarity uint

// DeltaConstraint describes how to fulfill a tone delta pair constraint.
//
// ENUM(constraint_exact, constraint_nearer, constraint_farther)
type DeltaConstraint uint
