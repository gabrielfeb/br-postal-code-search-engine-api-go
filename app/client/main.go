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
	Fonte      string `json:"fonte"`
}

func main() {
	cep := "01153000"
	url := fmt.Sprintf("http://localhost:8080/cep?cep=%s", cep)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Erro ao criar request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status HTTP não OK: %s", resp.Status)
	}

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		log.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	fmt.Printf("Resposta da API %s:\n", endereco.Fonte)
	fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n",
		endereco.Cep, endereco.Logradouro, endereco.Bairro, endereco.Localidade, endereco.Uf)
}
