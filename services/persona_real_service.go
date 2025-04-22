package services

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var ErrPersonaNoEncontrada = errors.New("persona no encontrada")

type Persona struct {
	Nombre           string `json:"nombre"`
	Apellido         string `json:"apellido"`
	Cuil             string `json:"cuil"`
	NumeroDoc        string `json:"numero_documento"`
	TipoDoc          string `json:"tipo_documento"`
	Sexo             string `json:"sexo"`
	FechaNacimiento  string `json:"fecha_nacimiento"`
	Provincia        string `json:"provincia"`
	Localidad        string `json:"localidad"`
	Direccion        string `json:"direccion"`
	EstadoCivil      string `json:"estado_civil"`
}

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Persona struct {
				Nombre          string `xml:"nombre"`
				Apellido        string `xml:"apellido"`
				Cuil            string `xml:"cuil"`
				NumeroDocumento string `xml:"numeroDocumento"`
				TipoDoc         string `xml:"idTipoDocumento"`
				Sexo            string `xml:"sexo>tipoSexo"`
				FechaNacimiento string `xml:"fechaNacimiento"`
				Provincia       string `xml:"domicilio>provincia"`
				Localidad       string `xml:"domicilio>localidad"`
				Direccion       string `xml:"domicilio>calle"`
				EstadoCivil     string `xml:"estadoCivil>nombreEstadoCivil"`
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
			Nombre:          "Juan",
			Apellido:        "Pérez",
			Cuil:            cuil,
			NumeroDoc:       "12345678",
			TipoDoc:         "DNI",
			Sexo:            "MASCULINO",
			FechaNacimiento: "1990-01-01",
			Provincia:       "CORDOBA",
			Localidad:       "SAN AGUSTÍN",
			Direccion:       "CALLE FICTICIA 123",
			EstadoCivil:     "SOLTERO/A",
		}, nil
	}

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:ns="https://cba.gov.ar/Maestros/PersonaFisica/1.0.0"
               xmlns:enc="https://cba.gov.ar/Comunes/Encabezado/1.0.0"
               xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
  <soap:Header>
    <wsse:Security soap:mustUnderstand="1">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">%s</wsse:Password>
      </wsse:UsernameToken>
    </wsse:Security>
  </soap:Header>
  <soap:Body>
    <ns:buscarPersonaFisica>
      <ns:personaFisicaRequest>
        <ns:encabezado>
          <enc:usuario>%s</enc:usuario>
          <enc:token></enc:token>
          <enc:sign></enc:sign>
          <enc:aplicacion>SICP</enc:aplicacion>
        </ns:encabezado>
        <ns:cuil>%s</ns:cuil>
      </ns:personaFisicaRequest>
    </ns:buscarPersonaFisica>
  </soap:Body>
</soap:Envelope>`, usuario, password, usuario, cuil)

	fmt.Println("\n📤 XML generado:")
	fmt.Println(soapRequest)

	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request SOAP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("\n📨 Respuesta SOAP:")
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta SOAP no exitosa (%d): %s", resp.StatusCode, string(body))
	}

	var envelope Envelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		return nil, fmt.Errorf("error parseando XML SOAP: %v", err)
	}

	p := envelope.Body.Response.Persona
	if p.Cuil == "" {
		return nil, ErrPersonaNoEncontrada
	}

	return &Persona{
		Nombre:          p.Nombre,
		Apellido:        p.Apellido,
		Cuil:            p.Cuil,
		NumeroDoc:       p.NumeroDocumento,
		TipoDoc:         p.TipoDoc,
		Sexo:            p.Sexo,
		FechaNacimiento: p.FechaNacimiento,
		Provincia:       p.Provincia,
		Localidad:       p.Localidad,
		Direccion:       p.Direccion,
		EstadoCivil:     p.EstadoCivil,
	}, nil
}
