

-- Drop tables jika sudah ada (untuk development/testing)
DROP TABLE IF EXISTS pekerjaan_alumni CASCADE;
DROP TABLE IF EXISTS alumni CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- Tabel Role
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample roles
INSERT INTO roles (id, nama) VALUES
(1, 'admin'),
(2, 'user');

-- Tabel Alumni (dengan role_id)
CREATE TABLE alumni (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(100) NOT NULL,
    jurusan VARCHAR(50) NOT NULL,
    angkatan INTEGER NOT NULL,
    tahun_lulus INTEGER NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    no_telepon VARCHAR(15),
    alamat TEXT,
    role_id INTEGER NOT NULL DEFAULT 2,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

-- Tabel Pekerjaan Alumni
CREATE TABLE pekerjaan_alumni (
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
    is_delete TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (alumni_id) REFERENCES alumni(id) ON DELETE CASCADE
);

-- Insert sample alumni dengan password hash (password: "123456")
INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, password_hash, no_telepon, alamat, role_id) VALUES
('2021001', 'Andino Ferdiansah', 'Teknik Informatika', 2021, 2025, 'andinoferdiansah@gmail.com', '$2y$10$ALjoqFyHgQXMrP/t/X661.cv5hYX7fVEQy5E8fXyog97BGaCn7066', '081359528944', 'JL BIBIS TAMA 1A NO 22', 1),
('2021002', 'Siti Nurhaliza', 'Sistem Informasi', 2021, 2025, 'siti.nurhaliza@email.com', '$2y$10$ALjoqFyHgQXMrP/t/X661.cv5hYX7fVEQy5E8fXyog97BGaCn7066', '081234567891', 'Jl. Diponegoro No. 2, Malang', 2),
('2020001', 'Budi Santoso', 'Teknik Informatika', 2020, 2024, 'budi.santoso@email.com', '$2y$10$ALjoqFyHgQXMrP/t/X661.cv5hYX7fVEQy5E8fXyog97BGaCn7066', '081234567892', 'Jl. Sudirman No. 3, Jakarta', 2),
('2022001', 'Maria Garcia', 'Teknik Informatika', 2022, 2026, 'maria.garcia@email.com', '$2y$10$ALjoqFyHgQXMrP/t/X661.cv5hYX7fVEQy5E8fXyog97BGaCn7066', '081234567893', 'Jl. Gatot Subroto No. 4, Bandung', 2),
('2022002', 'John Smith', 'Sistem Informasi', 2022, 2026, 'john.smith@email.com', '$2y$10$ALjoqFyHgQXMrP/t/X661.cv5hYX7fVEQy5E8fXyog97BGaCn7066', '081234567894', 'Jl. Thamrin No. 5, Medan', 2);

-- Insert sample pekerjaan
INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, status_pekerjaan, deskripsi_pekerjaan) VALUES
(1, 'PT. Tech Solutions', 'Software Developer', 'Teknologi', 'Jakarta', '5-8 juta', '2025-01-15', 'aktif', 'Mengembangkan aplikasi web menggunakan Go dan React'),
(2, 'PT. Digital Innovation', 'System Analyst', 'Teknologi', 'Surabaya', '6-9 juta', '2025-02-01', 'aktif', 'Menganalisis kebutuhan sistem dan merancang solusi IT'),
(3, 'PT. Data Analytics', 'Data Scientist', 'Teknologi', 'Bandung', '8-12 juta', '2024-12-01', 'aktif', 'Menganalisis data besar untuk insights bisnis'),
(4, 'PT. Cloud Computing', 'DevOps Engineer', 'Teknologi', 'Jakarta', '7-10 juta', '2025-03-01', 'aktif', 'Mengelola infrastruktur cloud dan CI/CD pipeline'),
(5, 'PT. Mobile Apps', 'Mobile Developer', 'Teknologi', 'Surabaya', '6-9 juta', '2025-02-15', 'aktif', 'Mengembangkan aplikasi mobile menggunakan Flutter');

-- Buat index untuk performa
CREATE INDEX IF NOT EXISTS idx_alumni_nim ON alumni(nim);
CREATE INDEX IF NOT EXISTS idx_alumni_email ON alumni(email);
CREATE INDEX IF NOT EXISTS idx_alumni_role_id ON alumni(role_id);
CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_id ON pekerjaan_alumni(alumni_id);
CREATE INDEX IF NOT EXISTS idx_pekerjaan_status ON pekerjaan_alumni(status_pekerjaan);
