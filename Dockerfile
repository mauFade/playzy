FROM golang:alpine

WORKDIR /usr/src/app


RUN go install github.com/air-verse/air@latest

# Copiar arquivos do projeto
COPY . .

# Baixar dependências do Go
RUN go mod tidy

# Expor porta para o contêiner
EXPOSE 8080

CMD ["air"]