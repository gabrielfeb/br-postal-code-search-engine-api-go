package main

import (
	"testing"
)

func FuzzSearchCEP(f *testing.F) {
	seeds := []string{
		"01001000", "01153000", "99999999", "00000000", "12345678",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, cep string) {
		if len(cep) != 8 {
			t.Skip()
		}
		for _, c := range cep {
			if c < '0' || c > '9' {
				t.Skip()
			}
		}

		endereco, err := SearchCEP(cep)
		if err != nil {
			t.Logf("Erro esperado: %v", err)
			return
		}

		if endereco.Cep == "" || endereco.Fonte == "" {
			t.Errorf("Endereço incompleto para CEP %q: %+v", cep, endereco)
		}
	})
}
