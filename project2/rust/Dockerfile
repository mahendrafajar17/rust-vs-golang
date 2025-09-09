FROM rust:1.75 AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src

RUN cargo build --release

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root

COPY --from=builder /app/target/release/project2-rust .
COPY config.yaml .

EXPOSE 8083

CMD ["./project2-rust"]