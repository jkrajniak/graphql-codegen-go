package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func inArray(item string, items []string) bool {
	itemLower := strings.ToLower(item)
	for _, v := range items {
		itemVal := strings.TrimSpace(strings.ToLower(v))
		if itemVal == itemLower {
			return true
		}
	}
	return false
}

type keywordMap map[string]bool

func (t keywordMap) Set(key string) error {
	if _, has := t[key]; has {
		return errors.New(fmt.Sprintf("keyword %s already declared", key))
	}
	t[key] = true
	return nil
}

func (t keywordMap) Has(key string) bool {
	if _, has := t[key]; has {
		return true
	}
	return false
}

func (t keywordMap) GetMissingKeys(keys []string) []string {
	var missingKeys []string
	for _, k := range keys {
		if !t.Has(k) {
			missingKeys = append(missingKeys, k)
		}
	}
	return missingKeys
}
