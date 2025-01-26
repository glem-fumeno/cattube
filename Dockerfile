FROM golang:1.23

WORKDIR /app

RUN mkdir -p /destination

RUN mkdir -p /out

RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go build -o /out/cattube .

EXPOSE 8090

CMD ["/out/cattube"]

