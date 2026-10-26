FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /servidor ./cmd/servidor

FROM alpine:3.22

# O binário conversa com ViaCEP e WeatherAPI por HTTPS, então os certificados
# precisam estar na imagem.
RUN apk add --no-cache ca-certificates

COPY --from=build /servidor /servidor

USER 65534:65534
EXPOSE 8080

ENTRYPOINT ["/servidor"]
