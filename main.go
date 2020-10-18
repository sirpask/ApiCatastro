package main

import (
	"fmt"

	_ "github.com/ApiCatastro/docs"
	"github.com/ApiCatastro/urlgoogle"
	swagger "github.com/arsmn/fiber-swagger"
	"github.com/gofiber/fiber"
)

//@title Fiber Catastro API
//@version 1.0
//@description metes un numero de catastro y devuelve la pagina del googlemaps
//@termsOfService http://swagger.io/terms/
//@contact.name sirpask
//@contact.email el_pask@hotmail.com
//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html
//@host localhost:8080
//@BasePath /

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/urlgoogle/:id", urlgoogle.GetUrl)
	//http://localhost:8080/api/v1/urlgoogle/42071A036260920000MO

}

func main() {
	app := fiber.New()

	// setupRoutes godoc
	// @Summary transforma un numero catastral en una pagina de googlemaps
	// @Description get URL by NC
	// @ID get-string-by-int
	// @Accept  json
	// @Produce  json // string
	// @Param id path string true "Numero referencia catastral"
	// @Success 200 {object} url
	// @Failure 400 {object} HTTPError
	// @Failure 404 {object} HTTPError
	// @Failure 500 {object} HTTPError
	// @Router /api/v1/urlgoogle/{id} [get]

	setupRoutes(app)

	app.Use("/swagger", swagger.Handler) // default
	err := app.Listen(":8080")
	//err := app.Listen(3000)
	if err != nil {
		fmt.Printf("%s caguen reus \n", err)
	}
}
