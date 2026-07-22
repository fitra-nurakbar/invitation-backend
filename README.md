# 💌 Invitation App — Backend API

Backend REST API untuk aplikasi undangan digital multi-domain, dibangun dengan Go (Gin), Supabase PostgreSQL, Redis, dan Xendit payment gateway.

---

## 📋 Daftar Isi

- [Fitur](#-fitur)
- [Teknologi](#-teknologi)
- [Arsitektur](#-arsitektur)
- [Struktur Project](#-struktur-project)
- [Prerequisites](#-prerequisites)
- [Instalasi & Setup](#-instalasi--setup)
- [Environment Variables](#-environment-variables)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Database Migration](#-database-migration)
- [Seeder](#-seeder)
- [API Dokumentasi](#-api-dokumentasi)
- [Autentikasi](#-autentikasi)
- [Payment Flow](#-payment-flow)
- [Docker](#-docker)
- [Deploy ke Production](#-deploy-ke-production)

---

## ✨ Fitur

- 🔐 **Autentikasi Multi-Method** — Magic Link, Google OAuth, Email+Password (admin)
- 💳 **Payment Gateway** — Integrasi Xendit (VA, QRIS, E-Wallet, Minimarket)
- 📧 **Email Notifikasi** — Konfirmasi pembayaran via SMTP
- 🎨 **Template Management** — Kelola template undangan dengan akses berbasis pembelian
- 📨 **Undangan Digital** — CRUD undangan dengan slug unik dan data JSONB fleksibel
- 💬 **Buku Tamu** — Sistem ucapan dari tamu undangan
- 🛡️ **Rate Limiting** — Redis-based rate limiter per endpoint
- 🌐 **CORS** — Konfigurasi origin per environment
- 🐳 **Docker Ready** — Docker Compose untuk development & production

---

## 🛠 Teknologi

| Kategori  | Teknologi                         |
| --------- | --------------------------------- |
| Language  | Go 1.25+                          |
| Framework | Gin                               |
| Database  | Supabase (PostgreSQL)             |
| ORM       | GORM                              |
| Cache     | Redis                             |
| Migration | golang-migrate                    |
| Auth      | JWT, Magic Link, Google OAuth 2.0 |
| Payment   | Xendit                            |
| Email     | gomail (SMTP)                     |
| Container | Docker, Docker Compose            |

---

## 🏗 Arsitektur

```
Internet
    │
    ▼
Nginx (port 80/443)        ← reverse proxy + SSL
    │
    ▼
Go App (port 8080)         ← REST API
    ├── Supabase PostgreSQL ← database utama
    ├── Redis               ← rate limiting, magic link token
    └── External Services:
        ├── Xendit          ← payment gateway
        ├── Google OAuth    ← autentikasi
        └── SMTP            ← email notifikasi
```

---

## 📁 Struktur Project

```
invitation-app/
├── config/
│   ├── database.go         # koneksi GORM + PostgreSQL
│   ├── migrate.go          # golang-migrate runner
│   ├── oauth.go            # Google OAuth config
│   └── redis.go            # koneksi Redis
├── handlers/
│   ├── init.go             # inisialisasi semua service
│   ├── auth.go             # login admin, change password
│   ├── magic_link.go       # magic link auth
│   ├── oauth.go            # Google OAuth handler
│   ├── user.go             # CRUD user
│   ├── template.go         # CRUD template
│   ├── invitation.go       # CRUD undangan
│   ├── message.go          # CRUD ucapan tamu
│   ├── order.go            # order & payment
│   └── webhook.go          # Xendit webhook handler
├── middleware/
│   ├── auth.go             # JWT middleware
│   ├── cors.go             # CORS middleware
│   ├── rate_limit.go       # rate limiting middleware
│   └── role.go             # role-based access control
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_templates_table.up.sql
│   ├── 000002_create_templates_table.down.sql
│   ├── 000003_create_invitations_table.up.sql
│   ├── 000003_create_invitations_table.down.sql
│   ├── 000004_create_messages_table.up.sql
│   ├── 000004_create_messages_table.down.sql
│   ├── 000005_create_orders_table.up.sql
│   ├── 000005_create_orders_table.down.sql
│   ├── 000006_create_user_templates_table.up.sql
│   ├── 000006_create_user_templates_table.down.sql
│   ├── 000007_make_password_nullable.up.sql
│   ├── 000007_make_password_nullable.down.sql
│   ├── 000008_add_google_id_to_users.up.sql
│   └── 000008_add_google_id_to_users.down.sql
├── models/
│   ├── user.go
│   ├── template.go
│   ├── invitation.go
│   ├── message.go
│   ├── order.go
│   └── user_template.go
├── routes/
│   └── routes.go
├── seeders/
│   └── seeder.go
├── services/
│   ├── email.go            # email service (gomail)
│   ├── template_access.go  # grant akses template
│   └── xendit.go           # Xendit invoice service
├── utils/
│   └── jwt.go              # JWT generate & validate
├── nginx/
│   └── nginx.conf
├── .env.example
├── .gitignore
├── docker-compose.yml
├── docker-compose.prod.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── main.go
```

---

## ✅ Prerequisites

Pastikan sudah terinstall:

- [Go 1.25+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Git](https://git-scm.com/)

Akun eksternal yang dibutuhkan:

- [Supabase](https://supabase.com) — database PostgreSQL
- [Xendit](https://xendit.co) — payment gateway
- [Google Cloud Console](https://console.cloud.google.com) — OAuth
- Gmail App Password atau [Mailtrap](https://mailtrap.io) — SMTP

---

## 🚀 Instalasi & Setup

### 1. Clone repository

```bash
git clone https://github.com/username/invitation-app.git
cd invitation-app
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup environment variables

```bash
cp .env.example .env
```

Edit `.env` dan isi semua value yang dibutuhkan (lihat [Environment Variables](#-environment-variables)).

### 4. Jalankan Redis (via Docker)

```bash
docker run -d --name redis-invitation -p 6379:6379 redis:alpine
```

### 5. Jalankan aplikasi

```bash
go run .
```

---

## 🔧 Environment Variables

Salin `.env.example` ke `.env` dan isi semua value:

```env
# ─── Database ────────────────────────────────────────────────
DATABASE_URL=postgresql://postgres:PASSWORD@db.REF.supabase.co:5432/postgres

# ─── JWT ─────────────────────────────────────────────────────
JWT_SECRET=random-string-min-32-characters
JWT_EXPIRED_HOURS=24

# ─── App ─────────────────────────────────────────────────────
APP_ENV=development
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000

# ─── Redis ───────────────────────────────────────────────────
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=

# ─── Xendit ──────────────────────────────────────────────────
XENDIT_SECRET_KEY=xnd_development_YOUR_KEY
XENDIT_WEBHOOK_TOKEN=YOUR_WEBHOOK_TOKEN
XENDIT_SUCCESS_URL=http://localhost:3000/payment/success
XENDIT_FAILURE_URL=http://localhost:3000/payment/failed

# ─── SMTP ────────────────────────────────────────────────────
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_SENDER_NAME=Invitation App

# ─── Magic Link ──────────────────────────────────────────────
MAGIC_LINK_BASE_URL=http://localhost:3000
MAGIC_LINK_EXPIRED_MINUTES=15

# ─── Google OAuth ────────────────────────────────────────────
GOOGLE_CLIENT_ID=YOUR_CLIENT_ID.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=YOUR_CLIENT_SECRET
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

---

## 🗄 Database Migration

Migration dijalankan **otomatis saat app start**. Untuk mengelola manual:

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Jalankan semua migration
migrate -path migrations -database $DATABASE_URL up

# Rollback 1 langkah
migrate -path migrations -database $DATABASE_URL down 1

# Cek versi saat ini
migrate -path migrations -database $DATABASE_URL version

# Buat migration baru
migrate create -ext sql -dir migrations -seq nama_migration
```

---

## 🌱 Seeder

```bash
# Jalankan seeder (isi data awal)
go run . --seed

# Via Docker
docker-compose exec app ./invitation-app --seed
```

Data yang di-seed:

| Tipe        | Data                                                                   |
| ----------- | ---------------------------------------------------------------------- |
| Admin       | admin@invitation.com / Admin@12345                                     |
| Users       | 3 user dummy (password: User@12345)                                    |
| Templates   | 5 template (Free, Minimalist, Elegant Rose, Rustic Garden, Royal Gold) |
| Invitations | 3 undangan contoh                                                      |
| Messages    | 5 ucapan contoh                                                        |

---

## 📡 API Dokumentasi

Base URL: `http://localhost:8080/api/v1`

### 🔓 Public Endpoints

| Method | Endpoint                      | Deskripsi               |
| ------ | ----------------------------- | ----------------------- |
| GET    | `/invitations/:slug`          | Detail undangan by slug |
| GET    | `/invitations/:slug/messages` | List ucapan tamu        |
| POST   | `/invitations/:slug/messages` | Kirim ucapan            |

### 🔐 Auth Endpoints

| Method | Endpoint                      | Deskripsi                    |
| ------ | ----------------------------- | ---------------------------- |
| POST   | `/auth/magic-link`            | Request magic link           |
| GET    | `/auth/verify`                | Verify magic link → redirect |
| GET    | `/auth/verify-api`            | Verify magic link → JSON     |
| GET    | `/auth/google`                | Login dengan Google          |
| GET    | `/auth/google/callback`       | Callback Google OAuth        |
| GET    | `/auth/me`                    | Profile user (butuh token)   |
| POST   | `/admin/auth/login`           | Login admin                  |
| POST   | `/admin/auth/create`          | Buat admin baru              |
| POST   | `/admin/auth/change-password` | Ganti password admin         |

### 👤 User Endpoints (butuh JWT)

| Method | Endpoint             | Deskripsi                  |
| ------ | -------------------- | -------------------------- |
| GET    | `/users/:id`         | Detail user                |
| PUT    | `/users/:id`         | Update user                |
| POST   | `/invitations`       | Buat undangan              |
| PUT    | `/invitations/:id`   | Update undangan            |
| DELETE | `/invitations/:id`   | Hapus undangan             |
| POST   | `/orders`            | Buat order template        |
| GET    | `/orders`            | List order saya            |
| GET    | `/orders/:id`        | Detail order               |
| POST   | `/orders/:id/cancel` | Batalkan order             |
| GET    | `/my-templates`      | Template yang sudah dibeli |

### 👑 Admin Endpoints (butuh JWT + role admin)

| Method | Endpoint               | Deskripsi           |
| ------ | ---------------------- | ------------------- |
| GET    | `/admin/users`         | List semua user     |
| DELETE | `/admin/users/:id`     | Hapus user          |
| GET    | `/admin/templates`     | List semua template |
| POST   | `/admin/templates`     | Buat template       |
| PUT    | `/admin/templates/:id` | Update template     |
| DELETE | `/admin/templates/:id` | Hapus template      |
| GET    | `/admin/invitations`   | List semua undangan |
| DELETE | `/admin/messages/:id`  | Hapus ucapan        |
| GET    | `/admin/orders`        | List semua order    |

### 🔔 Webhook

| Method | Endpoint          | Deskripsi                  |
| ------ | ----------------- | -------------------------- |
| POST   | `/webhook/xendit` | Callback pembayaran Xendit |

---

## 🔐 Autentikasi

### User — Magic Link

```bash
# 1. Request magic link
curl -X POST http://localhost:8080/api/v1/auth/magic-link \
  -H "Content-Type: application/json" \
  -d '{"email":"user@gmail.com","name":"Nama User"}'

# 2. Cek email → klik link → dapat JWT

# 3. Gunakan JWT
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### User — Google OAuth

```
Buka browser: http://localhost:8080/api/v1/auth/google
→ Pilih akun Google
→ Redirect ke frontend dengan JWT
```

### Admin — Email + Password

```bash
curl -X POST http://localhost:8080/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@invitation.com","password":"Admin@12345"}'
```

---

## 💳 Payment Flow

```
1. User beli template
   POST /api/v1/orders
   {"template_id": "uuid"}

2. App buat invoice di Xendit
   → Return invoice_url

3. Redirect user ke invoice_url
   → User bayar (VA / QRIS / E-Wallet / Minimarket)

4. Xendit kirim webhook ke /webhook/xendit
   → App update order → PAID
   → App grant akses template ke user
   → App kirim email konfirmasi

5. User bisa pakai template
   GET /api/v1/my-templates
```

### Setup Xendit Webhook

Di Xendit Dashboard → Settings → Webhooks:

- URL: `https://yourdomain.com/webhook/xendit`
- Events: `invoice.paid`, `invoice.expired`

---

## 🐳 Docker

### Development

```bash
# Jalankan semua service
docker-compose up -d

# Lihat log
docker-compose logs -f app

# Masuk ke container
docker-compose exec app sh

# Seed database
docker-compose exec app ./invitation-app --seed

# Stop
docker-compose down
```

### Build Manual

```bash
# Build image
docker build -t invitation-app:latest .

# Cek image
docker images | grep invitation-app
```

### Makefile Shortcuts

```bash
make dev        # go run . (development lokal)
make build      # build Docker image
make up         # docker-compose up -d
make down       # docker-compose down
make logs       # docker-compose logs -f app
make seed       # jalankan seeder di Docker
make rebuild    # down + build + up
```

---

## 🚀 Deploy ke Production

### 1. Siapkan VPS

```bash
# Install Docker di VPS (Ubuntu)
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

### 2. Clone & Setup

```bash
git clone https://github.com/username/invitation-app.git
cd invitation-app
cp .env.example .env
nano .env  # isi semua value production
```

### 3. Setup SSL (Let's Encrypt)

```bash
# Install certbot
sudo apt install certbot

# Generate SSL
sudo certbot certonly --standalone -d yourdomain.com

# Copy ke folder nginx
mkdir -p nginx/ssl
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/
```

### 4. Deploy

```bash
# Build & jalankan
docker build -t invitation-app:latest .
docker-compose -f docker-compose.prod.yml up -d

# Seed (pertama kali)
docker-compose -f docker-compose.prod.yml exec app ./invitation-app --seed

# Lihat log
docker-compose -f docker-compose.prod.yml logs -f app
```

### 5. Update aplikasi

```bash
git pull
docker build -t invitation-app:latest .
docker-compose -f docker-compose.prod.yml up -d --no-deps app
```

---

## ⚡ Rate Limiting

| Endpoint                | Limit   | Per     |
| ----------------------- | ------- | ------- |
| Semua endpoint          | 300 req | 1 menit |
| Auth (login/magic link) | 10 req  | 1 menit |
| API umum                | 60 req  | 1 menit |
| Kirim ucapan            | 20 req  | 1 menit |

---

## 🗃 Skema Database

```
users
├── id (uuid, PK)
├── email (text, unique)
├── name (text)
├── password (text, nullable — hanya admin)
├── google_id (text, nullable)
├── avatar (text, nullable)
├── role (enum: admin | user)
└── created_at

templates
├── id (uuid, PK)
├── name (text)
├── price (integer)
├── is_active (boolean)
├── order_deadline_days (integer)
└── active_days_after (integer)

invitations
├── id (uuid, PK)
├── user_id (FK → users)
├── template_id (FK → templates)
├── slug (text, unique)
├── event_date (date)
├── status (enum: draft | active | expired)
├── expires_at (timestamptz)
└── detail (jsonb)

messages
├── id (uuid, PK)
├── invitation_id (FK → invitations)
├── name (text)
├── message (text)
├── ip_address (inet)
└── created_at

orders
├── id (uuid, PK)
├── user_id (FK → users)
├── template_id (FK → templates)
├── invoice_id (text, unique)
├── invoice_url (text)
├── amount (integer)
├── status (enum: pending | paid | expired | failed)
├── expires_at (timestamptz)
├── paid_at (timestamptz, nullable)
├── payment_method (text)
├── payment_channel (text)
├── created_at
└── updated_at

user_templates
├── id (uuid, PK)
├── user_id (FK → users)
├── template_id (FK → templates)
├── order_id (FK → orders)
└── granted_at
```

---

## 📝 License

MIT License — bebas digunakan dan dimodifikasi.

---

## 👤 Author

Dibuat dengan ❤️ menggunakan Go + Supabase + Xendit

```

```
