package service

import (
	"context"
	"ecommerce/internal/config"
	"io"
	"log/slog"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudService interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
}

var _ CloudService = (*CloudinaryService)(nil)

type CloudinaryService struct {
	Client *cloudinary.Cloudinary
	Logger *slog.Logger
}

func NewCloudinaryService(cfg *config.Config, logger *slog.Logger) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryURL)
	if err != nil {
		logger.Error("Failed to initialize Cloudinary client", "error", err)
		return nil, err
	}

	logger.Info("Cloudinary service initialized successfully")
	return &CloudinaryService{
		Client: cld,
		Logger: logger,
	}, nil
}

func (s *CloudinaryService) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	s.Logger.Info("Uploading image to Cloudinary", "filename", filename)

	uploadParams := uploader.UploadParams{
		Folder: "ecommerce_products",

		// PublicID: filename,
	}

	result, err := s.Client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		s.Logger.Error("Cloudinary upload failed", "error", err)
		return "", err
	}

	s.Logger.Info("Image uploaded successfully", "url", result.SecureURL)
	return result.SecureURL, nil
}
