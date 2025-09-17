package presto

import (
    "encoding/json"
    "fmt"
    "reflect"
)

// normalizeToInterfaceSlice accepts the raw `value` that Scan receives and returns a []interface{}
// representation, accepting these input shapes:
//  - []interface{}                -> returned as-is
//  - JSON string (e.g. '["a","b"]') or []byte -> json.Unmarshal into []interface{}
//  - []T where T is concrete (e.g. []string, []int64) -> convert by reflection
// Returns an error if it cannot be converted.
func normalizeToInterfaceSlice(value interface{}) ([]interface{}, error) {
    if value == nil {
        return nil, nil
    }

    switch v := value.(type) {
    case []interface{}:
        return v, nil
    case string:
        var tmp []interface{}
        if err := json.Unmarshal([]byte(v), &tmp); err == nil {
            return tmp, nil
        } else {
            return nil, fmt.Errorf("presto: cannot unmarshal string into slice: %v", err)
        }
    case []byte:
        var tmp []interface{}
        if err := json.Unmarshal(v, &tmp); err == nil {
            return tmp, nil
        } else {
            return nil, fmt.Errorf("presto: cannot unmarshal bytes into slice: %v", err)
        }
    default:
        // If it's a slice of a concrete type (e.g. []string, []int64), use reflection to
        // copy elements into []interface{}.
        rv := reflect.ValueOf(value)
        if rv.Kind() == reflect.Slice {
            l := rv.Len()
            out := make([]interface{}, l)
            for i := 0; i < l; i++ {
                out[i] = rv.Index(i).Interface()
            }
            return out, nil
        }
    }

    return nil, fmt.Errorf("presto: cannot convert %v (%T) to []interface{}", value, value)
}
