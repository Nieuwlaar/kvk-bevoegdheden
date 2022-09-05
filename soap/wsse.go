package soap

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Implements the WS-Security standard using X.509 certificate signatures.
// https://www.di-mgt.com.au/xmldsig2.html is a handy reference to the WS-Security signing process.

const (
	wsseNS = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	wsuNS  = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	dsigNS = "http://www.w3.org/2000/09/xmldsig#"

	encTypeBinary    = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary"
	valTypeX509Token = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3"

	canonicalizationExclusiveC14N = "http://www.w3.org/2001/10/xml-exc-c14n#"
	rsaSha1Sig                    = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	sha1Sig                       = "http://www.w3.org/2000/09/xmldsig#sha1"
)

// WSSEAuthInfo contains the information required to use WS-Security X.509 signing.
type WSSEAuthInfo struct {
	certDER string
	key     *rsa.PrivateKey
}

// WSSEAuthIDs contains generated IDs used in WS-Security X.509 signing.
type WSSEAuthIDs struct {
	securityTokenID string
	bodyID          string
}

// NewWSSEAuthInfo retrieves the supplied certificate path and key path for signing SOAP requests.
// These requests will be secured using the WS-Security X.509 security standard.
// If the supplied certificate path does not point to a DER-encoded X.509 certificate, or
// if the supplied key path does not point to a PEM-encoded X.509 certificate, an error will be returned.
func NewWSSEAuthInfo(certDer string, keyContent string) (*WSSEAuthInfo, error) {
	// certFileContents, err := ioutil.ReadFile(certPath)
	// if err != nil {
	// 	return nil, err
	// }

	// certDer := string(certFileContents)

	// Super ugly way of getting the contents, but this works
	newlineRegex := regexp.MustCompile(`\r?\n`)
	certDer = newlineRegex.ReplaceAllString(certDer, "")
	certDer = strings.TrimPrefix(certDer, "-----BEGIN CERTIFICATE-----")
	certDer = strings.TrimSuffix(certDer, "-----END CERTIFICATE-----")

	// keyFileContents, err := ioutil.ReadFile(keyPath)
	// if err != nil {
	// 	return nil, err
	// }

	keyPemBlock, _ := pem.Decode([]byte(keyContent))

	if keyPemBlock == nil || keyPemBlock.Type != "RSA PRIVATE KEY" {
		return nil, ErrInvalidPEMFileSpecified
	} else if x509.IsEncryptedPEMBlock(keyPemBlock) {
		return nil, ErrEncryptedPEMFileSpecified
	}

	key, err := x509.ParsePKCS1PrivateKey(keyPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &WSSEAuthInfo{
		certDER: certDer,
		key:     key,
	}, nil
}

type binarySecurityToken struct {
	XMLName xml.Name `xml:"wsse:BinarySecurityToken"`
	XMLNS   string   `xml:"xmlns:wsu,attr"`

	WsuID string `xml:"wsu:Id,attr"`

	EncodingType string `xml:"EncodingType,attr"`
	ValueType    string `xml:"ValueType,attr"`

	Value string `xml:",chardata"`
}

type canonicalizationMethod struct {
	XMLName   xml.Name `xml:"CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type signatureMethod struct {
	XMLName   xml.Name `xml:"SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type digestMethod struct {
	XMLName   xml.Name `xml:"DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type digestValue struct {
	XMLName xml.Name `xml:"DigestValue"`
	Value   string   `xml:",chardata"`
}

type transform struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type transforms struct {
	XMLName   xml.Name `xml:"Transforms"`
	Transform transform
}

type signatureReference struct {
	XMLName xml.Name `xml:"Reference"`
	URI     string   `xml:"URI,attr"`

	Transforms transforms

	DigestMethod digestMethod
	DigestValue  digestValue
}

type signedInfo struct {
	XMLName xml.Name `xml:"SignedInfo"`
	XMLNS   string   `xml:"xmlns,attr"`

	CanonicalizationMethod canonicalizationMethod
	SignatureMethod        signatureMethod
	References             []signatureReference
}

type strReference struct {
	XMLName   xml.Name `xml:"wsse:Reference"`
	ValueType string   `xml:"ValueType,attr"`
	URI       string   `xml:"URI,attr"`
}

type securityTokenReference struct {
	XMLName xml.Name `xml:"wsse:SecurityTokenReference"`
	XMLNS   string   `xml:"xmlns:wsu,attr"`

	Reference strReference
}

type keyInfo struct {
	XMLName xml.Name `xml:"KeyInfo"`

	SecurityTokenReference securityTokenReference
}

type signature struct {
	XMLName xml.Name `xml:"Signature"`
	XMLNS   string   `xml:"xmlns,attr"`

	SignedInfo     signedInfo
	SignatureValue string `xml:"SignatureValue"`
	KeyInfo        keyInfo
}

type security struct {
	XMLName        xml.Name `xml:"wsse:Security"`
	XMLNS          string   `xml:"xmlns:wsse,attr"`
	MustUnderstand string   `xml:"http://schemas.xmlsoap.org/soap/envelope/ mustUnderstand,attr"`

	BinarySecurityToken binarySecurityToken
	Signature           signature
	Timestamp           timestamp
}

type timestamp struct {
	XMLName xml.Name  `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd Timestamp"`
	ID      string    `xml:"Id,attr,omitempty"`
	Created time.Time `xml:"Created,omitempty"`
	Expires time.Time `xml:"Expires,omitempty"`
}

func (w *WSSEAuthIDs) generateToken() ([]byte, error) {
	// We use a concatentation of the time and 10 securely generated random numbers to be the tokens.
	b := make([]byte, 10)

	token := sha1.New()
	token.Write([]byte(time.Now().Format(time.RFC3339)))

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token.Write(b)
	tokenHex := token.Sum(nil)

	return tokenHex, nil
}

func generateWSSEAuthIDs() (*WSSEAuthIDs, error) {
	w := &WSSEAuthIDs{}

	securityTokenHex, err := w.generateToken()
	if err != nil {
		return nil, err
	}

	w.securityTokenID = fmt.Sprintf("SecurityToken-%x", securityTokenHex)

	bodyTokenHex, err := w.generateToken()
	if err != nil {
		return nil, err
	}

	w.bodyID = fmt.Sprintf("Body-%x", bodyTokenHex)
	return w, nil
}

func generateRef(ID string, data any, rootEl string) signatureReference {
	refEnc, _ := xml.Marshal(data)

	canonTimestampEnc, _ := canonicalize(refEnc, rootEl)

	refHasher := sha1.New()
	refHasher.Write(canonTimestampEnc)
	encodedTimestampDigest := base64.StdEncoding.EncodeToString(refHasher.Sum(nil))
	ref := signatureReference{URI: "#" + ID,
		Transforms: transforms{
			Transform: transform{
				Algorithm: canonicalizationExclusiveC14N,
			},
		},
		DigestMethod: digestMethod{
			Algorithm: sha1Sig,
		},
		DigestValue: digestValue{
			Value: encodedTimestampDigest,
		},
	}
	return ref
}

func (w *WSSEAuthInfo) sign(body Body, header Header, ids *WSSEAuthIDs) (security, error) {
	// 0. We create the body_id and security_token_id values
	body.ID = "_0"

	ref0 := generateRef(body.ID, body, "Body")

	ts := timestamp{
		ID:      "_1",
		Created: time.Now().UTC(),
		Expires: time.Now().UTC().Add(10 * time.Minute),
	}

	ref1 := generateRef(ts.ID, ts, "Timestamp")

	action := header.Headers[0]
	ref2 := generateRef(action.(ActionHeader).ID, action, "Action")

	messageID := header.Headers[1]
	ref3 := generateRef(messageID.(MessageIDHeader).ID, messageID, "MessageID")

	to := header.Headers[2]
	ref4 := generateRef(to.(ToHeader).ID, to, "To")

	// if err != nil {
	// 	return security{}, err
	// }

	// 2. Set the DigestValue then sign the 'SignedInfo' struct
	signedInfo := signedInfo{
		XMLNS: dsigNS,
		CanonicalizationMethod: canonicalizationMethod{
			Algorithm: canonicalizationExclusiveC14N,
		},
		SignatureMethod: signatureMethod{
			Algorithm: rsaSha1Sig,
		},
		References: []signatureReference{
			ref0,
			ref1,
			ref2,
			ref3,
			ref4,
		},
	}

	signedInfoEnc, err := xml.Marshal(signedInfo)
	if err != nil {
		return security{}, err
	}

	signedInfoHasher := sha1.New()
	signedInfoHasher.Write(signedInfoEnc)
	signedInfoDigest := signedInfoHasher.Sum(nil)

	signatureValue, err := rsa.SignPKCS1v15(rand.Reader, w.key, crypto.SHA1, signedInfoDigest)
	if err != nil {
		return security{}, err
	}

	encodedSignatureValue := base64.StdEncoding.EncodeToString(signatureValue)

	secHeader := security{
		XMLNS:          wsseNS,
		MustUnderstand: "1",
		BinarySecurityToken: binarySecurityToken{
			XMLNS:        wsuNS,
			WsuID:        ids.securityTokenID,
			EncodingType: encTypeBinary,
			ValueType:    valTypeX509Token,
			Value:        w.certDER,
		},
		Timestamp: ts,
		Signature: signature{
			XMLNS:          dsigNS,
			SignedInfo:     signedInfo,
			SignatureValue: encodedSignatureValue,
			KeyInfo: keyInfo{
				SecurityTokenReference: securityTokenReference{
					XMLNS: wsuNS,
					Reference: strReference{
						ValueType: valTypeX509Token,
						URI:       "#" + ids.securityTokenID,
					},
				},
			},
		},
	}

	return secHeader, nil
}
