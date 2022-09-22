FROM golang:1.19

# Copy application data into image
COPY . /go/src/code
WORKDIR /go/src/code
#
#COPY go.mod ./
#COPY go.sum ./
RUN go mod download

# Copy only `.go` files, if you want all files to be copied then replace `with `COPY . .` for the code below.
#COPY *.go .

# Install our third-party application for hot-reloading capability.
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

EXPOSE 9502

ENTRYPOINT CompileDaemon -polling -log-prefix=false -build="go build -o main ." -command="./main" -directory="./"