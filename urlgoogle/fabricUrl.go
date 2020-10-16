package urlgoogle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Primera funcion para formar la url que ataca el api que nos va a dar los datos de parcela y municipio.
func httpExampleGetJson(refCatastral string) (bodyHtml []byte) {
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ")

	rootUrl := "http://ovc.catastro.meh.es/ovcservweb/OVCSWLocalizacionRC/OVCCallejero.asmx/Consulta_DNPRC?Provincia=&Municipio=&RC="

	urlRefCatas := strings.Join([]string{rootUrl, refCatastral}, "")

	//	get http example
	resp, err := http.Get(urlRefCatas)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//print json body
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyHtml = contents
	//fmt.Println(string(contents))
	return bodyHtml
}

//funcion que transforma los bytes en string
func BytesToString(data []byte) string {
	return string(data[:])
}

//funcion que selecciona el elemento que vamos a separar del formato html, va leyendo linea  linea y devuelve el primero en formato html
func Body(doc *html.Node, etiqueta string) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == etiqueta {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <etiqueta> in the node tree")
}

//coge el formato html que hemos seleccionado, y lo revisa y saca el valor.
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)
	if err != nil {
		// handle your error here
		log.Fatal(err)
	}
	domDocTest := html.NewTokenizer(strings.NewReader(buf.String()))
	for tokenType := domDocTest.Next(); tokenType != html.ErrorToken; {
		if tokenType != html.TextToken {
			tokenType = domDocTest.Next()
			continue
		}
		TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))

		return TxtContent
	}
	return buf.String()

}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func buscaCoordenadas(doc *html.Node) (string, error) {

	//https://geekflare.com/es/json-online-tools/    19 herramientas para json on line
	//provincia
	bn, err := Body(doc, "cp")
	if err != nil {
		return "0", errors.New("provincia no encontrada")
	}
	provincia := renderNode(bn)

	//municipio
	bn, err = Body(doc, "cmc")
	if err != nil {
		return "0", errors.New("municipio no encontrado")
	}
	municipio := renderNode(bn)

	//municipio agregado cma
	bn, err = Body(doc, "cma")
	var agregado string
	if err != nil {
		//return 0, 0, errors.New("agregado no encontrado")
		agregado = "0"
	} else {
		agregado = renderNode(bn)
	}

	//zona concentracion czc
	bn, err = Body(doc, "czc")
	if err != nil {
		return "0", errors.New("zona no encontrada")
	}
	zona := renderNode(bn)

	//poligono cpo
	bn, err = Body(doc, "cpo")
	if err != nil {
		return "0", errors.New("poligono no encontrado")
	}
	poligono := renderNode(bn)

	//parcela cpa
	bn, err = Body(doc, "cpa")
	if err != nil {
		return "0", errors.New("parcela no encontrada")
	}
	parcela := renderNode(bn)

	rootUrl := "http://sigpac.mapa.es/fega/ServiciosVisorSigpac/query/parcelabox/"

	urlRefCatas := strings.Join([]string{rootUrl, provincia, "/", municipio, "/", agregado, "/", zona, "/", poligono, "/", parcela, ".geojson"}, "")

	//	get http example

	resp, err := http.Get(urlRefCatas)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//print json body
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type NumeroCoordenada struct {
		X1 float64 `json:"x1"`
		Y1 float64 `json:"y1"`
		X2 float64 `json:"x2"`
		Y2 float64 `json:"Y2"`
	}

	type Caracteristicas struct {
		Geometry   string                 `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
		Tipo       string                 `json:"type"`
	}

	type Zero struct {
		Tabla Caracteristicas `json:"0"`
	}

	type Coordenada struct {
		Tipo     string        `json:"type"`
		Features []interface{} `json:"features"`
		Crs      string        `json:"crs"`
	}

	dec := json.NewDecoder(strings.NewReader(string(contents)))
	if dec != nil {
		var m Coordenada
		if err := dec.Decode(&m); err == io.EOF {
			log.Fatal(err)
		} else if err != nil {
			log.Fatal(err)
		}

		//converting a []interface{} to a []string
		z := make([]string, len(m.Features))
		for i, v := range m.Features {
			z[i] = fmt.Sprint(v)
		}

		posX1 := strings.LastIndex(z[0], "x1:")

		posX2 := strings.LastIndex(z[0], " x2:")

		posY1 := strings.LastIndex(z[0], " y1:")

		posY2 := strings.LastIndex(z[0], " y2:")

		posFin := strings.LastIndex(z[0], "] ")

		var coordenada1_inicio = posX1 + 3
		var coordenada1_fin = posX2

		var coordenada2_inicio = posX2 + 4
		var coordenada2_fin = posY1

		var coordenada3_inicio = posY1 + 4
		var coordenada3_fin = posY2

		var coordenada4_inicio = posY2 + 4
		var coordenada4_fin = posFin

		var sCoordenada1 = z[0][coordenada1_inicio:coordenada1_fin]

		var sCoordenada2 = z[0][coordenada2_inicio:coordenada2_fin]

		var sCoordenada3 = z[0][coordenada3_inicio:coordenada3_fin]

		var sCoordenada4 = z[0][coordenada4_inicio:coordenada4_fin]

		Coordenada1, _ := strconv.ParseFloat(sCoordenada1, 64)

		Coordenada2, _ := strconv.ParseFloat(sCoordenada2, 64)

		Coordenada3, _ := strconv.ParseFloat(sCoordenada3, 64)

		Coordenada4, _ := strconv.ParseFloat(sCoordenada4, 64)

		cordX_sum := ((Coordenada1 + Coordenada2) / 2)

		cordY_sum := ((Coordenada3 + Coordenada4) / 2)

		sCordX := fmt.Sprintf("%.15f", cordX_sum)

		sCordY := fmt.Sprintf("%.15f", cordY_sum)

		GoogleUrl := "http://maps.google.com/maps?q=loc:"

		urlRefGoogle := strings.Join([]string{GoogleUrl, sCordY, ",", sCordX}, "")
		return urlRefGoogle, nil

	}
	return "=)", nil
}
