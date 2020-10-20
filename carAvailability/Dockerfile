FROM golang:1.15

WORKDIR /app/carAvailability
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build

RUN ls

# Run the executable
CMD ["./carAvailability"]
