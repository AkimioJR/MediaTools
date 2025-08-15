package storage

import (
	"encoding/json"
	"strings"
)

type StorageType uint8

const (
	StorageUnknown StorageType = iota // 未知文件系统
	StorageLocal                      // 本地文件系统
)

func (t StorageType) String() string {
	switch t {
	case StorageLocal:
		return "LocalStorage"
	default:
		return "UnknownStorage"
	}
}

func ParseStorageType(s string) StorageType {
	switch strings.ToLower(s) {
	case "localstorage":
		return StorageLocal
	default:
		return StorageUnknown
	}
}

func (t StorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *StorageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = ParseStorageType(s)
	return nil
}

func (t StorageType) MarshalYAML() (any, error) {
	return t.String(), nil
}

func (t *StorageType) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	*t = ParseStorageType(s)
	return nil
}
