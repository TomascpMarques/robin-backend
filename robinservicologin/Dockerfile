FROM golang:alpine

# normaliza a dir corrente para a dir /app
# e copia os conteúdos da dir onde o dockerfile está para /app
WORKDIR /app
COPY . .

# instala as dependências do projeto
RUN go mod download

# Defenição de valores para variaveis env
ENV AUTH_SERVER_REDIS_PORT="6379"
ENV REDIS_USER1_PASS=""
ENV REDIS_USER1_NAME=""
ENV LOGIN_SERV_PORT=8080

# compila o programa para um executável
RUN go build

# expõe a porta 8080 do container 
EXPOSE 8080

# Corre o executável compilado do código fonte
CMD [ "./robinservicologin" ]