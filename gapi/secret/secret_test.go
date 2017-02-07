package secret

import "testing"

const tempSecretFile = "temp.secret"

func TestCervice_ReadFile(t *testing.T) {
	ser := new(Service)
	if err := ser.ReadFile(tempSecretFile); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", ser)
	if ser.Type != "type" {
		t.Error("parse error")
	}
}
