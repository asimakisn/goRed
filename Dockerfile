FROM debian:bookworm-slim  AS base

EXPOSE 6379

RUN apt-get update && apt-get install -y \
    curl \
    vim \
    git \
    build-essential

RUN curl -fsSL https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -o go.tar.gz
RUN rm -rf /usr/local/go
RUN tar -C /usr/local -xzf go.tar.gz
RUN rm go.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /app

COPY . /app

RUN go build .

WORKDIR /app

CMD ["/app/main"]
