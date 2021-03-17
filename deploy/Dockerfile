FROM golang:alpine AS builder
COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /work/tournament cmd/tournament/main.go

FROM scratch
EXPOSE 3000
COPY --from=builder /work/tournament /tournament
COPY --from=builder /work/sql /sql
CMD ["/tournament", "--port", "3000"]
