package codec

import "github.com/gorilla/schema"

var (
	DefaultQueryUnmarshaler = schema.NewDecoder()

	queryUnmarshalers = map[string]*schema.Decoder{}
)

func init() {
	DefaultQueryUnmarshaler.SetAliasTag("json")
	DefaultQueryUnmarshaler.IgnoreUnknownKeys(true)
}

func D() {

}

func UnmarshalQuery(endpoint string, src map[string][]string, dst interface{}) error {
	decoder, ok := queryUnmarshalers[endpoint]
	if !ok {
		return DefaultQueryUnmarshaler.Decode(dst, src)
	}
	return decoder.Decode(dst, src)
}

func RegQueryUnmarshaler(endpoint string, d *schema.Decoder) {
	queryUnmarshalers[endpoint] = d
}
