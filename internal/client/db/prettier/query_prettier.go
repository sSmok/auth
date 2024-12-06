package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

// DollarPlaceholder - символ, который используется текущей БД при подставлении аргументов запроса
const DollarPlaceholder = "$"

// Pretty - функция, форматируящая запрос в БД под формат логов
func Pretty(query string, placeholder string, args ...interface{}) string {
	for i, arg := range args {
		var value string
		switch v := arg.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), value, -1)
	}
	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
