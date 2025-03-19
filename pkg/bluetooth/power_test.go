package bluetooth

import "testing"

func TestDecode(t *testing.T) {
	bytes := []byte{16, 0, 161, 0, 3, 0, 0, 0, 66, 8}
	power := decode(bytes)
	if power != 161 {
		t.Error("Power should be equal to 161")
	}
}

func TestEncode(t *testing.T) {
	power := 300
	expected := []byte{5, 44, 1}
	actual := encode(power)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("expected %d, got %d\n", expected[i], actual[i])
		}
	}
}
