package masker

import (
	"errors"
	"strings"

	"github.com/docker/docker/pkg/namesgenerator"
)

var ErrAlreadyMasked = errors.New("this variable is already masked")

type Masker struct {
	mask   map[string]string
	maskFn func(string) string
}

func New() *Masker {
	return &Masker{
		mask: make(map[string]string),
		maskFn: func(s string) string {
			return namesgenerator.GetRandomName(0)
		},
	}
}

func NewWithFn(fn func(string) string) *Masker {
	return &Masker{
		mask:   make(map[string]string),
		maskFn: fn,
	}
}

func (m Masker) Add(s string) error {
	if _, ok := m.mask[s]; !ok {
		m.mask[s] = m.maskFn(s)
		return nil
	}
	return ErrAlreadyMasked
}

func (m Masker) MustAdd(s string) {
	if err := m.Add(s); err != nil {
		panic(err)
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
