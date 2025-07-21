# マルチステージビルド: Goアプリケーションのビルド
FROM golang:1.23-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# 統合実行環境: PostgreSQL + Goアプリ
FROM debian:bullseye-slim

# 必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    postgresql-13 \
    postgresql-client-13 \
    postgresql-contrib-13 \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

# PostgreSQL設定
USER postgres
RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER webapp WITH SUPERUSER PASSWORD 'password';" && \
    createdb -O webapp webapp

USER root

# Goアプリケーションのコピー
COPY --from=builder /app/main /usr/local/bin/app

# Supervisord設定
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# PostgreSQL設定ファイルの調整
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/13/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/13/main/postgresql.conf

# ポート公開
EXPOSE 8080 5432

# Supervisordで両プロセス起動
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]