package handler

import (
	"encoding/base32"
	"testing"
)

func TestBase32(t *testing.T) {
	code := base32.StdEncoding.EncodeToString([]byte("this is hello"))
	t.Log(code)
}

var id = "1264567"
var idCode = "XGMEP7NFFDTB25S2ZPU5F5MYYCAKR2UPL5LYW==="

func TestIdEncode(t *testing.T) {
	idCode, err := idEncode(id)
	if err != nil {
		t.Log(err)
	}
	t.Log("encode", idCode)
}
func TestIdDeoce(t *testing.T) {
	decodeId, err := idDecode(idCode)
	if err != nil {
		t.Log(err)
	}
	if id != decodeId {
		t.Errorf("idDecode  %s not equal %s", decodeId, id)
	}
	t.Log(id)
}
