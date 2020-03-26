package senml

// Senml Media Type Registrations
// https://tools.ietf.org/html/rfc8428#section-12.3

// Sensor Measurement Lists (SenML) Media Types
const (
	MediaTypeSenmlJSON = "application/senml+json"
	MediaTypeSenmlCBOR = "application/senml+cbor"
	MediaTypeSenmlXML  = "application/senml+xml"
	MediaTypeSenmlEXI  = "application/senml-exi"
	// Custom types
	MediaTypeCustomSenmlCSV = "text/vnd.senml.v2+csv"
)

// Sensor Streaming Measurement Lists (SenSML) Media Types
const (
	MediaTypeSensmlJSON = "application/sensml+json"
	MediaTypeSensmlCBOR = "application/sensml+cbor"
	MediaTypeSensmlXML  = "application/sensml+xml"
	MediaTypeSensmlEXI  = "application/sensml-exi"
	// Custom types
	MediaTypeCustomSensmlCSV = "text/vnd.sensml.v2+csv"
)
