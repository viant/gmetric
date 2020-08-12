package counter

import "github.com/viant/gmetric/stat"

//NormalizeValue normalize stat values
func NormalizeValue(values []interface{}) []interface{} {
	if len(values) == 1 {
		switch statValue := values[0].(type) {
		case *stat.Values:
			values = *statValue
		case stat.Values:
			values = statValue
		}
	}
	return values
}
