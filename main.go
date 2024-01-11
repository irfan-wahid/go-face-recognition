package main

import (
	"go-face-recognition-tutorial/controllers"
	"github.com/gofiber/fiber/v2"
)

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
func main() {
	app := fiber.New()

	app.Post("/", controllers.InitRecognition)

	app.Listen(":3000")
}
