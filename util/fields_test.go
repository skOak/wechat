package util

import (
	"testing"
)

func TestStructFieldsMap(t *testing.T) {
	s := struct {
		AxxBxx string
		Cxx    []string          `xml:"cXX,omitempty"`
		Dxx    map[string]string `xml:"xyz"`
		Exx    string
		Fxx    string `xml:"fuck, omitempty"`
	}{
		"value_axx_bxx",
		[]string{
			"1_value_cxx",
			"2_value_cxx",
		},
		nil,
		"",
		"",
	}

	t.Log(StructFieldsMap(s, "xml"))
}
