package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type PersonaGovResponse struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Cuil     string `json:"cuil"`
}

func ConsultarPersonaPorCUIL(cuil string) (*PersonaGovResponse, error) {
	useMock := os.Getenv("USE_MOCK")
	if useMock == "true" {
		return consultarMock(cuil)
	}

	url := os.Getenv("SOAP_ENDPOINT")
	usuario := os.Getenv("SOAP_USER")
	password := os.Getenv("SOAP_PASSWORD")

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
	xmlns:ser="http://servicios.maestros.cba.gov.ar/">
	<soapenv:Header>
		<wsse:Security soapenv:mustUnderstand="1"
			xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
			xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
			<wsse:UsernameToken>
				<wsse:Username>%s</wsse:Username>
				<wsse:Password>%s</wsse:Password>
			</wsse:UsernameToken>
		</wsse:Security>
	</soapenv:Header>
	<soapenv:Body>
		<ser:buscarPersonaFisica>
			<personaFisicaRequest>
				<encabezado>
					<usuario>%s</usuario>
					<aplicacion>SICP</aplicacion>
					<token></token>
					<sign></sign>
				</encabezado>
				<cuil>%s</cuil>
			</personaFisicaRequest>
		</ser:buscarPersonaFisica>
	</soapenv:Body>
</soapenv:Envelope>`, usuario, password, usuario, cuil)

	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "") // Algunos servicios lo requieren vacío

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request SOAP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en respuesta SOAP: %d\nBody: %s", resp.StatusCode, string(body))
	}

	// 👇 Acá deberías parsear la respuesta XML según lo que devuelva el servicio
	// De momento te devuelvo el raw del body para debug
	fmt.Println(string(body))

	return nil, fmt.Errorf("parsing aún no implementado")
}

func consultarMock(cuil string) (*PersonaGovResponse, error) {
	return &PersonaGovResponse{
		Nombre:   "Juan",
		Apellido: "Pérez",
		Cuil:     cuil,
	}, nil
}
