
-- ========================================
-- DATABASE SETUP SCRIPT
-- ========================================
-- IMPORTANT: This script will DROP existing tables and recreate them
-- Make sure to backup your data before running this script
-- ========================================

-- Drop existing tables in correct order (respecting foreign key constraints)
DROP TABLE IF EXISTS pekerjaan_alumni CASCADE;
DROP TABLE IF EXISTS alumni CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- ========================================
-- CREATE TABLES
-- ========================================

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS alumni (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(100) NOT NULL,
    jurusan VARCHAR(50) NOT NULL,
    angkatan INTEGER NOT NULL,
    tahun_lulus INTEGER NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id INTEGER NOT NULL,
    no_telepon VARCHAR(15),
    alamat TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

-- Drop existing unique constraint on user_id if it exists
ALTER TABLE alumni DROP CONSTRAINT IF EXISTS alumni_user_id_key;

CREATE TABLE IF NOT EXISTS pekerjaan_alumni (
    id SERIAL PRIMARY KEY,
    alumni_id INTEGER NOT NULL,
    nama_perusahaan VARCHAR(100) NOT NULL,
    posisi_jabatan VARCHAR(100) NOT NULL,
    bidang_industri VARCHAR(50) NOT NULL,
    lokasi_kerja VARCHAR(100) NOT NULL,
    gaji_range VARCHAR(50),
    tanggal_mulai_kerja DATE NOT NULL,
    tanggal_selesai_kerja DATE,
    status_pekerjaan VARCHAR(20) DEFAULT 'aktif' CHECK (status_pekerjaan IN ('aktif', 'selesai', 'resigned')),
    deskripsi_pekerjaan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_delete TIMESTAMP NULL,
    FOREIGN KEY (alumni_id) REFERENCES alumni(id) ON DELETE CASCADE
);

-- ========================================
-- INSERT SAMPLE DATA
-- ========================================
INSERT INTO roles (id, nama) VALUES
(1, 'admin'),
(2, 'user')
ON CONFLICT (id) DO NOTHING;

INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, password_hash, role_id, no_telepon, alamat) VALUES
('123456789', 'Andino Ferdiansah', 'Teknik Informatika', 2023, 2028, 'andinoferdiansah@gmail.com', '$2a$10$PXIVDVXOsMef3Pftrw9jSODs1IQY6eW/b7xQhK8CeBj5Ey6.Ke8Ry', 1, '081359528944', 'JL BIBIS TAMA 1A NO 22'),
('987654321', 'Sahrul Alaudin', 'Sistem Informasi', 2019, 2023, 'sahrul@gmail.com', '$2a$10$pdazUPEXpHBmjwLYvn97aex/SBO3/H30r5E/XaxyEA3tHFpI6hCeu', 2, '081987654321', 'JL BUBUTAN'),
('112233445', 'Ahmad Rahman', 'Teknik Komputer', 2020, 2024, 'ahmad@gmail.com', '$2a$10$pdazUPEXpHBmjwLYvn97aex/SBO3/H30r5E/XaxyEA3tHFpI6hCeu', 2, '081112233445', 'Jl SRIKANA'),
('556677889', 'Siti Nurhaliza', 'Manajemen Informatika', 2017, 2021, 'sitinurhaliza@gmail.com', '$2a$10$pdazUPEXpHBmjwLYvn97aex/SBO3/H30r5E/XaxyEA3tHFpI6hCeu', 2, '085556677889', 'JL KARTINI')
ON CONFLICT (nim) DO NOTHING;

INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, status_pekerjaan, deskripsi_pekerjaan) VALUES
(1, 'PT. Tech Indonesia', 'Software Engineer', 'Teknologi', 'Jakarta', '10-15 juta', '2022-08-01', 'aktif', 'Mengembangkan aplikasi web menggunakan Go dan React'),
(1, 'PT. Startup Digital', 'Backend Developer', 'Teknologi', 'Bandung', '8-12 juta', '2022-03-01', 'selesai', 'Mengembangkan API dan microservices'),
(2, 'PT. Data Analytics', 'Data Analyst', 'Data Science', 'Surabaya', '12-18 juta', '2023-09-01', 'aktif', 'Analisis data dan pembuatan dashboard'),
(2, 'PT. E-commerce Pro', 'Business Analyst', 'E-commerce', 'Yogyakarta', '9-13 juta', '2023-01-15', 'selesai', 'Analisis bisnis dan optimasi proses'),
(3, 'PT. Cloud Solutions', 'DevOps Engineer', 'Cloud Computing', 'Jakarta', '15-20 juta', '2024-06-01', 'aktif', 'Mengelola infrastruktur cloud dan CI/CD'),
(4, 'PT. Digital Marketing', 'Digital Marketing Specialist', 'Marketing', 'Bandung', '8-12 juta', '2021-10-01', 'aktif', 'Strategi digital marketing dan social media'),
(4, 'PT. Creative Agency', 'UI/UX Designer', 'Design', 'Jakarta', '7-10 juta', '2021-03-01', 'selesai', 'Desain interface dan user experience');

-- ========================================
-- CREATE INDEXES
-- ========================================

CREATE INDEX IF NOT EXISTS idx_roles_nama ON roles(nama);
CREATE INDEX IF NOT EXISTS idx_alumni_nim ON alumni(nim);
CREATE INDEX IF NOT EXISTS idx_alumni_email ON alumni(email);
CREATE INDEX IF NOT EXISTS idx_alumni_role_id ON alumni(role_id);
CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_id ON pekerjaan_alumni(alumni_id);

-- ========================================
-- VERIFICATION QUERIES
-- ========================================

-- Check if tables were created successfully
SELECT 'Roles table created' as status, COUNT(*) as count FROM roles;
SELECT 'Alumni table created' as status, COUNT(*) as count FROM alumni;
SELECT 'Pekerjaan Alumni table created' as status, COUNT(*) as count FROM pekerjaan_alumni;

-- Show sample data
SELECT 'Sample Roles:' as info;
SELECT * FROM roles;

SELECT 'Sample Alumni:' as info;
SELECT id, nim, nama, email, role_id FROM alumni;

SELECT 'Sample Pekerjaan Alumni:' as info;
SELECT id, alumni_id, nama_perusahaan, posisi_jabatan FROM pekerjaan_alumni LIMIT 5;