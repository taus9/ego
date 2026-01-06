package object

import "testing"

func TestStringMapKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

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
