package urlgoogle

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber"
	"golang.org/x/net/html"
)

func GetUrl(c *fiber.Ctx) {
	id := c.Params("id")
	varHtml := httpExampleGetJson(id)
	StrVarHtml := BytesToString(varHtml)

	//varHtml es un []byte y hay que transformarlos a string usamos una funcioncita:
	doc, _ := html.Parse(strings.NewReader(StrVarHtml))

	urlGoogle, err := buscaCoordenadas(doc)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("abriendo la web de google....")
	openBrowser(urlGoogle)

	err = c.JSON(urlGoogle)
	if err != nil {
		log.Fatal(err)
	}
}
