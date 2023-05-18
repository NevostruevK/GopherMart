package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const secretKeySize = 16

const (
	errTokenPartsCount       = "wrong count of token's parts"
	errTokenBearerPartsCount = "wrong count of Bearer token's parts"
	errTokenHeader           = "wrong token's header"
	errTokenSign             = "wrong token's sign"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Sub uint64 `json:"sub"`
}

type Token struct {
	key              []byte
	header           string
	headerPlus       string
	bearerHeaderPlus string
}

func setHeader(tk *Token) error {
	h := header{Alg: "HS256", Typ: "JWT"}
	jHeader, err := json.Marshal(&h)
	if err != nil {
		return err
	}
	tk.header = base64.URLEncoding.EncodeToString(jHeader)
	tk.headerPlus = tk.header + "."
	tk.bearerHeaderPlus = "Bearer " + tk.headerPlus
	return nil
}

func generateKey(tk *Token) error {
	tk.key = make([]byte, secretKeySize)
	rand.Seed(time.Now().UnixNano())
	_, err := rand.Read(tk.key)
	return err
}

func NewToken() (*Token, error) {
	tk := Token{}
	if err := setHeader(&tk); err != nil {
		return nil, err
	}
	if err := generateKey(&tk); err != nil {
		return nil, err
	}
	return &tk, nil
}

func (tk Token) sign(payload string) string {
	str := tk.headerPlus + payload
	sign := hmac.New(sha256.New, tk.key)
	sign.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(sign.Sum(nil))
}

func (tk Token) Make(ID uint64) (string, error) {
	p := payload{Sub: ID}
	jPayload, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}
	sPayload := base64.URLEncoding.EncodeToString(jPayload)
	return tk.bearerHeaderPlus + sPayload + "." + tk.sign(sPayload), nil
}

func (tk Token) Parse(token string) (uint64, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf(errTokenBearerPartsCount)
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf(errTokenPartsCount)
	}
	if parts[0] != tk.header {
		return 0, fmt.Errorf(errTokenHeader)
	}
	if tk.sign(parts[1]) != parts[2] {
		return 0, fmt.Errorf(errTokenSign)
	}
	bPayload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}
	p := payload{}
	err = json.Unmarshal(bPayload, &p)
	if err != nil {
		return 0, err
	}
	return p.Sub, nil
}
