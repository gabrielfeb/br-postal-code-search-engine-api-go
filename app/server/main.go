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

func buscarBrasilAPI(ctx context.Context, cep string, ch chan<- Endereco) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Erro BrasilAPI: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status BrasilAPI: %s", resp.Status)
		return
	}

	var r struct {
		CEP          string `json:"cep"`
		Street       string `json:"street"`
		Neighborhood string `json:"neighborhood"`
		City         string `json:"city"`
		State        string `json:"state"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("Erro decode BrasilAPI: %v", err)
		return
	}

	select {
	case ch <- Endereco{
		Fonte:      "BrasilAPI",
		Cep:        r.CEP,
		Logradouro: r.Street,
		Bairro:     r.Neighborhood,
		Localidade: r.City,
		Uf:         r.State,
	}:
	case <-ctx.Done():
		log.Println("Context cancelado antes de enviar BrasilAPI")
	}
}

func buscarViaCEP(ctx context.Context, cep string, ch chan<- Endereco) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Erro ViaCEP: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status ViaCEP: %s", resp.Status)
		return
	}

	var r struct {
		CEP        string `json:"cep"`
		Logradouro string `json:"logradouro"`
		Bairro     string `json:"bairro"`
		Localidade string `json:"localidade"`
		Uf         string `json:"uf"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf("Erro decode ViaCEP: %v", err)
		return
	}

	select {
	case ch <- Endereco{
		Fonte:      "ViaCEP",
		Cep:        r.CEP,
		Logradouro: r.Logradouro,
		Bairro:     r.Bairro,
		Localidade: r.Localidade,
		Uf:         r.Uf,
	}:
	case <-ctx.Done():
		log.Println("Context cancelado antes de enviar ViaCEP")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan Endereco, 1)

	go buscarBrasilAPI(ctx, cep, ch)
	go buscarViaCEP(ctx, cep, ch)

	select {
	case res := <-ch:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	case <-ctx.Done():
		http.Error(w, "Timeout: nenhuma API respondeu em 1 segundo", http.StatusGatewayTimeout)
	}
}

func main() {
	http.HandleFunc("/cep", handler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
