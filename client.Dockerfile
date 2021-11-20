FROM golang:1.17.2 as client-build

# Создаём директорию server и переходим в неё.
WORKDIR /app

# Копируем файлы go.mod и go.sum и делаем загрузку, чтобы вовремя пересборки контейнера зависимости
# подтягивались из слоёв.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем все файлы из директории ./client локальной машины в текущую директорию (client) образа.
COPY ./client ./client

# Запускаем компиляцию программы на go и сохраняем полученный бинарный файл server в директорию / образа.
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../rental/server .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../short-cli ./client/cmd/

# Устанавливаем в переменную окружения свою таймзону.
ENV TZ=Europe/Moscow

# Информационная команда показывающая на каком порту будет работать приложение.
EXPOSE 8080

# Устанавливаем по дефолту переменные окружения, которые можно переопределить при запуске контейнера.
ENV SRV_PORT=8035
ENV CLI_PORT=8080

# Запускаем приложение.
CMD ["/short-cli"]
