package ali

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

func urlValues(c *Client, param PayParam) url.Values {
	values := url.Values{}
	values.Add("app_id", c.config.AppID)
	values.Add("method", param.URI())
	values.Add("format", "JSON")
	values.Add("charset", "utf-8")
	values.Add("sign_type", c.config.SignType)
	values.Add("timestamp", generateTimestampStr())
	values.Add("version", "1.0")
	values.Add("biz_content", param.BizContent())

	for k, v := range param.ExtraParams() {
		values.Add(k, v)
	}

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values.Add("sign", signature(keys, values, c.config.AppPrivateKey, c.config.SignType))
	return values
}

func signature(keys []string, values url.Values, privateKey []byte, signType string) string {
	if values == nil {
		values = url.Values{}
	}

	var valueList []string
	for _, k := range keys {
		v := strings.TrimSpace(values.Get(k))
		if v != "" {
			valueList = append(valueList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	concat := strings.Join(valueList, "&")

	var sign string
	if signType == "RSA" {
		sign = signPKCS1v15([]byte(concat), privateKey, crypto.SHA1)
	} else if signType == "RSA2" {
		sign = signPKCS1v15([]byte(concat), privateKey, crypto.SHA256)
	}
	return sign
}

func signPKCS1v15(source, privateKey []byte, hash crypto.Hash) string {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		fmt.Println("BLOCK", block, len(privateKey))
		return ""
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// rsaPrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("RSAPRIVATE", err)
		return ""
	}

	h := hash.New()
	h.Write(source)
	hashed := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, hash, hashed)
	// s, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey.(*rsa.PrivateKey), hash, hashed)
	if err != nil {
		fmt.Println("SIGNPKCS", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(s)
}

func verify(values url.Values, publicKey []byte, signType string) bool {
	var excluded []string
	for k := range values {
		if k == "sign" || k == "sign_type" {
			continue
		}
		excluded = append(excluded, k)
	}
	sort.Strings(excluded)

	var valueList []string
	for _, k := range excluded {
		v := values.Get(k)
		if v != "" {
			valueList = append(valueList, fmt.Sprintf("%s=%s", k, v))
		}
	}
	concat := strings.Join(valueList, "&")

	decoded, err := base64.StdEncoding.DecodeString(values.Get("sign"))
	if err != nil {
		return false
	}

	var ok bool
	if signType == "RSA" {
		ok = verifyPKCS1v15([]byte(concat), decoded, publicKey, crypto.SHA1)
	} else if signType == "RSA2" {
		ok = verifyPKCS1v15([]byte(concat), decoded, publicKey, crypto.SHA256)
	}
	return ok
}

func verifyPKCS1v15(source, sign, publicKey []byte, hash crypto.Hash) bool {
	h := hash.New()
	h.Write(source)
	hashed := h.Sum(nil)

	block, _ := pem.Decode(publicKey)
	if block == nil {
		fmt.Println("VERIFY BLOCK", block)
		return false
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKIXPublicKey", err)
		return false
	}

	rsaPublicKey := pub.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(rsaPublicKey, hash, hashed, sign)
	if err != nil {
		fmt.Println("VerifyPKCS1v15", err)
		return false
	}
	return true
}

func generateTimestampStr() string {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func marshalJSON(val interface{}) string {
	data, err := json.Marshal(val)
	if err != nil {
		return ""
	}
	return string(data)
}

func toValues(st interface{}) (url.Values, error) {
	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("need a struct type, got %T", st)
	}

	typ := val.Type()
	result := url.Values{}

	for i := 0; i < val.NumField(); i++ {
		sf := typ.Field(i)
		if tag, ok := sf.Tag.Lookup("json"); ok && tag != "" {
			result.Add(tag, val.Field(i).String())
		}
	}
	return result, nil
}
