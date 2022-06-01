FROM --platform=amd64 golang:alpine as builder
ENV VERSION v1.0.0
WORKDIR /app
COPY go.mod go.sum  /app/
RUN go mod download
COPY . .
RUN ls
RUN go build -buildvcs=false -v -ldflags="-X 'main.Version=$VERSION'" -o app

FROM --platform=amd64 alpine as prod
WORKDIR /app
COPY --from=builder /app/app /app
CMD [ "/app/app" ]
