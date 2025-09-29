package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	r2client "github.com/daiyanuthsa/grpc-ecom-be/pkg/r2"
)

type IStorageService interface {
	UploadProductImage(ctx context.Context, file *multipart.FileHeader) (url string, key string, err error)
	CheckIfObjectExists(ctx context.Context, key string) (bool, error)
    DeleteObject(ctx context.Context, key string) error
}

type storageService struct {
    r2Client *s3.Client // Injeksi klien R2
    validator *utils.ImageValidator // Injeksi validator
}

func (s *storageService) UploadProductImage(ctx context.Context, file *multipart.FileHeader) (string, string, error) {
    // 1. Validasi File (SRP dipenuhi oleh ImageValidator)
    if err := s.validator.Validate(file); err != nil { 
        return "", "", err 
    }

    // 2. Logika Bisnis: Tentukan Key Unik
	timestamp := time.Now().UnixNano()
    extension := filepath.Ext(file.Filename)
    objectKey := fmt.Sprintf("products/product_%d%s", timestamp, extension)
    
    // 3. Kontrol Alur: Panggil Private Method untuk Upload
    // Service hanya peduli bahwa upload berhasil.
    publicURL, err := s.putObjectToR2(ctx, objectKey, file) // 👈 Memanggil private method
    if err != nil {
        return "", "", err // Error saat upload ke storage
    }
    
    return publicURL, objectKey, nil
}


func (s *storageService) putObjectToR2(ctx context.Context, key string, fileHeader *multipart.FileHeader) (string, error) {
    
    file, err := fileHeader.Open()
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    bucketName := os.Getenv("R2_BUCKET_NAME")
    
    // Gunakan s.r2Client yang sudah di-inject
    _, err = s.r2Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
        Body:   file,
        ContentLength: aws.Int64(fileHeader.Size),
        ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
    })

    if err != nil {
        return "", err
    }
    return key, nil
}

func (s *storageService) CheckIfObjectExists(ctx context.Context, key string) (bool, error) {
    bucketName := os.Getenv("R2_BUCKET_NAME")
    
    input := &s3.HeadObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
    }

    // Gunakan klien R2 (s.r2Client) untuk memanggil HeadObject
    _, err := s.r2Client.HeadObject(ctx, input)
    
    if err != nil {
        // AWS SDK mengembalikan error tipe 'NotFound' jika objek tidak ada.
        var noSuchKey *types.NotFound
        if errors.As(err, &noSuchKey) {
            // Objek tidak ditemukan
            return false, nil 
        }
        
        // Jika ada error lain (misalnya, masalah koneksi atau permission), kembalikan error
        return false, err
    }
    
    // Jika tidak ada error, HeadObject berhasil dan objek ada
    return true, nil
}
func (s *storageService) DeleteObject(ctx context.Context, key string) error {
	bucketName := os.Getenv("R2_BUCKET_NAME")

	_, err := s.r2Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object from R2: %w", err)
	}
	return nil
}



func NewStorageService(ctx context.Context) IStorageService {
	r2, err := r2client.NewR2Client(ctx)
	if err != nil {
		panic(err) // Atau tangani error dengan lebih baik
	}
	validator := &utils.ImageValidator{
		MaxSizeBytes: 5 * 1024 * 1024, // 5 MB
		AllowedMimeTypes: map[string]string{
			"image/jpeg": ".jpeg",
			"image/png":  ".png",
			"image/gif":  ".gif",
			"image/webp": ".webp",
		},
	}
	return &storageService{r2Client: r2, validator: validator}
}