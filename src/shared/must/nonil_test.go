package must

import "testing"

func TestInitializedStruct(t *testing.T) {
	subject := testSubjectStruct{
		InterfaceField: "not nil",
	}

	result := NoNilInterfaceFields(subject)

	if result != subject {
		t.Error("expected the same struct to be returned")
	}
}

func TestUnInitializedStruct(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic")
		}
	}()

	s := testSubjectStruct{}
	NoNilInterfaceFields(s)
}

type testSubjectStruct struct {
	InterfaceField any
}
