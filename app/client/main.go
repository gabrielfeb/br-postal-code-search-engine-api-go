package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cep?cep=01153000", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Erro da API: status %d", resp.StatusCode)
	}

	var result struct {
		Origem   string                 `json:"origem"`
		Endereco map[string]interface{} `json:"endereco"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	fmt.Printf("API usada: %s\n", result.Origem)
	fmt.Println("Dados do endereço:")
	for k, v := range result.Endereco {
		fmt.Printf("  %s: %v\n", k, v)
	}

	// opcional: salvar em arquivo
	file, _ := os.Create("endereco.txt")
	defer file.Close()
	json.NewEncoder(file).Encode(result)
}
