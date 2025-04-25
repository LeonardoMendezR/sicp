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

// Envelope mapea el sobre SOAP de respuesta
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		BuscarPersonaFisicaResponse struct {
			Return PersonaSOAP `xml:"return"`
		} `xml:"buscarPersonaFisicaResponse"`
	} `xml:"Body"`
}

// PersonaSOAP refleja la estructura interna del XML de respuesta
type PersonaSOAP struct {
	Nombre           string `xml:"nombre"`
	Apellido         string `xml:"apellido"`
	Cuil             string `xml:"cuil"`
	NumeroDocumento  string `xml:"numeroDocumento"`
	TipoDoc          string `xml:"idTipoDocumento"`
	Sexo             struct {
		TipoSexo string `xml:"tipoSexo"`
	} `xml:"sexo"`
	FechaNacimiento  string `xml:"fecNacimiento"`
	Domicilio        struct {
		Provincia  string `xml:"provincia"`
		Localidad  string `xml:"localidad"`
		Calle      string `xml:"calle"`
	} `xml:"domicilio"`
	EstadoCivil      struct {
		NombreEstadoCivil string `xml:"nombreEstadoCivil"`
	} `xml:"estadoCivil"`
}

// FaultEnvelope para capturar faults de negocio
type FaultEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Fault struct {
			FaultCode   string `xml:"faultcode"`
			FaultString string `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

// ConsultarPersonaPorCUIL invoca el servicio SOAP de PersonaFisica por CUIL
func ConsultarPersonaPorCUIL(cuil string) (*Persona, error) {
	endpoint := os.Getenv("SOAP_ENDPOINT")
	usuario := os.Getenv("SOAP_USER")
	password := os.Getenv("SOAP_PASSWORD")
	usuarioHeader := os.Getenv("SOAP_USUARIO_HEADER")
	aplicacionHeader := os.Getenv("SOAP_APLICACION_HEADER")
	useMock := os.Getenv("USE_MOCK")

	// Validar configuración
	if endpoint == "" || usuario == "" || password == "" || usuarioHeader == "" || aplicacionHeader == "" {
		return nil, fmt.Errorf("falta configuración en variables de entorno")
	}

	// Modo mock para desarrollo
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

	// Construir request SOAP
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:ser="http://servicios.maestros2.idecor.cba.gov.ar/"
                  xmlns:ns1="http://servicios.idecor.cba.gov.ar/"
                  xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
  <soapenv:Header>
    <wsse:Security soapenv:mustUnderstand="1">
      <wsse:UsernameToken>
        <wsse:Username>%s</wsse:Username>
        <wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">%s</wsse:Password>
      </wsse:UsernameToken>
    </wsse:Security>
    <ns1:encabezado>
      <ns1:usuario>%s</ns1:usuario>
      <ns1:aplicacion>%s</ns1:aplicacion>
    </ns1:encabezado>
  </soapenv:Header>
  <soapenv:Body>
    <ser:buscarPersonaFisica>
      <ser:cuil>%s</ser:cuil>
    </ser:buscarPersonaFisica>
  </soapenv:Body>
</soapenv:Envelope>`, usuario, password, usuarioHeader, aplicacionHeader, cuil)

	fmt.Println("📤 XML Generado:")
	fmt.Println(soapRequest)

	// Crear y enviar request HTTP
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error haciendo solicitud SOAP: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta completa
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta SOAP: %v", err)
	}
	fmt.Println("📨 Respuesta SOAP:")
	fmt.Println(string(body))

	// Manejar fault de negocio
	if bytes.Contains(body, []byte("<soap:Fault>")) {
		var fault FaultEnvelope
		if err := xml.Unmarshal(body, &fault); err == nil {
			return nil, fmt.Errorf("error SOAP: %s", fault.Body.Fault.FaultString)
		}
		return nil, ErrPersonaNoEncontrada
	}

	// Parsear respuesta exitosa
	var envelope Envelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("error parseando XML SOAP: %v", err)
	}

	p := envelope.Body.BuscarPersonaFisicaResponse.Return

	// Mapear a estructura de salida
	return &Persona{
		Nombre:          p.Nombre,
		Apellido:        p.Apellido,
		Cuil:            p.Cuil,
		NumeroDoc:       p.NumeroDocumento,
		TipoDoc:         p.TipoDoc,
		Sexo:            p.Sexo.TipoSexo,
		FechaNacimiento: p.FechaNacimiento,
		Provincia:       p.Domicilio.Provincia,
		Localidad:       p.Domicilio.Localidad,
		Direccion:       p.Domicilio.Calle,
		EstadoCivil:     p.EstadoCivil.NombreEstadoCivil,
	}, nil
}