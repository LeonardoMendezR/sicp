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

// Persona representa los datos que vamos a devolver al cliente
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
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Persona personaData `xml:"Body>personaFisicaResponse"`
		Fault   *soapFault  `xml:"Body>Fault"`
	} `xml:"Envelope"`
}

type personaData struct {
	Nombre          string `xml:"nombre"`
	Apellido        string `xml:"apellido"`
	Cuil            string `xml:"cuil"`
	NumeroDocumento string `xml:"numeroDocumento"`
	TipoDocumento   string `xml:"idTipoDocumento"`
	Sexo struct {
		TipoSexo string `xml:"tipoSexo"`
	} `xml:"sexo"`
	FechaNacimiento string `xml:"fecNacimiento"`
	Domicilio struct {
		Provincia string `xml:"provincia"`
		Localidad string `xml:"localidad"`
		Calle     string `xml:"calle"`
	} `xml:"domicilio"`
	EstadoCivil struct {
		NombreEstadoCivil string `xml:"nombreEstadoCivil"`
	} `xml:"estadoCivil"`
}

type soapFault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
}

// ConsultarPersonaPorCUIL invoca el servicio SOAP de PersonaFisica por CUIL
func ConsultarPersonaPorCUIL(cuil string) (*Persona, error) {
	endpoint := os.Getenv("SOAP_ENDPOINT")
	usuario := os.Getenv("SOAP_USER")
	password := os.Getenv("SOAP_PASSWORD")
	appUser := os.Getenv("SOAP_USUARIO_HEADER")
	appName := os.Getenv("SOAP_APLICACION_HEADER")
	soapAction := os.Getenv("SOAP_ACTION")
	useMock := os.Getenv("USE_MOCK")

	if endpoint == "" || usuario == "" || password == "" || appUser == "" || appName == "" || soapAction == "" {
		return nil, errors.New("falta configuración en variables de entorno")
	}

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
			Localidad:       "CIUDAD",
			Direccion:       "CALLE FICTICIA 123",
			EstadoCivil:     "SOLTERO/A",
		}, nil
	}

	s := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:ns="https://cba.gov.ar/Maestros/PersonaFisica/1.0.0"
               xmlns:ns1="https://cba.gov.ar/Comunes/Encabezado/1.0.0"
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
    <ns:personaFisicaRequest>
      <ns:encabezado>
        <ns1:usuario>%s</ns1:usuario>
        <ns1:token></ns1:token>
        <ns1:sign></ns1:sign>
        <ns1:aplicacion>%s</ns1:aplicacion>
      </ns:encabezado>
      <ns:cuil>%s</ns:cuil>
    </ns:personaFisicaRequest>
  </soap:Body>
</soap:Envelope>`, usuario, password, appUser, appName, cuil)

	fmt.Println("📤 XML SOAP enviado:")
	fmt.Println(s)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(s))
	if err != nil {
		return nil, fmt.Errorf("error creando request SOAP: %v", err)
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", soapAction)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando request SOAP: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("🛰️ Status HTTP:", resp.StatusCode)
	fmt.Println("📩 Headers de respuesta:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta SOAP: %v", err)
	}
	fmt.Println("📨 Body SOAP crudo:")
	fmt.Println(string(body))

	var env soapEnvelope
	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("error parseando SOAP: %v", err)
	}

	if env.Body.Fault != nil {
		return nil, ErrPersonaNoEncontrada
	}

	p := env.Body.Persona
	if p.Cuil == "" {
		return nil, ErrPersonaNoEncontrada
	}

	return &Persona{
		Nombre:          p.Nombre,
		Apellido:        p.Apellido,
		Cuil:            p.Cuil,
		NumeroDocumento: p.NumeroDocumento,
		TipoDocumento:   p.TipoDocumento,
		Sexo:            p.Sexo.TipoSexo,
		FechaNacimiento: p.FechaNacimiento,
		Provincia:       p.Domicilio.Provincia,
		Localidad:       p.Domicilio.Localidad,
		Direccion:       p.Domicilio.Calle,
		EstadoCivil:     p.EstadoCivil.NombreEstadoCivil,
	}, nil
}
