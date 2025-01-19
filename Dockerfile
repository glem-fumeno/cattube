FROM golang:1.23

WORKDIR /app

RUN mkdir -p /destination

RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8090

CMD ["air"]

