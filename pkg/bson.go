package pkg

import "encoding/json"

func Serialize[K any](obj K) ([]byte, error) {
	return json.Marshal(obj)
}

func Deserialize[K any](b []byte) (obj K, err error) {

	err = json.Unmarshal(b, &obj)
	if err != nil {
		return obj, err
	}
	return
}
