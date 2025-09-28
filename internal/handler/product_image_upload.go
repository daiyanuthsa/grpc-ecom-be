package handler

import (
	"context"
	"net/http"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/service"
	"github.com/gofiber/fiber/v2"
)



func UploadProductImageHandler(c *fiber.Ctx) error {
    ctx := context.Context(c.Context()) // Gunakan context Fiber
    
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get image file",
		})
	}
	storageService := service.NewStorageService(ctx)

	publicURL, objectKey, serviceErr := storageService.UploadProductImage(c.Context(), file)
	if serviceErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": serviceErr.Error(),
		})
	}

    // 3. Kembalikan respons sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image uploaded successfully",
		"url": publicURL,
        "key": objectKey, // Ini yang akan disimpan di DB gRPC Anda
	})
}

