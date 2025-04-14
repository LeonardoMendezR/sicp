package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Persona struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Cuil     string `json:"cuil"`
}
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Persona struct {
				Nombre   string `xml:"nombre"`
				Apellido string `xml:"apellido"`
				Cuil     string `xml:"cuil"`
			} `xml:"personaFisicaResponse"`
		} `xml:"buscarPersonaFisicaResponse"`
	} `xml:"Body"`
}
func ConsultarPersonaPorCUIL(cuil string) (*Persona, error) {
	url := os.Getenv("SOAP_ENDPOINT")
	usuario := os.Getenv("SOAP_USER")
	password := os.Getenv("SOAP_PASSWORD")
	useMock := os.Getenv("USE_MOCK")

	if useMock == "true" {
		return &Persona{
			Nombre:   "Juan",
			Apellido: "Pérez",
			Cuil:     cuil,
		}, nil
	}
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:ser="https://cba.gov.ar/Maestros/PersonaFisica/1.0.0"
                  xmlns:enc="https://cba.gov.ar/Comunes/Encabezado/1.0.0"
                  xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
  <soapenv:Header>
    <wsse:Security soapenv:mustUnderstand="1">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">%s</wsse:Password>
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

	fmt.Println("XML generado:")
	fmt.Println(soapRequest)

	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "PersonaFisica_buscarPersonaFisica")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request SOAP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta SOAP no exitosa (%d): %s", resp.StatusCode, string(body))
	}

	var envelope Envelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		return nil, fmt.Errorf("error parseando XML SOAP: %v", err)
	}

	persona := envelope.Body.Response.Persona
	return &Persona{
		Nombre:   persona.Nombre,
		Apellido: persona.Apellido,
		Cuil:     persona.Cuil,
	}, nil
}
