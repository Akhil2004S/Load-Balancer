FROM golang:1.26.4

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/loadbalancer ./

ENTRYPOINT [ "/build/loadbalancer" ]

