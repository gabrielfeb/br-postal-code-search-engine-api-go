package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Endereco struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Origem     string `json:"-"`
}

func main() {
	http.HandleFunc("/cep", func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		if cep == "" {
			http.Error(w, "CEP não informado", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		resChan := make(chan Endereco, 2)

		go fetchFromBrasilAPI(ctx, cep, resChan)
		go fetchFromViaCEP(ctx, cep, resChan)

		select {
		case res := <-resChan:
			resJson := struct {
				Origem   string   `json:"origem"`
				Endereco Endereco `json:"endereco"`
			}{Origem: res.Origem, Endereco: res}

			json.NewEncoder(w).Encode(resJson)
		case <-ctx.Done():
			http.Error(w, "Timeout ao buscar CEP", http.StatusGatewayTimeout)
		}
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchFromBrasilAPI(ctx context.Context, cep string, ch chan<- Endereco) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		return
	}
	endereco.Origem = "BrasilAPI"
	ch <- endereco
}

func fetchFromViaCEP(ctx context.Context, cep string, ch chan<- Endereco) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		return
	}
	endereco.Origem = "ViaCEP"
	ch <- endereco
}
