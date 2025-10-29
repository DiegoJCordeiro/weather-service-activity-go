# Weather Service - Serviço de Clima por CEP

Sistema em Go que recebe um CEP brasileiro, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## 📋 Índice

- URL GCP: https://weather-service-activity-go-926876726731.us-central1.run.app/

- [Funcionalidades](#-funcionalidades)
- [Requisitos](#-requisitos)
- [Instalação e Configuração](#-instalação-e-configuração)
- [Uso Local](#-uso-local)
- [Endpoints da API](#-endpoints-da-api)
- [Testes](#-testes)
- [Deploy no Google Cloud Run](#-deploy-no-google-cloud-run)
- [Exemplos de Uso](#-exemplos-de-uso)
- [Troubleshooting](#-troubleshooting)

## 🚀 Funcionalidades

- ✅ Validação de CEP (8 dígitos)
- ✅ Consulta de localização via API ViaCEP
- ✅ Consulta de temperatura via WeatherAPI
- ✅ Conversão automática para Fahrenheit e Kelvin
- ✅ Tratamento de erros apropriado (422, 404, 500)
- ✅ Health check endpoint
- ✅ Suporte a Docker e Docker Compose
- ✅ Pronto para deploy no Google Cloud Run

## 📦 Requisitos

- Docker e Docker Compose (para execução containerizada)
- Go 1.21+ (para desenvolvimento local)
- Conta WeatherAPI (gratuita) - [Criar conta](https://www.weatherapi.com/signup.aspx)
- Google Cloud SDK (apenas para deploy)

## 🔧 Instalação e Configuração

### 1. Estrutura do Projeto

```
weather-service/
├── main.go              # Código principal
├── main_test.go         # Testes automatizados
├── go.mod               # Dependências Go
├── go.sum               # Checksums
├── Dockerfile           # Configuração Docker
├── docker-compose.yml   # Orquestração
├── .env                 # Variáveis de ambiente
├── .gitignore          # Arquivos ignorados
└── README.md           # Documentação
```

### 2. Obter API Key do WeatherAPI

1. Acesse [https://www.weatherapi.com/](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Copie sua API key do dashboard

### 3. Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```bash
WEATHER_API_KEY=sua_chave_aqui
PORT=8080
```

## 🐳 Uso Local com Docker

### Iniciar o Serviço

```bash
# Build e iniciar
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar o serviço
docker-compose down
```

### Testar o Serviço

```bash
# CEP válido (Av. Paulista, São Paulo)
curl http://localhost:8080/weather/01310100

# Resposta esperada:
# {"temp_C":25.0,"temp_F":77.0,"temp_K":298.0}

# CEP inválido (formato incorreto)
curl http://localhost:8080/weather/123
# {"message":"invalid zipcode"}

# CEP não encontrado
curl http://localhost:8080/weather/99999999
# {"message":"can not find zipcode"}

# Health check
curl http://localhost:8080/health
# OK
```

## 🏃 Uso Local sem Docker

```bash
# Baixar dependências
go mod download

# Configurar variável de ambiente
export WEATHER_API_KEY=sua_chave_aqui

# Executar
go run main.go

# Em outro terminal, testar
curl http://localhost:8080/weather/01310100
```

## 📡 Endpoints da API

### GET /weather/{cep}

Retorna a temperatura atual para o CEP informado.

**Parâmetros:**
- `cep`: CEP brasileiro de 8 dígitos (com ou sem hífen)

**Respostas:**

#### ✅ 200 OK - Sucesso
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### ❌ 422 Unprocessable Entity - CEP Inválido
```json
{
  "message": "invalid zipcode"
}
```

**Quando ocorre:**
- CEP com menos ou mais de 8 dígitos
- CEP contendo letras ou caracteres especiais

#### ❌ 404 Not Found - CEP Não Encontrado
```json
{
  "message": "can not find zipcode"
}
```

**Quando ocorre:**
- CEP não existe na base da ViaCEP

### GET /health

Health check endpoint para monitoramento.

**Resposta:**
```
HTTP/1.1 200 OK
OK
```

## 🧪 Testes

### Executar Testes Unitários

```bash
# Todos os testes
go test -v

# Com cobertura
go test -v -cover

# Relatório de cobertura HTML
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Script de Teste Automatizado

Crie o arquivo `test.sh`:

```bash
#!/bin/bash

BASE_URL="${1:-http://localhost:8080}"

echo "Testing valid CEP..."
curl -s $BASE_URL/weather/01310100 | jq

echo -e "\nTesting invalid CEP format..."
curl -s $BASE_URL/weather/123 | jq

echo -e "\nTesting CEP not found..."
curl -s $BASE_URL/weather/99999999 | jq

echo -e "\nTesting health check..."
curl -s $BASE_URL/health
```

Executar:
```bash
chmod +x test.sh
./test.sh
```

## ☁️ Deploy no Google Cloud Run

### Pré-requisitos

1. **Instalar Google Cloud SDK**
   ```bash
   # Linux/macOS
   curl https://sdk.cloud.google.com | bash
   exec -l $SHELL
   gcloud init
   ```

2. **Configurar Projeto**
   ```bash
   gcloud auth login
   gcloud config set project SEU_PROJECT_ID
   ```

3. **Habilitar APIs**
   ```bash
   gcloud services enable run.googleapis.com
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable containerregistry.googleapis.com
   ```

### Deploy em 2 Comandos

```bash
# 1. Build da imagem
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/weather-service

# 2. Deploy no Cloud Run
gcloud run deploy weather-service \
  --image gcr.io/SEU_PROJECT_ID/weather-service \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SUA_CHAVE_AQUI \
  --memory 256Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10
```

### Obter URL do Serviço

```bash
gcloud run services describe weather-service \
  --region us-central1 \
  --format 'value(status.url)'
```

### Testar Serviço no Cloud Run

```bash
# Substitua pela URL retornada
curl https://weather-service-xxxxx-xx.a.run.app/weather/01310100
```

### Atualizar o Serviço

```bash
# Rebuild e redeploy
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/weather-service
gcloud run deploy weather-service \
  --image gcr.io/SEU_PROJECT_ID/weather-service \
  --region us-central1
```

### Gerenciar com Secrets Manager (Recomendado)

```bash
# Criar secret
echo -n "SUA_WEATHER_API_KEY" | \
  gcloud secrets create weather-api-key \
  --data-file=- \
  --replication-policy="automatic"

# Deploy usando secret
gcloud run deploy weather-service \
  --image gcr.io/SEU_PROJECT_ID/weather-service \
  --update-secrets=WEATHER_API_KEY=weather-api-key:latest \
  --region us-central1
```

### Monitoramento

```bash
# Ver logs em tempo real
gcloud run services logs tail weather-service --region us-central1

# Ver métricas no console
# https://console.cloud.google.com/run
```

### Deletar o Serviço

```bash
gcloud run services delete weather-service --region us-central1
```

## 📝 Exemplos de Uso

### CEPs para Teste

| CEP | Localização |
|-----|-------------|
| `01310100` | Av. Paulista, São Paulo/SP |
| `20040020` | Centro, Rio de Janeiro/RJ |
| `30130100` | Centro, Belo Horizonte/MG |
| `70040902` | Brasília/DF |
| `88015100` | Centro, Florianópolis/SC |

### cURL

```bash
# Básico
curl http://localhost:8080/weather/01310100

# Com formatação JSON (requer jq)
curl -s http://localhost:8080/weather/01310100 | jq

# Ver headers completos
curl -v http://localhost:8080/weather/01310100
```

### JavaScript (Fetch API)

```javascript
async function getWeather(cep) {
  try {
    const response = await fetch(`http://localhost:8080/weather/${cep}`);
    const data = await response.json();
    
    if (response.ok) {
      console.log(`Temperatura: ${data.temp_C}°C`);
    } else {
      console.error(`Erro: ${data.message}`);
    }
  } catch (error) {
    console.error('Erro de conexão:', error);
  }
}

getWeather('01310100');
```

### Python (requests)

```python
import requests

def get_weather(cep):
    try:
        response = requests.get(f'http://localhost:8080/weather/{cep}')
        response.raise_for_status()
        data = response.json()
        print(f"Temperatura: {data['temp_C']}°C")
        return data
    except requests.exceptions.HTTPError as e:
        print(f"Erro HTTP: {e}")
        print(f"Mensagem: {response.json()['message']}")
    except requests.exceptions.RequestException as e:
        print(f"Erro de conexão: {e}")

weather = get_weather('01310100')
```

## 🧮 Fórmulas de Conversão

O sistema utiliza as seguintes fórmulas para conversão de temperatura:

- **Celsius para Fahrenheit**: `F = C × 1.8 + 32`
- **Celsius para Kelvin**: `K = C + 273`

## 🛠️ Desenvolvimento

### Comandos Úteis (Makefile)

Crie um `Makefile` com os comandos:

```makefile
.PHONY: run test docker-build docker-run docker-stop

run:
	go run main.go

test:
	go test -v -cover

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

test-api:
	./test.sh
```

Uso:
```bash
make run          # Executar localmente
make test         # Rodar testes
make docker-run   # Iniciar com Docker
make test-api     # Testar API
```

### Estrutura do Código

**main.go** - Funções principais:
- `handleWeather()` - Handler principal do endpoint /weather
- `isValidCEPFormat()` - Valida formato do CEP
- `getLocationByCEP()` - Consulta ViaCEP
- `getTemperature()` - Consulta WeatherAPI
- `celsiusToFahrenheit()` - Conversão C → F
- `celsiusToKelvin()` - Conversão C → K

**main_test.go** - Testes:
- Validação de formato de CEP
- Conversões de temperatura
- Testes de endpoints HTTP
- Casos de sucesso e erro

## 🔍 Troubleshooting

### Porta 8080 já em uso

```bash
# Mudar porta no .env
echo "PORT=8081" > .env

# Reiniciar
docker-compose down && docker-compose up -d
```

### Erro de API Key

```bash
# Verificar se a variável está configurada
docker-compose exec weather-service env | grep WEATHER

# Verificar se a chave é válida
curl "http://api.weatherapi.com/v1/current.json?key=SUA_CHAVE&q=London"
```

### Build Docker falha

```bash
# Limpar cache e rebuildar
docker-compose down
docker system prune -f
docker-compose build --no-cache
docker-compose up -d
```

### CEP retorna temperatura zerada

Verifique:
1. API Key configurada corretamente
2. Cota da WeatherAPI não esgotada (plano free: 1M req/mês)
3. Logs do serviço: `docker-compose logs`

### Erro 500 Internal Server Error

```bash
# Ver logs detalhados
docker-compose logs -f

# Verificar conectividade com APIs externas
curl https://viacep.com.br/ws/01310100/json/
curl "http://api.weatherapi.com/v1/current.json?key=SUA_CHAVE&q=SaoPaulo"
```

## 💰 Custos - Google Cloud Run

### Free Tier (por mês):
- 2 milhões de requisições
- 360,000 GB-seconds de memória
- 180,000 vCPU-seconds

### Estimativa:
- **10,000 req/mês**: Grátis
- **100,000 req/mês**: ~$0.04
- **1,000,000 req/mês**: ~$0.40

Com `min-instances=0`, o serviço escala para zero quando não há tráfego = **custo zero**.

## 🔐 Segurança

- ✅ Validação rigorosa de input (CEP)
- ✅ Não expõe API keys nos logs
- ✅ Suporte a HTTPS automático no Cloud Run
- ✅ Rate limiting configurável
- ✅ Health checks para monitoramento

**Recomendações:**
- Use Secrets Manager em produção
- Implemente rate limiting
- Configure CORS se necessário
- Monitore logs e métricas

## 📊 Performance

**Métricas esperadas:**
- Latência média: 200-500ms (inclui chamadas externas)
- Throughput: 100+ req/s
- Taxa de erro: < 1%
- Cold start: ~2s (Cloud Run)

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/NovaFuncionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/NovaFuncionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto foi desenvolvido para fins educacionais como parte de um desafio técnico.

## 🆘 Suporte

- **Documentação Cloud Run**: [https://cloud.google.com/run/docs](https://cloud.google.com/run/docs)
- **WeatherAPI Docs**: [https://www.weatherapi.com/docs/](https://www.weatherapi.com/docs/)
- **ViaCEP**: [https://viacep.com.br/](https://viacep.com.br/)
