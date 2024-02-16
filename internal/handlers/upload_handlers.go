package handlers

import (
	"Proyectos-UTEQ/api-ortografia/internal/services"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UploadHandler struct {
	config *viper.Viper
}

func NewUploadHandler(config *viper.Viper) *UploadHandler {
	return &UploadHandler{
		config: config,
	}
}

// UploadFiles para subir archivos
func (h *UploadHandler) UploadFiles(c *fiber.Ctx) error {

	filesPath := make([]string, 0)

	if form, err := c.MultipartForm(); err == nil {
		pathString := "/"

		if pathValues := form.Value["path"]; len(pathValues) > 0 {
			if len(pathValues[0]) > 0 {
				pathString = pathValues[0]
			}
		}

		files := form.File["file"]

		for _, file := range files {

			// Se crea la ruta donde se guardaran los archivos
			pathString = path.Join("./uploads", pathString)
			controllingFolders(pathString)
			// sacar la extencio del archivo archiv.png > .png
			extension := path.Ext(file.Filename)
			newuuid := uuid.NewString()

			pathString = path.Join(pathString, newuuid+extension)

			if err := c.SaveFile(file, pathString); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": err.Error(),
				})
			}

			completePaht := fmt.Sprintf("%s/api/%s", h.config.GetString("APP_HOST"), pathString)

			filesPath = append(filesPath, completePaht)
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al recuperar los archivos",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Files uploaded successfully",
		"data":    filesPath,
	})
}

// controllingFolders controla la creacion de carpetas para subir los archivos.
func controllingFolders(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error al crear el directorio")
		}

	}
	return nil
}

func (h *UploadHandler) UploadFileToGoogle(c *fiber.Ctx) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer client.Close()

	objectName := ""

	if form, err := c.MultipartForm(); err == nil {
		if form.File["file"] == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		fileHeader := form.File["file"][0]
		file, _ := fileHeader.Open()
		ext := path.Ext(fileHeader.Filename)
		//name := fileHeader.Filename[0 : len(fileHeader.Filename)-len(ext)]

		objectName = fmt.Sprintf("%s_%s%s", time.Now().Format("2006-01-02"), uuid.NewString(), ext)
		wc := client.Bucket(services.NAME_BUCKET).Object(objectName).NewWriter(ctx)
		defer wc.Close()

		if _, err := io.Copy(wc, file); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Files uploaded successfully",
		"url":     fmt.Sprintf("https://storage.googleapis.com/%s/%s", services.NAME_BUCKET, objectName),
	})
}
