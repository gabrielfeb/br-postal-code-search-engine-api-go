# br-postal-code-search-engine-api-go

## Descrição

Aplicação cliente-servidor em Go para buscar dados de endereços a partir do CEP (Código de Endereçamento Postal) usando APIs públicas (BrasilAPI e ViaCEP).

---

## Como executar o servidor

1. Navegue até o diretório do servidor:

```bash
cd app/server
```

2. Execute o servidor:

```bash
go run main.go
```

O servidor ficará disponível em: `http://localhost:8080`.

---

## Como executar o cliente

1. Navegue até o diretório do cliente:

```bash
cd app/client
```

2. Execute o cliente:

```bash
go run main.go
```

O cliente fará uma requisição para o servidor local na rota `/cep` com um CEP de exemplo (`01153000`) e imprimirá o resultado.

---

## Testes

### Rodar testes unitários do servidor

1. No diretório do servidor, execute:

```bash
go test -v
```

### Rodar testes com fuzzing

```bash
go test -v -fuzz=.
```

---

## Observações importantes

- O CEP deve ser sempre composto por 8 dígitos numéricos.
- O servidor faz chamadas concorrentes para as APIs ViaCEP e BrasilAPI e retorna a primeira resposta válida.
- Timeout padrão para busca é 1 segundo no servidor e 2 segundos no cliente.