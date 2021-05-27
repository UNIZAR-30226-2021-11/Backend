FROM golang:1.16 as builder


ENV GO111MODULE=on
WORKDIR /usr/local/go/src/
RUN git clone --branch WebSockets https://github.com/UNIZAR-30226-2021-11/Backend.git
WORKDIR /usr/local/go/src/Backend
RUN go build -o /backend cmd/backend/main.go 


FROM debian:jessie
COPY --from=builder /backend /
COPY --from=builder /usr/local/go/src/Backend/database/models.sql /database/models.sql
ENV PORT="11050"
ENV HOST="15.188.14.243"
ENV DATABASE_URI="postgres://javiadmin:javiadmin@15.188.14.243:11000/guinote?sslmode=disable"
CMD /backend
