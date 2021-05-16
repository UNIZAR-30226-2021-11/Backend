FROM golang:1.16 as builder


ENV GO111MODULE=on
WORKDIR /usr/local/go/src/
RUN git clone --branch WS-Rejoin https://github.com/UNIZAR-30226-2021-11/Backend.git
WORKDIR /usr/local/go/src/Backend
RUN go build -o /backend cmd/backend/main.go 


FROM debian:jessie
COPY --from=builder /backend /
CMD /backend
