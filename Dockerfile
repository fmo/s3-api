FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make s3_api_players_amd

CMD ["./s3ApiApp"]
