# Desafio Rate Limiter

### Funcionamento

1. O arquivo `.env` possui todas as configurações necessárias para rodar o projeto.
   - Por ele é possível configurar o número máximo de requisições (`RATE_LIMIT_IP`) por IP em um determinado período (`RATE_LIMIT_WINDOW_IP`) e ainda configurar quantos segundos esse IP ficará bloqueado (`RATE_LIMIT_BLOCK_WINDOW_IP`).
   - Além disso, é possível configurar o número máximo de requisições por token baseado na variável (`TOKENS_CONFIG_LIMIT`) que possui um array de objetos com as configurações de cada token.
   - A aplicação também possui um Rate Limit global, que pode ser configurado através da variável `RATE_LIMIT_GLOBAL`, que limita a quantidade de requests recebidos independente da origem.
   - É possível configurar também dr o rate limit por IP ou por token está ativado através das variáveis `ALLOW_IP_LIMIT` e `ALLOW_TOKEN_LIMIT`, respectivamente.
   - A porta da aplicação pode ser configurada através da variável `WEB_SERVER_PORT`. Atualmente é `8080`.
   - A variável `REDIS_HOST` é a URL do Redis, que deve ser configurada para que a aplicação funcione corretamente.
   - A variável `REDIS_PORT` é a porta do Redis, que deve ser configurada para que a aplicação funcione corretamente.

2. O arquivo `docker-compose.yml` possui as configurações necessárias para rodar o projeto.
    - Para rodar o projeto, basta executar o comando `docker-compose up` na raiz do projeto.
    - Para testar as requisições por IP, basta executar o comando `curl localhost:8080`, por exemplo.
    - Para testar as requisições por token, basta executar o comando ` curl -H "API_KEY: 1c5986a9-44d2-4369-ad21-40cafccbb56c" http://localhost:8080`, por exemplo.