package wx

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func toMap(st interface{}) (map[string]string, error) {
	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("need a struct type, got %T", st)
	}

	typ := val.Type()
	result := map[string]string{}

	for i := 0; i < val.NumField(); i++ {
		sf := typ.Field(i)
		if tag, ok := sf.Tag.Lookup("xml"); ok && tag != "" && tag != "xml" {
			result[tag] = val.Field(i).String()
		}
	}
	return result, nil
}

func generateNonceStr() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

func signature(params map[string]string, key string) string {
	excluded := map[string]string{}
	for k, v := range params {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		excluded[k] = v
	}

	var keyValueConcat []string
	for k, v := range excluded {
		keyValueConcat = append(keyValueConcat, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(keyValueConcat)
	keyValueStr := strings.Join(keyValueConcat, "&")

	keyValueSecret := keyValueStr + "&key=" + key
	return fmt.Sprintf("%X", md5.Sum([]byte(keyValueSecret)))
}

func toXMLStr(params map[string]string) string {
	xml := "<xml>"
	for k, v := range params {
		xml += fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xml += "</xml>"
	return xml
}
