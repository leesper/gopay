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

	var sortedkeys []string
	for k := range excluded {
		sortedkeys = append(sortedkeys, k)
	}
	sort.Strings(sortedkeys)

	var keyValueConcat []string
	for _, k := range sortedkeys {
		keyValueConcat = append(keyValueConcat, fmt.Sprintf("%s=%s", k, excluded[k]))
	}
	keyValueStr := strings.Join(keyValueConcat, "&")

	keyValueSecret := keyValueStr + "&key=" + key
	fmt.Printf("params %v, keyValueSecret %s\n", params, keyValueSecret)
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

func generateTimestampStr() string {
	const ChinaTimeZoneOffset = 8 * 60 * 60 // UTC + 8
	return fmt.Sprintf("%d", time.Now().Unix()+ChinaTimeZoneOffset)
}
