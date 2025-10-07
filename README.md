# Alumni Management System API

Sistem manajemen alumni dengan API menggunakan Go Fiber dan PostgreSQL. Sistem ini memungkinkan pengelolaan data alumni dan pekerjaan alumni dengan fitur autentikasi berbasis role.

## Fitur Utama

- **Autentikasi JWT** - Login dengan email dan password
- **Role-Based Access Control** - Admin dan User memiliki akses berbeda
- **CRUD Alumni** - Kelola data alumni lengkap
- **CRUD Pekerjaan Alumni** - Kelola riwayat pekerjaan alumni
- **Soft Delete** - Tandai data sebagai terhapus tanpa menghilangkan dari database
- **Pagination & Search** - Tampilkan data dengan pagination dan pencarian
- **Sorting** - Urutkan data berdasarkan kolom tertentu

## Teknologi yang Digunakan

- **Backend**: Go dengan framework Fiber
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Token)
- **Password Hashing**: bcrypt
- **Environment**: godotenv

## Instalasi dan Setup

### 1. Clone Repository
```bash
git clone https://github.com/andinoferdi/go-fiber
cd go-fiber
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Database
- Install PostgreSQL
- Buat database baru
- Jalankan script setup:
```bash
psql -U postgres -f database/setup.sql
```

### 4. Setup Environment
Buat file `.env` di root project:
```env
DB_DSN=postgres://username:password@localhost:5432/alumni_db?sslmode=disable
APP_PORT=3000
JWT_SECRET=your-secret-key-here
```

### 5. Jalankan Aplikasi
```bash
go run main.go
```

Server akan berjalan di `http://localhost:3000`

## Struktur Database

### Tabel Roles
- `id` - Primary key
- `nama` - Nama role (admin, user)

### Tabel Alumni
- `id` - Primary key
- `nim` - Nomor Induk Mahasiswa
- `nama` - Nama lengkap
- `jurusan` - Jurusan kuliah
- `angkatan` - Tahun masuk
- `tahun_lulus` - Tahun lulus
- `email` - Email (unique)
- `password_hash` - Password terenkripsi
- `no_telepon` - Nomor telepon
- `alamat` - Alamat lengkap
- `role_id` - Foreign key ke tabel roles

### Tabel Pekerjaan Alumni
- `id` - Primary key
- `alumni_id` - Foreign key ke tabel alumni
- `nama_perusahaan` - Nama perusahaan
- `posisi_jabatan` - Posisi/jabatan
- `bidang_industri` - Bidang industri
- `lokasi_kerja` - Lokasi kerja
- `gaji_range` - Range gaji
- `tanggal_mulai_kerja` - Tanggal mulai kerja
- `tanggal_selesai_kerja` - Tanggal selesai kerja
- `status_pekerjaan` - Status (aktif, selesai, resigned)
- `deskripsi_pekerjaan` - Deskripsi pekerjaan
- `is_delete` - Timestamp soft delete

## API Endpoints

### Autentikasi
- `POST /go-fiber/login` - Login dan dapatkan token
- `GET /go-fiber/profile` - Ambil data profile dari token

### Alumni (Memerlukan Token)
- `GET /go-fiber/alumni` - Ambil semua alumni dengan pagination
- `GET /go-fiber/alumni/:id` - Ambil alumni berdasarkan ID
- `POST /go-fiber/alumni` - Buat alumni baru (Admin only)
- `PUT /go-fiber/alumni/:id` - Update alumni (Admin only)
- `DELETE /go-fiber/alumni/:id` - Hapus alumni (Admin only)

### Pekerjaan Alumni (Memerlukan Token)
- `GET /go-fiber/pekerjaan` - Ambil semua pekerjaan dengan pagination
- `GET /go-fiber/pekerjaan/:id` - Ambil pekerjaan berdasarkan ID
- `GET /go-fiber/pekerjaan/alumni/:alumni_id` - Ambil pekerjaan berdasarkan alumni ID
- `POST /go-fiber/pekerjaan` - Buat pekerjaan baru (Admin only)
- `PUT /go-fiber/pekerjaan/:id` - Update pekerjaan (Admin only)
- `PUT /go-fiber/pekerjaan/soft-delete/:id` - Soft delete pekerjaan (User/Admin)
- `DELETE /go-fiber/pekerjaan/:id` - Hard delete pekerjaan (Admin only)

### Roles (Memerlukan Token)
- `GET /go-fiber/roles` - Ambil semua roles
- `GET /go-fiber/roles/:id` - Ambil role berdasarkan ID
- `POST /go-fiber/roles` - Buat role baru (Admin only)
- `PUT /go-fiber/roles/:id` - Update role (Admin only)
- `DELETE /go-fiber/roles/:id` - Hapus role (Admin only)

## Parameter Query

### Pagination
- `page` - Halaman (default: 1)
- `limit` - Data per halaman (default: 10, max: 100)

### Sorting
- `sortBy` - Kolom untuk sorting (id, nama, email, dll)
- `order` - Urutan (asc/desc, default: asc)

### Search
- `search` - Kata kunci pencarian (case-insensitive)

### Contoh URL
```
GET /go-fiber/alumni?page=1&limit=5&sortBy=nama&order=desc&search=informatika
```

## Role dan Permission

### Admin
- Dapat mengakses semua endpoint
- Dapat membuat, mengupdate, dan menghapus data
- Dapat soft delete semua pekerjaan alumni

### User
- Dapat membaca data alumni dan pekerjaan
- Dapat soft delete hanya pekerjaan alumni miliknya sendiri
- Tidak dapat membuat, mengupdate, atau menghapus data

## Format Request/Response

### Login Request
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Login Response
```json
{
  "success": true,
  "message": "Login berhasil. Token JWT telah dibuat.",
  "data": {
    "alumni": {
      "id": 1,
      "nama": "John Doe",
      "email": "john@example.com",
      "role": {
        "nama": "admin"
      }
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Create Alumni Request
```json
{
  "nim": "2021001",
  "nama": "John Doe",
  "jurusan": "Teknik Informatika",
  "angkatan": 2021,
  "tahun_lulus": 2025,
  "email": "john@example.com",
  "password": "password123",
  "no_telepon": "081234567890",
  "alamat": "Jl. Contoh No. 1",
  "role_id": 2
}
```

### Pagination Response
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "pages": 5,
    "sortBy": "nama",
    "order": "asc",
    "search": ""
  }
}
```

## Error Handling

Semua error response memiliki format konsisten:
```json
{
  "success": false,
  "message": "Deskripsi error yang jelas"
}
```

### Status Code
- `200` - Success
- `201` - Created
- `400` - Bad Request (validasi error)
- `401` - Unauthorized (token tidak valid)
- `403` - Forbidden (tidak ada permission)
- `404` - Not Found (data tidak ditemukan)
- `500` - Internal Server Error

## Testing API

### Menggunakan Postman
1. Login untuk mendapatkan token
2. Set header `Authorization: Bearer YOUR_TOKEN`
3. Test endpoint sesuai kebutuhan

### Menggunakan curl
```bash
# Login
curl -X POST http://localhost:3000/go-fiber/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"123456"}'

# Ambil data alumni
curl -X GET http://localhost:3000/go-fiber/alumni \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Data Sample

Sistem sudah dilengkapi dengan data sample:
- 5 alumni (1 admin, 4 user)
- 5 pekerjaan alumni
- 2 roles (admin, user)

### Login Credentials
- **Admin**: ahmad.fauzi@email.com / 123456
- **User**: siti.nurhaliza@email.com / 123456

## Troubleshooting

### Database Connection Error
- Pastikan PostgreSQL sudah running
- Cek koneksi string di file .env
- Pastikan database sudah dibuat

### Token Error
- Pastikan token tidak expired
- Cek format header Authorization
- Login ulang jika token invalid

### Permission Error
- Pastikan menggunakan role yang benar
- Admin dapat akses semua endpoint
- User hanya dapat akses terbatas

### Soft Delete vs Hard Delete
- **Soft Delete** (`PUT /soft-delete/:id`): Data ditandai sebagai terhapus dengan timestamp
- **Hard Delete** (`DELETE /:id`): Data benar-benar dihapus dari database
- Data yang di-soft delete tidak akan muncul di hasil query
- **Admin**: Dapat melakukan soft delete dan hard delete semua pekerjaan alumni
- **User**: Hanya dapat soft delete pekerjaan alumni miliknya sendiri

## Kontribusi

1. Fork repository
2. Buat branch fitur baru
3. Commit perubahan
4. Push ke branch
5. Buat Pull Request

## Lisensi

MIT License
