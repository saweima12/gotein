FROM golang:1.21 as builder
LABEL stage=builder

WORKDIR /data
# COPY project into build image.
COPY . .

# Download depedencies.
RUN go mod download

# Build project
RUN CGO_ENABLED=0 go build -o ./app/ ./cmd/...

# Build finish, Copy to runtime
FROM alpine as runtime
# FROM gcr.io/distroless/static-debian12 as runtime
LABEL stage=runtime


WORKDIR /app

COPY --from=builder /data/app ./
COPY --from=builder /data/config.yml ./
COPY --from=builder /data/lang.yml ./

CMD ["./gotein"]
