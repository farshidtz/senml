package codec

import (
	"encoding/xml"

	"github.com/farshidtz/senml/v2"
)

type xmlPack struct {
	senml.Pack `xml:"senml"`
	XMLName    *bool  `xml:"sensml"`
	XMLNS      string `xml:"xmlns,attr"`
}


// EncodeXML serializes the SenML pack into XML bytes
func EncodeXML(p senml.Pack, pretty bool) ([]byte, error) {

	xmlPack := xmlPack{
		Pack:  p,
		XMLNS: "urn:ietf:params:xml:ns:senml",
	}

	if pretty {
		return xml.MarshalIndent(&xmlPack, "", "  ")
	}

	return xml.Marshal(&xmlPack)
}

// DecodeXML takes a SenML pack in XML bytes and decodes it into a Pack
func DecodeXML(b []byte) (senml.Pack, error) {
	var temp xmlPack
	err := xml.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	return temp.Pack, nil
}