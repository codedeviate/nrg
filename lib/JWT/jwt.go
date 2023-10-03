package JWT

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type JWTToken struct {
	Header    JWTTokenHeader
	Payload   JWTTokenPayload
	Signature string
}

type JWTTokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTTokenPayload struct {
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Nbf int64  `json:"nbf"`
	Jti string `json:"jti"`
}

func UnpackJWTToken(token string, mode int) (string, error) {
	// Check if the token is empty
	if len(token) == 0 {
		// Check if the token is being piped in
		fi, statErr := os.Stdin.Stat()
		if statErr == nil && (fi.Mode()&os.ModeCharDevice == 0) {
			stdin, inErr := io.ReadAll(os.Stdin)
			if inErr == nil && len(stdin) > 0 {
				token = string(stdin)
			}
		}
	}
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return "", errors.New("Invalid JWT token")
	}
	// Decode the first part of the token
	decodedFirst, err1 := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if err1 != nil {
		return "", err1
	}
	// Decode the second part of the token
	decodedSecond, err2 := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err2 != nil {
		return "", err2
	}
	// Print the decoded parts
	// Prettify the JSON
	var data interface{}
	json.Unmarshal(decodedFirst, &data)
	var pretty1 []byte
	var err3 error
	err3 = nil
	if mode&1 == 1 {
		pretty1, err3 = json.MarshalIndent(data, "", "    ")
	} else if mode&2 == 2 {
		pretty1, err3 = json.Marshal(data)
	}
	if err3 != nil {
		return "", err3
	}
	var result string
	if mode&1 == 1 {
		result = "Header:\n"
		result += string(pretty1)
	} else if mode&2 == 2 {
		result = " Header: "
		result += string(pretty1)
	}

	json.Unmarshal(decodedSecond, &data)
	var pretty2 []byte
	var err4 error
	err4 = nil
	if mode&1 == 1 {
		pretty2, err4 = json.MarshalIndent(data, "", "    ")
	} else if mode&2 == 2 || mode&4 == 4 {
		pretty2, err4 = json.Marshal(data)
	}
	if err4 != nil {
		return "", err4
	}
	if mode&1 == 1 {
		result += "\nPayload:\n"
	} else if mode&2 == 2 {
		result += " Payload: "
	} else if mode&4 == 4 {
		result = ""
	}
	result += string(pretty2)
	if mode&1 == 1 || mode&2 == 2 {
		payload := JWTTokenPayload{}
		json.Unmarshal(decodedSecond, &payload)
		result += "\n"

		if payload.Iat != 0 {
			result += fmt.Sprintf("Issued at: %s\n", time.Unix(payload.Iat, 0).Format(time.RFC3339))
		}
		if payload.Nbf != 0 {
			result += fmt.Sprintf("Not before: %s", time.Unix(payload.Nbf, 0).Format(time.RFC3339))
			if payload.Nbf > time.Now().Unix() {
				result += " (not yet valid)\n"
			} else {
				result += "\n"
			}
		}
		if payload.Exp != 0 {
			result += fmt.Sprintf("Expires at: %s", time.Unix(payload.Exp, 0).Format(time.RFC3339))
			if payload.Exp < time.Now().Unix() {
				result += " (expired)\n"
			} else {
				result += "\n"
			}
		}
		if payload.Iss != "" {
			result += fmt.Sprintf("Issuer: %s\n", payload.Iss)
		}
		if payload.Aud != "" {
			result += fmt.Sprintf("Audience: %s\n", payload.Aud)
		}
		if payload.Sub != "" {
			result += fmt.Sprintf("Subject: %s\n", payload.Sub)
		}
		if payload.Jti != "" {
			result += fmt.Sprintf("JWT ID: %s\n", payload.Jti)
		}
	}
	return result, nil
}

func SignJWTToken(secret string, data []byte) (string, error) {
	data = []byte(strings.ReplaceAll(string(data), "{{NOW}}", fmt.Sprintf("%d", time.Now().Unix())))
	data = []byte(strings.ReplaceAll(string(data), "{{NOW+1H}}", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix())))
	data = []byte(strings.ReplaceAll(string(data), "{{NOW+1D}}", fmt.Sprintf("%d", time.Now().Add(time.Hour*24).Unix())))
	data = []byte(strings.ReplaceAll(string(data), "{{NOW+1W}}", fmt.Sprintf("%d", time.Now().Add(time.Hour*24*7).Unix())))
	data = []byte(strings.ReplaceAll(string(data), "{{NOW+1M}}", fmt.Sprintf("%d", time.Now().Add(time.Hour*24*30).Unix())))
	data = []byte(strings.ReplaceAll(string(data), "{{NOW+1Y}}", fmt.Sprintf("%d", time.Now().Add(time.Hour*24*365).Unix())))
	// Build the header
	header := base64.URLEncoding.EncodeToString([]byte("{\"alg\":\"HS256\",\"typ\":\"JWT\"}"))
	header = strings.ReplaceAll(header, "=", "")
	// Build the payload
	payload := base64.URLEncoding.EncodeToString(data)
	payload = strings.ReplaceAll(payload, "=", "")
	// Sign the token
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(header + "." + payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	signature = strings.ReplaceAll(signature, "=", "")
	// Print the token
	return header + "." + payload + "." + signature, nil
}

func ValidateJWTToken(token string, secret string) error {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return errors.New("Invalid JWT token")
	}

	header := JWTTokenHeader{}
	decodedHeader, errHeader := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if errHeader != nil {
		return errHeader
	}
	json.Unmarshal(decodedHeader, &header)
	if header.Alg != "HS256" {
		return errors.New("JWT token algorithm not supported")
	}
	if header.Typ != "JWT" {
		return errors.New("JWT token type not supported")
	}

	payload := JWTTokenPayload{}
	decodedPayload, errPayload := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if errPayload != nil {
		return errPayload
	}
	json.Unmarshal(decodedPayload, &payload)
	if payload.Exp > 0 && payload.Exp < time.Now().Unix() {
		return errors.New("JWT token expired")
	}
	if payload.Nbf > time.Now().Unix() {
		return errors.New("JWT token not valid yet")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tokenParts[0] + "." + tokenParts[1]))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	signature = strings.ReplaceAll(signature, "=", "")
	if signature != tokenParts[2] {
		return errors.New("Invalid JWT token")
	}
	return nil
}
