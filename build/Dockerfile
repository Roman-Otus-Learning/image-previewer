# Собираем в гошке
FROM golang:1.22 as build

ENV BIN_FILE /opt/image-previewer/image-previewer
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. ${CODE_DIR}

# Собираем статический бинарник Go,
# иначе он не будет работать в alpine образе.
RUN CGO_ENABLED=0 go build \
        -o ${BIN_FILE} cmd/*

FROM alpine:3.15

ENV HTTP_PORT "8000"

ENV CONFIG_FILE /etc/image-previewer/config.yaml
COPY /../configs/config.yaml ${CONFIG_FILE}

ENV BIN_FILE "/opt/image-previewer/image-previewer"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE} image-previewer