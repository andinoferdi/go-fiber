package database

import (
	"context"
	"log"
	"time"

	utilsmongo "go-fiber/utils/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunMigrations(db *mongo.Database) error {
	log.Println("Starting MongoDB migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dropCollections(ctx, db); err != nil {
		return err
	}

	if err := createIndexes(ctx, db); err != nil {
		return err
	}

	if err := seedData(ctx, db); err != nil {
		return err
	}

	log.Println("MongoDB migrations completed successfully!")
	return nil
}

func dropCollections(ctx context.Context, db *mongo.Database) error {
	log.Println("Dropping existing collections...")

	collections := []string{"alumni", "pekerjaan_alumni", "files"}
	
	for _, collectionName := range collections {
		collection := db.Collection(collectionName)
		
		names, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
		if err != nil {
			log.Printf("Warning: Could not check collection %s: %v", collectionName, err)
			continue
		}
		
		if len(names) > 0 {
			if err := collection.Drop(ctx); err != nil {
				log.Printf("Warning: Could not drop collection %s: %v", collectionName, err)
			} else {
				log.Printf("Dropped collection: %s", collectionName)
			}
		}
	}

	return nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	log.Println("Creating indexes...")


	alumniCollection := db.Collection("alumni")
	alumniIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "nim", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "role", Value: 1}},
		},
	}
	if _, err := alumniCollection.Indexes().CreateMany(ctx, alumniIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for alumni collection")

	pekerjaanCollection := db.Collection("pekerjaan_alumni")
	pekerjaanIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "alumni_info.alumni_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status_pekerjaan", Value: 1}},
		},
	}
	if _, err := pekerjaanCollection.Indexes().CreateMany(ctx, pekerjaanIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for pekerjaan_alumni collection")

	filesCollection := db.Collection("files")
	filesIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "alumni_info.alumni_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "file_type", Value: 1}},
		},
	}
	if _, err := filesCollection.Indexes().CreateMany(ctx, filesIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for files collection")

	return nil
}

func seedData(ctx context.Context, db *mongo.Database) error {
	log.Println("Seeding initial data...")

	alumniIDs, err := seedAlumni(ctx, db)
	if err != nil {
		return err
	}

	if err := seedPekerjaanAlumni(ctx, db, alumniIDs); err != nil {
		return err
	}

	log.Println("Data seeding completed successfully!")
	return nil
}


func seedAlumni(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	log.Println("Seeding alumni...")

	passwordHash, err := utilsmongo.HashPassword("123456")
	if err != nil {
		return nil, err
	}

	alumni := []interface{}{
		bson.M{
			"nim":           "2021001",
			"nama":          "Andino Ferdiansah",
			"jurusan":       "Teknik Informatika",
			"angkatan":      2021,
			"tahun_lulus":   2025,
			"email":         "andinoferdiansah@gmail.com",
			"password_hash": passwordHash,
			"no_telepon":    "081359528944",
			"alamat":        "JL BIBIS TAMA 1A NO 22",
			"role":          "admin",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"nim":           "2021002",
			"nama":          "Siti Nurhaliza",
			"jurusan":       "Sistem Informasi",
			"angkatan":      2021,
			"tahun_lulus":   2025,
			"email":         "siti.nurhaliza@email.com",
			"password_hash": passwordHash,
			"no_telepon":    "081234567891",
			"alamat":        "Jl. Diponegoro No. 2, Malang",
			"role":          "user",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"nim":           "2020001",
			"nama":          "Budi Santoso",
			"jurusan":       "Teknik Informatika",
			"angkatan":      2020,
			"tahun_lulus":   2024,
			"email":         "budi.santoso@email.com",
			"password_hash": passwordHash,
			"no_telepon":    "081234567892",
			"alamat":        "Jl. Sudirman No. 3, Jakarta",
			"role":          "user",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"nim":           "2022001",
			"nama":          "Maria Garcia",
			"jurusan":       "Teknik Informatika",
			"angkatan":      2022,
			"tahun_lulus":   2026,
			"email":         "maria.garcia@email.com",
			"password_hash": passwordHash,
			"no_telepon":    "081234567893",
			"alamat":        "Jl. Gatot Subroto No. 4, Bandung",
			"role":          "user",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"nim":           "2022002",
			"nama":          "John Smith",
			"jurusan":       "Sistem Informasi",
			"angkatan":      2022,
			"tahun_lulus":   2026,
			"email":         "john.smith@email.com",
			"password_hash": passwordHash,
			"no_telepon":    "081234567894",
			"alamat":        "Jl. Thamrin No. 5, Medan",
			"role":          "user",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
	}

	collection := db.Collection("alumni")
	result, err := collection.InsertMany(ctx, alumni)
	if err != nil {
		return nil, err
	}

	alumniIDs := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		alumniIDs[i] = id.(primitive.ObjectID)
	}

	log.Printf("Inserted %d alumni", len(result.InsertedIDs))
	return alumniIDs, nil
}

func seedPekerjaanAlumni(ctx context.Context, db *mongo.Database, alumniIDs []primitive.ObjectID) error {
	log.Println("Seeding pekerjaan alumni...")

	pekerjaan := []interface{}{
		bson.M{
			"alumni_info": bson.M{
				"alumni_id": alumniIDs[0],
				"nim":       "2021001",
				"nama":      "Andino Ferdiansah",
				"email":     "andinoferdiansah@gmail.com",
			},
			"nama_perusahaan":      "PT. Tech Solutions",
			"posisi_jabatan":       "Software Developer",
			"bidang_industri":      "Teknologi",
			"lokasi_kerja":         "Jakarta",
			"gaji_range":           "5-8 juta",
			"tanggal_mulai_kerja":  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			"status_pekerjaan":     "aktif",
			"deskripsi_pekerjaan":  "Mengembangkan aplikasi web menggunakan Go dan React",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
		bson.M{
			"alumni_info": bson.M{
				"alumni_id": alumniIDs[1], 
				"nim":       "2021002",
				"nama":      "Siti Nurhaliza",
				"email":     "siti.nurhaliza@email.com",
			},
			"nama_perusahaan":      "PT. Digital Innovation",
			"posisi_jabatan":       "System Analyst",
			"bidang_industri":      "Teknologi",
			"lokasi_kerja":         "Surabaya",
			"gaji_range":           "6-9 juta",
			"tanggal_mulai_kerja":  time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			"status_pekerjaan":     "aktif",
			"deskripsi_pekerjaan":  "Menganalisis kebutuhan sistem dan merancang solusi IT",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
		bson.M{
			"alumni_info": bson.M{
				"alumni_id": alumniIDs[2],
				"nim":       "2020001",
				"nama":      "Budi Santoso",
				"email":     "budi.santoso@email.com",
			},
			"nama_perusahaan":      "PT. Data Analytics",
			"posisi_jabatan":       "Data Scientist",
			"bidang_industri":      "Teknologi",
			"lokasi_kerja":         "Bandung",
			"gaji_range":           "8-12 juta",
			"tanggal_mulai_kerja":  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			"status_pekerjaan":     "aktif",
			"deskripsi_pekerjaan":  "Menganalisis data besar untuk insights bisnis",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
		bson.M{
			"alumni_info": bson.M{
				"alumni_id": alumniIDs[3], 
				"nim":       "2022001",
				"nama":      "Maria Garcia",
				"email":     "maria.garcia@email.com",
			},
			"nama_perusahaan":      "PT. Cloud Computing",
			"posisi_jabatan":       "DevOps Engineer",
			"bidang_industri":      "Teknologi",
			"lokasi_kerja":         "Jakarta",
			"gaji_range":           "7-10 juta",
			"tanggal_mulai_kerja":  time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			"status_pekerjaan":     "aktif",
			"deskripsi_pekerjaan":  "Mengelola infrastruktur cloud dan CI/CD pipeline",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
		bson.M{
			"alumni_info": bson.M{
				"alumni_id": alumniIDs[4], 
				"nim":       "2022002",
				"nama":      "John Smith",
				"email":     "john.smith@email.com",
			},
			"nama_perusahaan":      "PT. Mobile Apps",
			"posisi_jabatan":       "Mobile Developer",
			"bidang_industri":      "Teknologi",
			"lokasi_kerja":         "Surabaya",
			"gaji_range":           "6-9 juta",
			"tanggal_mulai_kerja":  time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			"status_pekerjaan":     "aktif",
			"deskripsi_pekerjaan":  "Mengembangkan aplikasi mobile menggunakan Flutter",
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		},
	}

	collection := db.Collection("pekerjaan_alumni")
	result, err := collection.InsertMany(ctx, pekerjaan)
	if err != nil {
		return err
	}

	log.Printf("Inserted %d pekerjaan alumni", len(result.InsertedIDs))
	return nil
}
