FROM golang:alpine AS builder
COPY . /work
WORKDIR /work
RUN go build -o /work/tournament cmd/tournament/main.go

FROM alpine
RUN apk add --no-cache postgresql-client
EXPOSE 3000
COPY --from=builder /work/tournament /tournament
COPY --from=builder /work/sql /sql
COPY --from=builder /work/tools/wait-for-postgres.sh /wait-for-postgres.sh
CMD ["/tournament", "--port", "3000"]
