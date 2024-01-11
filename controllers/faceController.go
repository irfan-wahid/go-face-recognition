package controllers

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Kagami/go-face"
	"github.com/gofiber/fiber/v2"
)

// Path to directory with models and test images. Here it's assumed it
// points to the <https://github.com/Kagami/go-face-testdata> clone.
const dataDir = "testdata"

var (
	modelsDir = filepath.Join(dataDir, "models")
	imagesDir = filepath.Join(dataDir, "images")
)

type JSONRequest struct{
	File string `json:"image_name"`
}

func InitRecognition(c *fiber.Ctx) error {
	var req JSONRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	
	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	// Test image with 10 faces.
	testImageSample := filepath.Join(imagesDir, "twice.jpg")
	// Recognize faces on that image.
	faces, err := rec.RecognizeFile(testImageSample)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if len(faces) != 9 {
		log.Fatalf("Wrong number of faces")
	}

	// Fill known samples. In the real world you would use a lot of images
	// for each person to get better classification results but in our
	// example we just get them from one big image.
	var samples []face.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		cats = append(cats, int32(i))
	}
	// Name the categories, i.e. people on the image.
	labels := []string{
		"Tzuyu", "Jeongyeon", "Sana", "Nayeon", "Momo",
		"Mina", "Jihyo", "Dahyun", "Chaeyoung",
	}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, cats)

	// Now let's try to classify some not yet known image.
	testImageModel := filepath.Join(imagesDir, req.File)
	modelFace, err := rec.RecognizeSingleFile(testImageModel)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if modelFace == nil {
		log.Fatalf("Not a single face on the image")
	}
	catID := rec.Classify(modelFace.Descriptor)
	if catID < 0 {
		log.Fatalf("Can't classify")
	}
	
	fmt.Println(labels[catID])

	return c.SendString("Greeting processed. Check the console for the message.")
}
