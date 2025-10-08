# Alumni Management System API

Sistem manajemen alumni dengan API menggunakan Go Fiber dan PostgreSQL.

## Fitur

- **JWT Authentication** - Login dengan email/password
- **Role-Based Access** - Admin dan User permissions
- **CRUD Alumni** - Kelola data alumni
- **CRUD Pekerjaan** - Kelola riwayat pekerjaan alumni
- **Soft Delete** - Tandai data terhapus tanpa menghilangkan
- **Pagination & Search** - Data dengan pagination dan pencarian

## Tech Stack

- **Backend**: Go Fiber
- **Database**: PostgreSQL
- **Auth**: JWT + bcrypt
- **Environment**: godotenv

## Setup

### 1. Install Dependencies
```bash
go mod download
```

### 2. Database Setup
```bash
psql -U postgres -f database/setup.sql
```

### 3. Environment Variables
Buat file `.env`:
```env
DB_DSN=postgres://username:password@localhost:5432/alumni_db?sslmode=disable
APP_PORT=3000
JWT_SECRET=your-jwt-secret-key-here-minimum-32-characters
API_KEY=your-api-key-for-check-endpoint
```

### 4. Run
```bash
go run main.go
```

Server: `http://localhost:3000`

## Database Schema

### Tables
- **roles** - User roles (admin, user)
- **alumni** - Alumni data with role_id FK
- **pekerjaan_alumni** - Job history with alumni_id FK

### Sample Data
- 5 alumni (1 admin, 4 users)
- 5 job records
- 2 roles

## API Endpoints

### Auth
- `POST /go-fiber/login` - Login
- `GET /go-fiber/profile` - Get profile

### Alumni (Requires Token)
- `GET /go-fiber/alumni` - List alumni (pagination)
- `GET /go-fiber/alumni/:id` - Get alumni by ID
- `POST /go-fiber/alumni` - Create alumni (Admin only)
- `PUT /go-fiber/alumni/:id` - Update alumni (Admin only)
- `DELETE /go-fiber/alumni/:id` - Delete alumni (Admin only)
- `POST /go-fiber/alumni/check/:key` - Check alumni by NIM

### Jobs (Requires Token)
- `GET /go-fiber/pekerjaan` - List jobs (pagination)
- `GET /go-fiber/pekerjaan/:id` - Get job by ID
- `GET /go-fiber/pekerjaan/alumni/:alumni_id` - Get jobs by alumni ID
- `POST /go-fiber/pekerjaan` - Create job (Admin only)
- `PUT /go-fiber/pekerjaan/:id` - Update job (Admin only)
- `PUT /go-fiber/pekerjaan/soft-delete/:id` - Soft delete job
- `DELETE /go-fiber/pekerjaan/:id` - Hard delete job (Admin only)

### Roles (Requires Token)
- `GET /go-fiber/roles` - List roles
- `GET /go-fiber/roles/:id` - Get role by ID
- `POST /go-fiber/roles` - Create role (Admin only)
- `PUT /go-fiber/roles/:id` - Update role (Admin only)
- `DELETE /go-fiber/roles/:id` - Delete role (Admin only)

## Query Parameters

### Pagination
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

### Sorting
- `sortBy` - Column name (id, nama, email, etc)
- `order` - Sort order (asc/desc, default: asc)

### Search
- `search` - Search keyword (case-insensitive)

### Example
```
GET /go-fiber/alumni?page=1&limit=5&sortBy=nama&order=desc&search=informatika
```

## Permissions

### Admin
- Full access to all endpoints
- Can create, update, delete all data
- Can soft/hard delete any job

### User
- Read-only access to alumni and jobs
- Can only soft delete own jobs
- Cannot create, update, or hard delete

## Request/Response Examples

### Login
```json
POST /go-fiber/login
{
  "email": "admin@example.com",
  "password": "123456"
}
```

### Create Alumni
```json
POST /go-fiber/alumni
{
  "nim": "2021001",
  "nama": "John Doe",
  "jurusan": "Teknik Informatika",
  "angkatan": 2021,
  "tahun_lulus": 2025,
  "email": "john@example.com",
  "password": "password123",
  "role_id": 2
}
```

### Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...}
}
```

## Error Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Testing

### Login Credentials
- **Admin**: andinoferdiansah@gmail.com / 123456
- **User**: siti.nurhaliza@email.com / 123456

### Using curl
```bash
# Login
curl -X POST http://localhost:3000/go-fiber/login \
  -H "Content-Type: application/json" \
  -d '{"email":"andinoferdiansah@gmail.com","password":"123456"}'

# Get alumni (with token)
curl -X GET http://localhost:3000/go-fiber/alumni \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Notes

- All endpoints except login require JWT token in Authorization header
- Soft delete marks data as deleted with timestamp
- Hard delete permanently removes data from database
- Users can only soft delete their own job records
- Admins have full access to all operations