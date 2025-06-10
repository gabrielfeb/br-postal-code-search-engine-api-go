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

func SearchCEP(cep string) (*Endereco, error) {
	url := fmt.Sprintf("http://localhost:8080/cep?cep=%s", cep)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP não OK: %s", resp.Status)
	}

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &endereco, nil
}

func main() {
	cep := "01153000"

	endereco, err := SearchCEP(cep)
	if err != nil {
		log.Fatalf("Erro ao buscar CEP: %v", err)
	}

	fmt.Printf("Resposta da API: %s\n", endereco.Fonte)
	fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n",
		endereco.Cep, endereco.Logradouro, endereco.Bairro, endereco.Localidade, endereco.Uf)
}
