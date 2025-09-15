# Alumni Management API

Sistem manajemen data alumni dan riwayat pekerjaan yang dibangun dengan Go Fiber menggunakan clean architecture pattern. API ini menyediakan operasi CRUD lengkap untuk mengelola informasi alumni dan data pekerjaan mereka dengan sistem autentikasi JWT dan role-based access control.

## Persiapan Database

Sebelum menjalankan aplikasi, Anda perlu menyiapkan database PostgreSQL terlebih dahulu.

1. Buat database baru dengan nama `alumni_db`


## Konfigurasi Environment

Buat file `.env` di root directory seperti env_example:


Sesuaikan username, password, dan host database dengan konfigurasi PostgreSQL Anda.

## Cara Menjalankan Aplikasi

```bash

go mod tidy

go run main.go
```

Server akan berjalan di `http://localhost:3000` (atau port yang dikonfigurasi di .env).

## Dokumentasi API

API ini menyediakan endpoint RESTful untuk mengelola data alumni dan pekerjaan dengan sistem autentikasi JWT. Semua endpoint menggunakan format JSON untuk request dan response.

### Sistem Autentikasi

API menggunakan JWT (JSON Web Token) untuk autentikasi dengan role-based access control:

- **Public Endpoints**: Login
- **Protected Endpoints**: Memerlukan token JWT di header Authorization
- **Role-based Access**: Admin dapat akses semua endpoint, User hanya dapat akses GET endpoints

### Login

**Endpoint:** `POST /go-fiber/login`

**Request Body:**
```json
{
    "username": "admin",
    "password": "123456"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Login berhasil",
    "data": {
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@gmail.com",
            "role": "admin",
            "created_at": "2023-01-01T00:00:00Z"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

### Menggunakan Token

Untuk semua endpoint protected, tambahkan header:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

### Access Control

| Endpoint | Admin | User | Keterangan |
|----------|-------|------|------------|
| `POST /login` | ✅ | ✅ | Public endpoint |
| `GET /profile` | ✅ | ✅ | Profile user yang login |
| `GET /alumni` | ✅ | ✅ | Lihat semua alumni |
| `GET /alumni/{id}` | ✅ | ✅ | Lihat alumni by ID |
| `POST /alumni` | ✅ | ❌ | Tambah alumni |
| `PUT /alumni/{id}` | ✅ | ❌ | Update alumni |
| `DELETE /alumni/{id}` | ✅ | ❌ | Hapus alumni |
| `GET /pekerjaan` | ✅ | ✅ | Lihat semua pekerjaan |
| `GET /pekerjaan/{id}` | ✅ | ✅ | Lihat pekerjaan by ID |
| `GET /pekerjaan/alumni/{alumni_id}` | ✅ | ❌ | Lihat pekerjaan by alumni |
| `POST /pekerjaan` | ✅ | ❌ | Tambah pekerjaan |
| `PUT /pekerjaan/{id}` | ✅ | ❌ | Update pekerjaan |
| `DELETE /pekerjaan/{id}` | ✅ | ❌ | Hapus pekerjaan |

### Manajemen Data Alumni

#### Mengambil Semua Alumni
```
GET /go-fiber/alumni
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Mengembalikan daftar lengkap semua alumni yang terdaftar dalam sistem. **Memerlukan autentikasi.**

**Response berhasil:**
```json
{
    "message": "Berhasil mengambil semua data alumni",
    "success": true,
    "data": [
        {
            "id": 1,
            "nim": "123456789",
            "nama": "John Doe",
            "jurusan": "Teknik Informatika",
            "angkatan": 2018,
            "tahun_lulus": 2022,
            "email": "john.doe@email.com",
            "no_telepon": "081234567890",
            "alamat": "Jl. Contoh No. 123",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        }
    ]
}
```

#### Mengambil Data Alumni Tertentu
```
GET /go-fiber/alumni/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Mengambil informasi detail alumni berdasarkan ID yang diberikan. **Memerlukan autentikasi.**

#### Menambah Alumni Baru
```
POST /go-fiber/alumni
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

**Admin Only** - Hanya admin yang dapat menambah alumni baru.

**Request body yang diperlukan:**
```json
{
    "nim": "123456789",
    "nama": "John Doe",
    "jurusan": "Teknik Informatika",
    "angkatan": 2018,
    "tahun_lulus": 2022,
    "email": "john.doe@email.com",
    "no_telepon": "081234567890",
    "alamat": "Jl. Contoh No. 123"
}
```

Field `nim`, `nama`, `jurusan`, dan `email` wajib diisi. Field lainnya bersifat opsional.

#### Memperbarui Data Alumni
```
PUT /go-fiber/alumni/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

**Admin Only** - Hanya admin yang dapat memperbarui data alumni. Menggunakan format request body yang sama dengan endpoint POST. Semua field akan diperbarui sesuai data yang dikirim.

#### Menghapus Alumni
```
DELETE /go-fiber/alumni/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

**Admin Only** - Hanya admin yang dapat menghapus data alumni beserta seluruh riwayat pekerjaan yang terkait (cascade delete).

### Manajemen Data Pekerjaan Alumni

#### Mengambil Semua Data Pekerjaan
```
GET /go-fiber/pekerjaan
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Mengembalikan daftar semua riwayat pekerjaan dari seluruh alumni. **Memerlukan autentikasi.**

#### Mengambil Data Pekerjaan Tertentu
```
GET /go-fiber/pekerjaan/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Mengambil detail pekerjaan berdasarkan ID pekerjaan yang diberikan. **Memerlukan autentikasi.**

#### Mengambil Pekerjaan Berdasarkan Alumni
```
GET /go-fiber/pekerjaan/alumni/{alumni_id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

**Admin Only** - Menampilkan semua riwayat pekerjaan dari alumni tertentu, diurutkan berdasarkan tanggal mulai kerja terbaru.

#### Menambah Data Pekerjaan Baru
```
POST /go-fiber/pekerjaan
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

**Admin Only** - Hanya admin yang dapat menambah data pekerjaan baru.

**Request body:**
```json
{
    "alumni_id": 1,
    "nama_perusahaan": "PT. Tech Indonesia",
    "posisi_jabatan": "Software Engineer",
    "bidang_industri": "Teknologi",
    "lokasi_kerja": "Jakarta",
    "gaji_range": "10-15 juta",
    "tanggal_mulai_kerja": "2022-08-01",
    "tanggal_selesai_kerja": null,
    "status_pekerjaan": "aktif",
    "deskripsi_pekerjaan": "Mengembangkan aplikasi web"
}
```

Field wajib: `alumni_id`, `nama_perusahaan`, `posisi_jabatan`, `bidang_industri`, `lokasi_kerja`, dan `tanggal_mulai_kerja`.

Status pekerjaan yang valid: `aktif`, `selesai`, atau `resigned`. Jika tidak diisi, default adalah `aktif`.

#### Memperbarui Data Pekerjaan
```
PUT /go-fiber/pekerjaan/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

**Admin Only** - Hanya admin yang dapat memperbarui data pekerjaan. Menggunakan format request body yang sama dengan endpoint POST.

#### Menghapus Data Pekerjaan
```
DELETE /go-fiber/pekerjaan/{id}
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

**Admin Only** - Hanya admin yang dapat menghapus riwayat pekerjaan tertentu dari sistem.

## Testing dan Pengujian

### Langkah-langkah Testing dengan Authentication

#### 1. Login untuk Mendapatkan Token

**Request:**
```
POST http://localhost:3000/go-fiber/login
Content-Type: application/json

{
    "username": "admin",
    "password": "123456"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Login berhasil",
    "data": {
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@gmail.com",
            "role": "admin"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

#### 2. Menggunakan Token di Header

Untuk semua endpoint protected, tambahkan header:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

#### 3. Testing Scenarios

**Scenario 1: Login sebagai Admin**
1. Login dengan username: `admin`, password: `123456`
2. Copy token dari response
3. Test semua endpoint dengan token admin
4. Pastikan admin bisa akses semua endpoint

**Scenario 2: Login sebagai User**
1. Login dengan username: `user1`, password: `123456`
2. Copy token dari response
3. Test GET endpoints → harus berhasil
4. Test POST/PUT/DELETE endpoints → harus mendapat 403 Forbidden
5. Test GET `/go-fiber/pekerjaan/alumni/{alumni_id}` → harus 403 Forbidden

**Scenario 3: Uji Token Invalid**
1. Gunakan token yang salah/expired
2. Test endpoint protected → harus mendapat 401 Unauthorized

**Scenario 4: Uji Tanpa Token**
1. Test endpoint protected tanpa header Authorization
2. Harus mendapat 401 Unauthorized

### Menggunakan Postman

Untuk memudahkan testing, ikuti langkah-langkah berikut:

1. **Setup Environment Variable** di Postman:
   - `base_url`: `http://localhost:3000`
   - `token`: (akan diisi setelah login)

2. **Login Request:**
   - Method: `POST`
   - URL: `{{base_url}}/go-fiber/login`
   - Body: JSON dengan username dan password

3. **Setup Authorization Header:**
   - Di setiap request protected, tambahkan header:
   - Key: `Authorization`
   - Value: `Bearer {{token}}`

4. **Test Endpoints:**
   - Jalankan request sesuai kebutuhan testing
   - Periksa response format dan status code

### Testing Manual dengan cURL

**Login untuk mendapatkan token:**
```bash
curl -X POST http://localhost:3000/go-fiber/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'
```

**Testing endpoint GET alumni (dengan token):**
```bash
curl -X GET http://localhost:3000/go-fiber/alumni \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```

**Testing POST alumni baru (admin only):**
```bash
curl -X POST http://localhost:3000/go-fiber/alumni \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "nim": "123456789",
    "nama": "John Doe",
    "jurusan": "Teknik Informatika",
    "angkatan": 2018,
    "tahun_lulus": 2022,
    "email": "john.doe@email.com"
  }'
```

### Expected HTTP Status Codes

- `200 OK`: Request berhasil
- `201 Created`: Resource berhasil dibuat
- `400 Bad Request`: Input validation error
- `401 Unauthorized`: Token missing/invalid/expired
- `403 Forbidden`: Insufficient privileges (wrong role)
- `404 Not Found`: Resource tidak ditemukan
- `500 Internal Server Error`: Server error

## Arsitektur dan Struktur Project

Aplikasi ini menggunakan clean architecture pattern dengan pemisahan layer yang jelas:

```
go-fiber/
├── app/
│   ├── model/           # Definisi struktur data (alumni, pekerjaan, auth)
│   ├── repository/      # Layer akses database (CRUD operations)
│   └── service/         # Logic bisnis aplikasi dan authentication
├── config/              # Konfigurasi aplikasi dan database
├── database/            # Koneksi database dan schema
│   ├── connection.go    # Database connection
│   ├── Alumni&Pekerjaan.sql    # Schema alumni dan pekerjaan + sample data
│   ├── Auth.sql         # Schema authentication + sample users
│   └── setup_complete.sql      # Setup lengkap semua tabel (recommended)
├── helper/              # Utility functions
│   ├── jwt.go          # JWT token generation dan validation
│   ├── password.go     # Password hashing dan verification
│   └── util.go         # General utilities
├── middleware/          # Custom middleware
│   ├── auth.go         # Authentication dan authorization middleware
│   └── logger.go       # Logging middleware
├── route/               # Definisi routing endpoint dengan access control
├── main.go             # Entry point aplikasi
└── .env                # Environment configuration
```

### Penjelasan Layer

- **Model**: Mendefinisikan struktur data alumni, pekerjaan, dan authentication (User, LoginRequest, JWTClaims)
- **Repository**: Menangani operasi database (CRUD) untuk alumni, pekerjaan, dan users
- **Service**: Berisi logic bisnis, validasi, dan authentication services
- **Route**: Mendefinisikan endpoint dengan access control dan menghubungkan dengan service
- **Config**: Konfigurasi aplikasi dan database connection
- **Helper**: Utility functions untuk JWT, password hashing, dan general utilities
- **Middleware**: Authentication, authorization, dan logging middleware

### Fitur Keamanan

- **JWT Authentication**: Token-based authentication dengan expiration 24 jam
- **Password Hashing**: Menggunakan bcrypt untuk keamanan password
- **Role-based Access Control**: Admin dan User dengan permission berbeda
- **Input Validation**: Validasi input untuk semua endpoint
- **Error Handling**: Response error yang konsisten dan informatif

Arsitektur ini memungkinkan code yang mudah dimaintain, testable, scalable, dan aman.
