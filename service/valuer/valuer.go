package valuer

import (
  "errors"
)

var (
  // ErrKeyNotFound - when the key cannot be found
  ErrKeyNotFound = errors.New("Could not find associated key in map")

  // ErrFormatInvalid - when the format is invalid
  ErrFormatInvalid = errors.New("Format is invalid")
)

// Service helps change interface{} from JSON to a specific value.

// ToFloat will convert a value to a float
func ToFloat64(m map[string]interface{}, key string) (float64, error) {
  val, ok := m[key]
  if !ok {
    return 0, ErrKeyNotFound
  }

  valAsFloat, ok := val.(float64)
  if !ok {
    return 0, ErrFormatInvalid
  }

  return valAsFloat, nil
}

// ToInt will convert a value to an integer
func ToInt(m map[string]interface{}, key string) (int, error) {
  val, err := ToFloat64(m, key)
  if err != nil {
    return 0, err
  }
  return int(val), nil
}

// ToInt64 will convert a value to an integer
func ToInt64(m map[string]interface{}, key string) (int64, error) {
  val, err := ToFloat64(m, key)
  if err != nil {
    return 0, err
  }
  return int64(val), nil
}

// ToUint will convert a value to an integer
func ToUint(m map[string]interface{}, key string) (uint, error) {
  val, err := ToFloat64(m, key)
  if err != nil {
    return 0, err
  }
  return uint(val), nil
}

// ToString will convert a value to a float
func ToString(m map[string]interface{}, key string) (string, error) {
  val, ok := m[key]
  if !ok {
    return "", ErrKeyNotFound
  }

  valAsString, ok := val.(string)
  if !ok {
    return "", ErrFormatInvalid
  }

  return valAsString, nil
}

// ToBoolean will convert a value to a bool
func ToBoolean(m map[string]interface{}, key string) (bool, error) {
  val, ok := m[key]
  if !ok {
    return false, ErrKeyNotFound
  }

  valAsBool, ok := val.(bool)
  if !ok {
    return false, ErrFormatInvalid
  }

  return valAsBool, nil
}

// ToMap will convert a value to a bool
func ToMap(m map[string]interface{}, key string) (map[string]interface{}, error) {
  val, ok := m[key]
  if !ok {
    return nil, ErrKeyNotFound
  }

  valAsMap, ok := val.(map[string]interface{})
  if !ok {
    return nil, ErrFormatInvalid
  }

  return valAsMap, nil
}

// ToArray will convert a value to a bool
func ToArray(m map[string]interface{}, key string) ([]interface{}, error) {
  val, ok := m[key]
  if !ok {
    return nil, ErrKeyNotFound
  }

  valAsArray, ok := val.([]interface{})
  if !ok {
    return nil, ErrFormatInvalid
  }

  return valAsArray, nil
}
