// Package nocopy adds contains a functionality which disallows the copying of custom types.
// This is possible by 'abusing' the go vet tool.
//
// Struct may be added to structs which must not be copied.
//
// Lock and Unlock are no-op used by -copylocks checker from `go vet`.
//
// # Should be used like this:
//
//	type MySensetiveStruct struct {
//		_ nocopy.Struct
//
//		// sensitive to copying fields...
//	}
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
package nocopy

type Struct struct{}

func (*Struct) Lock()   {}
func (*Struct) Unlock() {}
