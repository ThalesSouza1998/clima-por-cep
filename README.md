# Clima por CEP

API em Go que recebe um CEP, descobre a cidade no [ViaCEP](https://viacep.com.br) e busca a temperatura atual na [WeatherAPI](https://www.weatherapi.com). A resposta vem em Celsius, Fahrenheit e Kelvin.

Fiz o projeto para o desafio do curso Go Expert, mas ele serve como base para qualquer consulta que precise cruzar CEP com clima.

## Uso

```bash
curl "http://localhost:8080/?cep=01001000"
```

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

| Situação | Status | Corpo |
| --- | --- | --- |
| CEP fora do formato de 8 dígitos | 422 | `{"error":"invalid zipcode"}` |
| CEP válido que o ViaCEP não conhece | 404 | `{"error":"can not find zipcode"}` |

As mensagens de erro continuam em inglês porque o contrato do desafio é assim. O Fahrenheit sai com uma casa decimal e o Kelvin com duas.

## Rodando localmente

Você precisa de Go 1.24+ e de uma chave da WeatherAPI (o plano gratuito atende).

```bash
cp .env.example .env
# coloque sua chave no .env
go run ./cmd/servidor
```

O servidor lê `WEATHER_API_KEY` e `PORT` (8080 por padrão). O `.env` é carregado automaticamente quando existe; em produção, use as variáveis de ambiente do próprio Cloud Run.

Também dá para rodar sem o arquivo:

```bash
WEATHER_API_KEY=sua_chave go run ./cmd/servidor
```

## Docker

```bash
docker build -t clima-por-cep .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua_chave clima-por-cep
```

## Testes

```bash
go test ./...
```

Os testes cobrem as conversões de temperatura, as regras do serviço (CEP inválido, CEP inexistente e falha dos fornecedores) e as respostas HTTP.

## Deploy no Cloud Run

O serviço está publicado no projeto `clima-por-cep-512804`, região `us-central1`:

https://clima-por-cep-429122214318.us-central1.run.app

### APIs e repositório de imagens

```bash
gcloud config set project clima-por-cep-512804

gcloud services enable run.googleapis.com cloudbuild.googleapis.com \
  artifactregistry.googleapis.com secretmanager.googleapis.com

gcloud artifacts repositories create clima-por-cep \
  --repository-format=docker \
  --location=us-central1 \
  --description="Imagens do Clima por CEP"
```

### Chave da WeatherAPI

A chave não fica no código nem no comando de deploy. Ela é gravada no Secret Manager:

```bash
gcloud secrets create weather-api-key --replication-policy=automatic

printf '%s' "$WEATHER_API_KEY" | \
  gcloud secrets versions add weather-api-key --data-file=-
```

Quando a chave for rotacionada, basta adicionar uma nova versão ao mesmo segredo.

A conta de serviço usada pelo Cloud Run precisa poder ler o segredo:

```bash
gcloud secrets add-iam-policy-binding weather-api-key \
  --member="serviceAccount:429122214318-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### Build e deploy

```bash
gcloud builds submit \
  --tag us-central1-docker.pkg.dev/clima-por-cep-512804/clima-por-cep/clima-por-cep:latest

gcloud run deploy clima-por-cep \
  --image us-central1-docker.pkg.dev/clima-por-cep-512804/clima-por-cep/clima-por-cep:latest \
  --region us-central1 \
  --allow-unauthenticated \
  --set-secrets WEATHER_API_KEY=weather-api-key:latest \
  --min-instances 0 \
  --max-instances 2
```

O Cloud Run injeta a variável `PORT`, então a aplicação usa a porta recebida pelo ambiente. Com `min-instances 0`, o serviço não fica consumindo recursos quando está sem tráfego.

## Organização

- `cmd/servidor`: entrada da aplicação e configuração do servidor HTTP
- `internal/clima`: regras de negócio e conversões de temperatura
- `internal/viacep` e `internal/weatherapi`: clientes das APIs externas
- `internal/api`: handler HTTP e tradução dos erros para os status da API
