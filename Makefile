package:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o aws-registration-callback main.go
	zip aws-registration-callback.zip aws-registration-callback
