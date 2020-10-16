package main

import (
	"fmt"

	"github.com/ApiCatastro/urlgoogle"
	"github.com/gofiber/fiber"
)

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/urlgoogle/:id", urlgoogle.GetUrl)

}

func main() {
	app := fiber.New()

	setupRoutes(app)
	err := app.Listen(3000)
	if err != nil {
		fmt.Printf("%s caguen reus \n", err)
	}
}
