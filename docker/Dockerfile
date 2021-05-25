FROM golang:1.16 as builder


ENV GO111MODULE=on
WORKDIR /usr/local/go/src/
RUN git clone --branch WebSockets https://github.com/UNIZAR-30226-2021-11/Backend.git
WORKDIR /usr/local/go/src/Backend
RUN go build -o /backend cmd/backend/main.go 


FROM debian:jessie
COPY --from=builder /backend /
COPY --from=builder /database/models.sql /database/models.sql
ENV PORT="11050"
ENV DATABASE_URI="postgres://javiadmin:javiadmin@127.0.0.1:11000/guinote?sslmode=disable"
CMD /backend
