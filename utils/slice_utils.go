// utils/slice_utils.go
package utils

func ConvertToStringSlice(data interface{}) []string {
	slice, ok := data.([]interface{})
	if !ok {
		return []string{}
	}
	
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.(string)
	}
	return result
}