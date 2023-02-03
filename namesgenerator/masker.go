package namesgenerator

import (
	"strings"

	"github.com/docker/docker/pkg/namesgenerator"
)

type Masker struct {
	mask map[string]string
}

func NewMasker() Masker {
	return Masker{
		mask: make(map[string]string),
	}
}

func (m Masker) Register(s string) {
	if _, ok := m.mask[s]; !ok {
		m.mask[s] = namesgenerator.GetRandomName(0)
	}
}

func (m Masker) MaskString(s string) string {
	return m.replaceStrings(s, m.mask)
}

func (m Masker) UnmaskString(s string) string {
	return m.replaceStrings(s, m.reverseMaskMap())
}

func (m Masker) replaceStrings(s string, mask map[string]string) string {
	result := s
	for old, new := range mask {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

func (m Masker) reverseMaskMap() map[string]string {
	r := make(map[string]string)
	for k, v := range m.mask {
		r[v] = k
	}
	return r
}
