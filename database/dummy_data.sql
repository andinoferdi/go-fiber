-- Add more dummy data for alumni table to test pagination
INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat) VALUES
('20180001', 'Ahmad Budiman', 'Teknik Informatika', 2018, 2022, 'ahmad.budiman@email.com', '081234567891', 'Jl. Merdeka No. 1'),
('20180002', 'Siti Nurhaliza', 'Sistem Informasi', 2018, 2022, 'siti.nurhaliza@email.com', '081234567892', 'Jl. Sudirman No. 2'),
('20180003', 'Budi Santoso', 'Teknik Informatika', 2018, 2022, 'budi.santoso@email.com', '081234567893', 'Jl. Thamrin No. 3'),
('20180004', 'Ani Sari', 'Sistem Informasi', 2018, 2022, 'ani.sari@email.com', '081234567894', 'Jl. Gatot Subroto No. 4'),
('20180005', 'Dedi Rahman', 'Teknik Informatika', 2018, 2022, 'dedi.rahman@email.com', '081234567895', 'Jl. Kuningan No. 5'),
('20190001', 'Eka Putri', 'Sistem Informasi', 2019, 2023, 'eka.putri@email.com', '081234567896', 'Jl. Senayan No. 6'),
('20190002', 'Fajar Wijaya', 'Teknik Informatika', 2019, 2023, 'fajar.wijaya@email.com', '081234567897', 'Jl. Pancoran No. 7'),
('20190003', 'Gita Maharani', 'Sistem Informasi', 2019, 2023, 'gita.maharani@email.com', '081234567898', 'Jl. Kemang No. 8'),
('20190004', 'Hadi Kusuma', 'Teknik Informatika', 2019, 2023, 'hadi.kusuma@email.com', '081234567899', 'Jl. Blok M No. 9'),
('20190005', 'Ira Suhartono', 'Sistem Informasi', 2019, 2023, 'ira.suhartono@email.com', '081234567800', 'Jl. Pondok Indah No. 10'),
('20200001', 'Joko Widodo', 'Teknik Informatika', 2020, 2024, 'joko.widodo@email.com', '081234567801', 'Jl. Menteng No. 11'),
('20200002', 'Kartika Sari', 'Sistem Informasi', 2020, 2024, 'kartika.sari@email.com', '081234567802', 'Jl. Cikini No. 12'),
('20200003', 'Lukman Hakim', 'Teknik Informatika', 2020, 2024, 'lukman.hakim@email.com', '081234567803', 'Jl. Salemba No. 13'),
('20200004', 'Maya Sari', 'Sistem Informasi', 2020, 2024, 'maya.sari@email.com', '081234567804', 'Jl. Tebet No. 14'),
('20200005', 'Nur Hidayat', 'Teknik Informatika', 2020, 2024, 'nur.hidayat@email.com', '081234567805', 'Jl. Cempaka Putih No. 15')
ON CONFLICT (nim) DO NOTHING;

-- Add more dummy data for pekerjaan_alumni table
-- Use dynamic alumni_id based on existing alumni records
INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan) 
SELECT 
    a.id as alumni_id,
    'PT. Garuda Technology' as nama_perusahaan,
    'Full Stack Developer' as posisi_jabatan,
    'Teknologi' as bidang_industri,
    'Jakarta' as lokasi_kerja,
    '12-16 juta' as gaji_range,
    '2022-10-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengembangkan aplikasi enterprise dengan Java Spring Boot' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '123456789'
UNION ALL
SELECT 
    a.id as alumni_id,
    'CV. Web Solutions' as nama_perusahaan,
    'Frontend Developer' as posisi_jabatan,
    'Teknologi' as bidang_industri,
    'Bandung' as lokasi_kerja,
    '8-10 juta' as gaji_range,
    '2022-01-01'::date as tanggal_mulai_kerja,
    '2022-09-30'::date as tanggal_selesai_kerja,
    'selesai' as status_pekerjaan,
    'Membuat website responsive dengan React' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '123456789'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Bank Digital' as nama_perusahaan,
    'Business Analyst' as posisi_jabatan,
    'Perbankan' as bidang_industri,
    'Jakarta' as lokasi_kerja,
    '15-20 juta' as gaji_range,
    '2022-11-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Menganalisis kebutuhan bisnis untuk produk digital banking' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '987654321'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. E-Commerce Nusantara' as nama_perusahaan,
    'Backend Developer' as posisi_jabatan,
    'E-Commerce' as bidang_industri,
    'Surabaya' as lokasi_kerja,
    '10-14 juta' as gaji_range,
    '2022-12-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengembangkan microservices untuk platform e-commerce' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20180001'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Fintech Innovation' as nama_perusahaan,
    'Data Scientist' as posisi_jabatan,
    'Fintech' as bidang_industri,
    'Jakarta' as lokasi_kerja,
    '18-25 juta' as gaji_range,
    '2023-01-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Machine learning untuk credit scoring dan fraud detection' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20180002'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Mobile Solutions' as nama_perusahaan,
    'Mobile Developer' as posisi_jabatan,
    'Teknologi' as bidang_industri,
    'Bandung' as lokasi_kerja,
    '11-15 juta' as gaji_range,
    '2023-02-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengembangkan aplikasi mobile dengan Flutter' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20180003'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Cloud Computing' as nama_perusahaan,
    'DevOps Engineer' as posisi_jabatan,
    'Teknologi' as bidang_industri,
    'Jakarta' as lokasi_kerja,
    '16-22 juta' as gaji_range,
    '2023-03-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengelola infrastructure cloud dengan AWS dan Docker' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20180004'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Game Studio' as nama_perusahaan,
    'Game Developer' as posisi_jabatan,
    'Gaming' as bidang_industri,
    'Yogyakarta' as lokasi_kerja,
    '9-13 juta' as gaji_range,
    '2023-04-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengembangkan game mobile dengan Unity' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20180005'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. Cyber Security' as nama_perusahaan,
    'Security Analyst' as posisi_jabatan,
    'Keamanan Siber' as bidang_industri,
    'Jakarta' as lokasi_kerja,
    '14-20 juta' as gaji_range,
    '2023-05-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Monitoring dan analisis keamanan sistem' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20190001'
UNION ALL
SELECT 
    a.id as alumni_id,
    'PT. IoT Solutions' as nama_perusahaan,
    'IoT Developer' as posisi_jabatan,
    'IoT' as bidang_industri,
    'Surabaya' as lokasi_kerja,
    '12-17 juta' as gaji_range,
    '2023-06-01'::date as tanggal_mulai_kerja,
    NULL as tanggal_selesai_kerja,
    'aktif' as status_pekerjaan,
    'Mengembangkan solusi Internet of Things' as deskripsi_pekerjaan
FROM alumni a WHERE a.nim = '20190002';