FROM golang:1.26

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y zip bzip2 xz-utils && \
    apt-get purge -y && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

WORKDIR /app
COPY [".", "."]

RUN go mod tidy
RUN go build -ldflags "-s -w" -trimpath -o main .

CMD ["./main"]
