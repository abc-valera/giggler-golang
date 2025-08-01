// Package enum provides a helper interface for enums validation.
//
// # All Enums should be defined in the following ways
//
// # Enums without a zero value:
//
//	type Mood int
//
//	const (
//		// iota starts from 1, because enum does not has a zero value
//		Happy Mood = iota + 1
//		Sad
//		Angry
//
//		// end should always be the last field of the enum, signaling its end.
//		end
//	)
//
//	func (s Mood) IsValid() bool {
//		return s > 0 && s < end
//	}
//
//	// This is to ensure that the Status type implements the Interface interface
//	var _ enum.Interface = (*Mood)(nil)
//
// # Enums with a zero value:
//
//	type Mood int
//
//	const (
//		// DefaultMood is the zero value of the enum.
//		DefaultMood Mood = iota
//		Happy
//		Sad
//		Angry
//
//		// end should always be the last field of the enum, signaling its end.
//		end
//	)
//
//	func (s Mood) IsValid() bool {
//		return s >= DefaultMood && s < end
//	}
//
//	// This is to ensure that the Status type implements the Interface interface
//	var _ enum.Interface = (*Mood)(nil)
package enum

// Interface is the interface that all enums should implement.
type Interface interface {
	IsValid() bool
}
