package object

import (
	"ego/internal/backend/eval/object"
	"testing"
)

func TestStringMapKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}
	diff1 := &object.String{Value: "My name is johnny"}
	diff2 := &object.String{Value: "My name is johnny"}

	if hello1.HasKey() != hello2.HasKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if diff1.HasKey() != diff2.HasKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if hello1.HasKey() == diff1.HasKey() {
		t.Errorf("strings with different content have the same hash keys")
	}

}
