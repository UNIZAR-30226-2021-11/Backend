main:
	CGO_ENABLED=0 go build  -o server cmd/backend/main.go
	scp server aws: