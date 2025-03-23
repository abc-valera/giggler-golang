package nocopy

// Struct may be added to structs which must not be copied
// after the first use.
//
// Lock and Unlock are no-op used by -copylocks checker from `go vet`.
//
// # Should be used like this:
//
//	type myStruct struct {
//		_ nocopy.Struct
//
//		// sensitive to copying fields...
//	}
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type Struct struct{}

func (*Struct) Lock()   {}
func (*Struct) Unlock() {}
