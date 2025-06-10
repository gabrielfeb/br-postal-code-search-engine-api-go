package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ValidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep?cep=01153000", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado status 200, obteve %d", resp.StatusCode)
	}

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if endereco.Cep == "" || endereco.Fonte == "" {
		t.Errorf("dados incompletos: %+v", endereco)
	}
}

func TestHandler_InvalidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep?cep=INVALID", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("esperado erro com cep inválido, mas obteve 200")
	}
}

func TestHandler_EmptyCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cep", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400 para CEP vazio, obteve %d", resp.StatusCode)
	}
}

func FuzzHandlerCEP(f *testing.F) {

	seeds := []string{
		"01153000",
		"01001000",
		"12345678",
		"99999999",
		"00000000",
		"abcdefgh",
		"1234",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, cep string) {

		if len(cep) != 8 {
			t.Skipf("ignorando CEP com tamanho != 8: %q", cep)
		}
		for _, c := range cep {
			if c < '0' || c > '9' {
				t.Skipf("ignorando CEP com caractere não numérico: %q", cep)
			}
		}

		req := httptest.NewRequest(http.MethodGet, "/cep?cep="+cep, nil)
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("status inesperado %d para CEP %q", resp.StatusCode, cep)
		}
	})
}
