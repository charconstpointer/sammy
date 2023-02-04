package masker_test

import (
	"testing"

	"github.com/charconstpointer/sammy/masker"
)

var (
	maskerFn = func(s string) string {
		switch s {
		case "test":
			return "tset"
		case "foo":
			return "oof"
		case "oh_noez_my_creds":
			return "nothing to see here"
		default:
			return "tluafed"
		}
	}
)

func TestMasker_RoundTrip(t *testing.T) {
	t.Run(t.Name(), func(t *testing.T) {
		m := masker.NewWithFn(maskerFn)
		s := "foo bar oh_noez_my_creds test"
		if err := m.Add("oh_noez_my_creds"); err != nil {
			t.Fatalf("Masker.Add() = %v, want nil", err)
		}
		masked := m.MaskString(s)
		if masked != "foo bar nothing to see here test" {
			t.Errorf("Masker.MaskString() = %s, want 'foo bar nothing to see here test'", masked)
		}
		unmasked := m.UnmaskString(masked)
		if unmasked != s {
			t.Errorf("Masker.UnmaskString() = %s, want %s", unmasked, s)
		}
	})
}

func TestMasker_Add(t *testing.T) {
	t.Run(t.Name(), func(t *testing.T) {
		m := masker.NewWithFn(maskerFn)
		if err := m.Add("test"); err != nil {
			t.Errorf("Masker.Add() = %v, want nil", err)
		}
		if err := m.Add("test"); err == nil {
			t.Errorf("Masker.Add() = nil, want error")
		}
	})
}
