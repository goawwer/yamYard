package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goawwer/yamyard/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func CreateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), err
}

var (
	dbFieldCache = map[string]map[string]string{}
)

func GetDBFieldMap(model any) map[string]string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if cached, ok := dbFieldCache[t.Name()]; ok {
		return cached
	}

	m := buildDBFieldMap(model)
	dbFieldCache[t.Name()] = m
	return m
}

func buildDBFieldMap(model any) map[string]string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	m := make(map[string]string)

	for i := 1; i < t.NumField(); i++ {
		dbTags := t.Field(i).Tag.Get("db")

		if dbTags == "" || dbTags == "-" {
			continue
		}

		col := strings.Split(dbTags, ",")[0]

		m[strings.ToLower(t.Field(i).Name)] = col
	}

	return m
}

func GetGetDBFieldMaps() {
	GetDBFieldMap(domain.User{})
	GetDBFieldMap(domain.Recipe{})
}
