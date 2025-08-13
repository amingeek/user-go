# user-go

> Backend service template in Go — OTP login & user management

---

## 🔥 خلاصه

این مخزن (`user-go`) یک سرویس بک‌اند نوشته‌شده با Go است که پیاده‌سازی OTP برای ورود/ثبت‌نام و مدیریت کاربران را هدف دارد. در این README یک راهنمای کامل، دستورالعمل‌های راه‌اندازی محلی و با Docker، اجرای تست‌ها، و نکات رفع اشکال آورده شده — به‌صورت زیبا و با بخش‌های کد جدا برای کپی آسان.

---

## 📦 امکانات (Features)

* OTP بر پایه شماره تلفن (چاپ شده در کنسول — بدون ارسال SMS)
* ذخیره موقتی OTP (در حافظه یا دیتابیس configurable)
* انقضای OTP پس از 2 دقیقه
* ثبت‌نام / ورود بر پایه OTP
* مدیریت کاربران پایه (CRUD)
* تست‌های واحد و integration-ready

---

## 🧰 پیش‌نیازها

* Go 1.20+
* Docker & Docker Compose (برای راه‌اندازی سریع دیتابیس و محیط تست)
* PostgreSQL (در صورت اجرای محلی بدون Docker)

---

## ⚙️ ساختار پروژه (نمونه)

```
user-go/
├── cmd/               # entrypoints (main)
├── internal/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   ├── middleware/
│   └── cache/
├── migrations/        # SQL migrations (اختیاری)
├── Dockerfile
├── Dockerfile.test
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 🧭 متغیرهای محیطی (Environment variables)

لیست متغیرهای مهم که سرویس از آن‌ها استفاده می‌کند:

```env
# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
DB_HOST=db
DB_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=dbname

# App
PORT=8080
JWT_SECRET=your_jwt_secret_here
OTP_EXPIRATION_SECONDS=120

# (اختیاری) logging, debug
LOG_LEVEL=debug
```

> تست‌ها معمولاً دنبال `DATABASE_URL` می‌گردند — اگر این متغیر تنظیم نشود، خطایی شبیه `DATABASE_URL env var not set` خواهید دید.

---

## 🚀 راه‌اندازی سریع (با Docker Compose)

پیشنهاد می‌کنم از Docker Compose برای راه‌اندازی یک دیتابیس Postgres و اجرای اپ یا تست‌ها استفاده کنید.

### نمونه `docker-compose.yml`

```yaml
version: "3.8"

services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: usergolang
      POSTGRES_PASSWORD: admin1234
      POSTGRES_DB: userdb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U usergolang -d userdb"]
      interval: 2s
      timeout: 5s
      retries: 30

  app:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=postgres://usergolang:admin1234@db:5432/userdb?sslmode=disable
      - DB_HOST=db
      - DB_PORT=5432
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "8080:8080"
    tty: true

  test:
    build:
      context: .
      dockerfile: Dockerfile.test
    environment:
      - DATABASE_URL=postgres://usergolang:admin1234@db:5432/userdb?sslmode=disable
      - DB_HOST=db
      - DB_PORT=5432
      - POSTGRES_USER=usergolang
    depends_on:
      db:
        condition: service_healthy
    tty: true

volumes:
  pgdata:
```

### اجرای اپ

```bash
# ساخت و ران اپ و دیتابیس (برای توسعه)
docker-compose up --build app
# سپس اپ روی localhost:8080 در دسترس است (اگر اپ همین پورت را استفاده کند)
```

### اجرای تست‌ها

```bash
# اجرا کردن سرویس تست (همراه db) و برگشت exit code از کانتینر test
docker-compose up --build --abort-on-container-exit --exit-code-from test test
```

### پاک کردن داده‌ها

```bash
docker-compose down -v
```

---

## 🧩 نمونه Dockerfile (برای اپ)

```dockerfile
# Dockerfile
FROM golang:1.20-bullseye AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# اگر main در مسیر مشخصی است (مثلاً ./cmd/server) آن را تنظیم کنید
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app ./...

FROM gcr.io/distroless/static-debian11
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
EXPOSE 8080
ENV DATABASE_URL=""
CMD ["/usr/local/bin/app"]
```

## 🧪 Dockerfile.test (اجرای تست‌ها)

```dockerfile
FROM golang:1.20-bullseye
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD sh -c '\
  echo "Waiting for Postgres at ${DB_HOST:-db}:${DB_PORT:-5432} ..."; \
  until pg_isready -h "${DB_HOST:-db}" -p "${DB_PORT:-5432}" -U "${POSTGRES_USER:-usergolang}" -d "${POSTGRES_DB:-userdb}"; do \
    sleep 1; \
  done; \
  echo "Postgres is ready — running tests"; \
  go test ./...'
```

---

## 🧭 اجرای لوکال (بدون داکر)

1. PostgreSQL را نصب و اجرا کنید.
2. یک دیتابیس و یوزر بسازید و `DATABASE_URL` را export کنید:

```bash
export DATABASE_URL="postgres://usergolang:admin1234@localhost:5432/userdb?sslmode=disable"
go test ./...
# یا اجرای اپ
go run ./cmd/server
```

---

## ✅ نکات متداول و رفع اشکال

* ⚠️ `DATABASE_URL env var not set` — مطمئن شوید متغیر محیطی `DATABASE_URL` تنظیم شده است؛ برای تست‌ها یا docker-compose آن را تعریف کنید.
* اگر تست‌ها درباره جداول خطا دادند، ممکن است نیاز به migrations داشته باشید. می‌توانید یک اسکریپت `migrations/init.sql` بسازید و آن را قبل از `go test` اجرا کنید:

```bash
psql "$DATABASE_URL" -f migrations/init.sql
```

* برای مشاهده لاگ‌های کامل کانتینر هنگام اجرای docker-compose از `docker-compose logs -f` استفاده کنید.

---

## 🧪 تست و CI (نمونه GitHub Actions)

یک Workflow ساده برای اجرای تست‌ها در CI:

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: usergolang
          POSTGRES_PASSWORD: admin1234
          POSTGRES_DB: userdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 2s
          --health-timeout 5s
          --health-retries 10
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      - name: Wait for Postgres
        run: |
          for i in {1..20}; do
            pg_isready -h localhost -p 5432 && break || sleep 1
          done
      - name: Run tests
        env:
          DATABASE_URL: postgres://usergolang:admin1234@localhost:5432/userdb?sslmode=disable
        run: go test ./...
```

---

## 📚 پیشنهادات توسعه

* اضافه کردن migrations (برای مثال با Goose یا sql-migrate)
* ارسال OTP از طریق سرویس SMS (در صورت نیاز)
* رمزنگاری و مدیریت امن‌تر secretها (Vault / GitHub Secrets)
* مستندسازی API با Swagger/OpenAPI

---

## 💬 کمک بیشتر

اگر دوست داری من این README را به انگلیسی هم ترجمه کنم، یا فایل‌های Docker/Docker-Compose و اسکریپت migrations را برایت آماده کنم — فقط بگو. همچنین می‌تونم یک workflow کامل GitHub Actions یا یک makefile هم اضافه کنم.

---

## 📝 لایسنس

اضافه کن که قصد داری از چه لایسنسی استفاده کنی (مثلاً MIT). فعلاً فایل README اینجا بدون لایسنس خاص است.

---

با آرزوی موفقیت — اگر می‌خواهی من همین‌جا فایل‌ها را ایجاد/به‌روز کنم بگو تا انجام دهم.
