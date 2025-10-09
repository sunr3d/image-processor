# Image Processor

Микросервисная система для обработки изображений с асинхронной архитектурой на основе Kafka.

## Возможности

- Загрузка изображений через REST API и веб-интерфейс
- Асинхронная обработка через Kafka
- Три типа обработки: resize (800px), thumbnail (200x200px), watermark
- Веб-интерфейс для загрузки и просмотра результатов
- Поиск по ID изображения
- Graceful shutdown и структурированное логирование

## Архитектура

### Компоненты

- **API Service** - HTTP сервер для загрузки и получения изображений
- **Worker Service** - асинхронная обработка изображений
- **Kafka** - брокер сообщений для связи между сервисами
- **File Storage** - хранение оригинальных и обработанных изображений
- **Metadata Storage** - хранение метаданных изображений

## Технологии

- Go 1.24
- Kafka
- Docker & Docker Compose
- wb-go/wbf
- disintegration/imaging
- HTML/JavaScript

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.24+ (для разработки)

### Запуск

1. Клонируйте репозиторий:

   ```bash
   git clone <repository-url>
   cd image-processor
   ```

2. Запустите все сервисы:

   ```bash
   make up
   ```

3. Откройте веб-интерфейс: http://localhost:8080

### Доступные команды

```bash
make up      # Запуск всех сервисов
make down    # Остановка всех сервисов
make rebuild # Пересборка и запуск
make clean   # Остановка и удаление volumes
make logs    # Просмотр логов API сервиса
make test    # Запуск тестов
```

## API Endpoints

### Загрузка изображения

```http
POST /upload
Content-Type: multipart/form-data

Form data:
- image: файл изображения (JPEG, PNG, GIF)
```

Ответ:

```json
{
  "id": "uuid",
  "status": "uploaded",
  "message": "Изображение успешно загружено"
}
```

### Получение изображения

```http
GET /image/{id}?type={type}
```

Параметры:

- `id` - ID изображения
- `type` - тип изображения: `original`, `resized`, `thumbnail`, `watermarked`

### Получение статуса

```http
GET /status/{id}
```

Ответ:

```json
{
  "id": "uuid",
  "status": "completed|processing|failed",
  "message": "Описание ошибки (если есть)"
}
```

### Удаление изображения

```http
DELETE /image/{id}
```

Ответ:

```json
{
  "status": "deleted",
  "message": "Изображение успешно удалено"
}
```

## Веб-интерфейс

Веб-интерфейс доступен по адресу `http://localhost:8080` и предоставляет:

- Загрузку изображений через выбор файла
- Отображение ID загруженного изображения
- Поиск по ID для просмотра существующих изображений
- Просмотр статуса обработки в реальном времени
- Отображение результатов (оригинал, resized, thumbnail, watermarked)
- Удаление изображений

## Конфигурация

Конфигурация задается через переменные окружения:

```bash
HTTP_PORT=8080                    # Порт HTTP сервера
LOG_LEVEL=info                    # Уровень логирования
KAFKA_BROKERS=kafka:29092         # Адреса Kafka брокеров
KAFKA_TOPIC=image-processing      # Топик Kafka
KAFKA_GROUP=image-processor-group # Группа потребителей
STORAGE_PATH=/app/storage         # Путь к хранилищу файлов
METADATA_PATH=/app/metadata       # Путь к хранилищу метаданных
THUMBNAIL_SIZE=200                # Размер миниатюры
RESIZE_WIDTH=800                  # Ширина для resize
```

## Тестирование

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
go test -v -cover ./...

# Запуск тестов конкретного пакета
go test -v ./internal/services/imagesvc
```

## Структура проекта

```
├── cmd/                    # Точки входа
│   ├── app/               # API сервис
│   └── worker/            # Worker сервис
├── internal/              # Внутренняя логика
│   ├── config/           # Конфигурация
│   ├── entrypoint/       # Инициализация сервисов
│   ├── handlers/         # HTTP обработчики
│   ├── infra/            # Инфраструктурный слой
│   │   ├── broker/       # Kafka (Publisher/Subscriber)
│   │   └── storage/      # Хранилища (File/Metadata)
│   ├── interfaces/       # Интерфейсы
│   ├── server/           # HTTP сервер
│   └── services/         # Бизнес-логика
├── models/               # Доменные модели
├── web/                  # Веб-интерфейс
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile.app        # Dockerfile для API
├── Dockerfile.worker     # Dockerfile для Worker
└── Makefile             # Команды для разработки
```

## Архитектурные принципы

Проект следует принципам Clean Architecture:

- **Domain Layer** (`models/`) - доменные модели и бизнес-правила
- **Application Layer** (`services/`) - бизнес-логика и use cases
- **Infrastructure Layer** (`infra/`) - внешние зависимости (Kafka, файловая система)
- **Presentation Layer** (`handlers/`) - HTTP API и веб-интерфейс

### SOLID принципы

- **Single Responsibility** - каждый сервис отвечает за одну задачу
- **Open/Closed** - интерфейсы позволяют расширять функциональность
- **Liskov Substitution** - реализации интерфейсов взаимозаменяемы
- **Interface Segregation** - интерфейсы разделены по назначению
- **Dependency Inversion** - зависимости направлены к абстракциям

## Отладка

### Просмотр логов

```bash
# Логи API сервиса
make logs

# Логи всех сервисов
docker compose logs -f

# Логи конкретного сервиса
docker compose logs -f app
docker compose logs -f worker
```
