# SPDX-FileCopyrightText: © 2022 ChiselStrike <info@chiselstrike.com>

FROM golang:1.19 as builder

COPY . /app/
WORKDIR /app
RUN go build -o main

FROM gcr.io/distroless/base
COPY --from=builder /app/main/ replayer

EXPOSE 8080
CMD ["./replayer"]
