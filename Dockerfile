FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -gcflags="all=-l -B" -ldflags="-s -w -extldflags '-static'" -trimpath -o /app/bot .

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/bot /app/bot

CMD ["/app/bot"]
