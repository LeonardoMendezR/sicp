package services

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ErrPersonaNoEncontrada se devuelve cuando el servicio no encuentra datos de la persona
var ErrPersonaNoEncontrada = errors.New("persona no encontrada")

type Persona struct {
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	Cuil            string `json:"cuil"`
	NumeroDocumento string `json:"numero_documento"`
	TipoDocumento   string `json:"tipo_documento"`
	Sexo            string `json:"sexo"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Provincia       string `json:"provincia"`
	Localidad       string `json:"localidad"`
	Direccion       string `json:"direccion"`
	EstadoCivil     string `json:"estado_civil"`
}

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    struct {
		Response struct {
			XMLName xml.Name `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 personaFisicaResponse"`

			Nombre           string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 nombre"`
			Apellido         string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 apellido"`
			Cuil             string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 cuil"`
			NumeroDocumento  string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 numeroDocumento"`
			TipoDocumento    string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 idTipoDocumento"`
			Sexo struct {
				Tipo string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 tipoSexo"`
			} `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 sexo"`
			FechaNacimiento string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 fecNacimiento"`
			Domicilio struct {
				Provincia string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 provincia"`
				Localidad string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 localidad"`
				Calle     string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 calle"`
			} `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 domicilio"`
			EstadoCivil struct {
				Nombre string `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 nombreEstadoCivil"`
			} `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 estadoCivil"`
		} `xml:"https://cba.gov.ar/Maestros/PersonaFisica/1.0.0 personaFisicaResponse"`

		Fault *struct {
			Code   string `xml:"faultcode"`
			String string `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}




func ConsultarPersonaPorCUIL(cuil string) (*Persona, error) {
	endpoint := os.Getenv("SOAP_ENDPOINT")
	usuario := os.Getenv("SOAP_USER")
	password := os.Getenv("SOAP_PASSWORD")
	appUser := os.Getenv("SOAP_USUARIO_HEADER")
	appName := os.Getenv("SOAP_APLICACION_HEADER")
	soapAction := os.Getenv("SOAP_ACTION")
	useMock := os.Getenv("USE_MOCK")

	if useMock == "true" {
		return &Persona{
			Nombre:          "Juan",
			Apellido:        "Pérez",
			Cuil:            cuil,
			NumeroDocumento: "12345678",
			TipoDocumento:   "DNI",
			Sexo:            "MASCULINO",
			FechaNacimiento: "1990-01-01",
			Provincia:       "CORDOBA",
			Localidad:       "SAN AGUSTÍN",
			Direccion:       "CALLE FICTICIA 123",
			EstadoCivil:     "SOLTERO/A",
		}, nil
	}

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
                  xmlns:ns="https://cba.gov.ar/Maestros/PersonaFisica/1.0.0" 
                  xmlns:ns1="https://cba.gov.ar/Comunes/Encabezado/1.0.0">
  <soapenv:Header>
    <wsse:Security soapenv:mustUnderstand="1" xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">%s</wsse:Password>
      </wsse:UsernameToken>
    </wsse:Security>
  </soapenv:Header>
  <soapenv:Body>
    <ns:personaFisicaRequest>
      <ns:encabezado>
        <ns1:usuario>%s</ns1:usuario>
        <ns1:token>?</ns1:token>
        <ns1:sign>?</ns1:sign>
        <ns1:aplicacion>%s</ns1:aplicacion>
      </ns:encabezado>
      <ns:cuil>%s</ns:cuil>
    </ns:personaFisicaRequest>
  </soapenv:Body>
</soapenv:Envelope>`, usuario, password, appUser, appName, cuil)


	fmt.Println("📤 XML SOAP enviado:")
	fmt.Println(soapRequest)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", soapAction)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request SOAP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	fmt.Println("📨 Body SOAP crudo:")
	fmt.Println(string(body))

	var envelope soapEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("error parseando XML: %v", err)
	}

	if envelope.Body.Fault != nil {
		return nil, fmt.Errorf("SOAP fault: %s", envelope.Body.Fault.String)
	}

	p := envelope.Body.Response
	if p.Cuil == "" {
		return nil, ErrPersonaNoEncontrada
	}

	return &Persona{
		Nombre:          p.Nombre,
		Apellido:        p.Apellido,
		Cuil:            p.Cuil,
		NumeroDocumento: p.NumeroDocumento,
		TipoDocumento:   p.TipoDocumento,
		Sexo:            p.Sexo.Tipo,
		FechaNacimiento: p.FechaNacimiento,
		Provincia:       p.Domicilio.Provincia,
		Localidad:       p.Domicilio.Localidad,
		Direccion:       p.Domicilio.Calle,
		EstadoCivil:     p.EstadoCivil.Nombre,
	}, nil
}
