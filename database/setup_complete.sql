
CREATE TABLE IF NOT EXISTS alumni (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(100) NOT NULL,
    jurusan VARCHAR(50) NOT NULL,
    angkatan INTEGER NOT NULL,
    tahun_lulus INTEGER NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    no_telepon VARCHAR(15),
    alamat TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
    FOREIGN KEY (alumni_id) REFERENCES alumni(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO users (username, email, password_hash, role) VALUES
('admin', 'admin@gmail.com', 
 '$2a$10$PXIVDVXOsMef3Pftrw9jSODs1IQY6eW/b7xQhK8CeBj5Ey6.Ke8Ry', 'admin'),
('user1', 'user1@gmail.com', 
 '$2a$10$pdazUPEXpHBmjwLYvn97aex/SBO3/H30r5E/XaxyEA3tHFpI6hCeu', 'user'),
('user2', 'user2@gmail.com', 
 '$2a$10$EMT6R8YFy7mUaz6rFfuWGeBtoFP1gg1etwTBMVjIQs2wpdTHdcW', 'user')
ON CONFLICT (username) DO NOTHING;

INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat) VALUES
('123456789', 'John Doe', 'Teknik Informatika', 2018, 2022, 'john.doe@email.com', '081234567890', 'Jl. Contoh No. 123'),
('987654321', 'Jane Smith', 'Sistem Informasi', 2019, 2023, 'jane.smith@email.com', '081987654321', 'Jl. Sample No. 456')
ON CONFLICT (nim) DO NOTHING;

INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, status_pekerjaan, deskripsi_pekerjaan) VALUES
(1, 'PT. Tech Indonesia', 'Software Engineer', 'Teknologi', 'Jakarta', '10-15 juta', '2022-08-01', 'aktif', 'Mengembangkan aplikasi web menggunakan Go dan React'),
(1, 'PT. Startup Digital', 'Backend Developer', 'Teknologi', 'Bandung', '8-12 juta', '2022-03-01', 'selesai', 'Mengembangkan API dan microservices'),
(2, 'PT. Data Analytics', 'Data Analyst', 'Data Science', 'Surabaya', '12-18 juta', '2023-09-01', 'aktif', 'Analisis data dan pembuatan dashboard');

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_alumni_nim ON alumni(nim);
CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_id ON pekerjaan_alumni(alumni_id);