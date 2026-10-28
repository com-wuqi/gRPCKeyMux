package config

import (
	"encoding/json"
	"testing"
)

func TestMakcBackendNameUnmarshal(t *testing.T) {
	cases := []struct {
		in      string
		want    MakcBackendName
		wantErr bool
	}{
		{`"auto"`, MakcBackendAuto, false},
		{`"uinput"`, MakcBackendUInput, false},
		{`"xdgportal"`, MakcBackendXDGPortal, false},
		{`""`, "", false},
		{`"sendinput"`, "", true},
	}
	for _, c := range cases {
		var got MakcBackendName
		err := json.Unmarshal([]byte(c.in), &got)
		if (err != nil) != c.wantErr {
			t.Errorf("Unmarshal(%s) err = %v, wantErr %v", c.in, err, c.wantErr)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("Unmarshal(%s) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMakcBackendDefaults(t *testing.T) {
	var c Config
	c.applyDefaults()
	if c.MakcKeyboardBackend != MakcBackendAuto || c.MakcMouseBackend != MakcBackendAuto {
		t.Fatalf("defaults = (%q, %q), want auto/auto", c.MakcKeyboardBackend, c.MakcMouseBackend)
	}

	c.MakcKeyboardBackend = MakcBackendXDGPortal
	c.MakcMouseBackend = MakcBackendUInput
	c.applyDefaults()
	if c.MakcKeyboardBackend != MakcBackendXDGPortal || c.MakcMouseBackend != MakcBackendUInput {
		t.Fatalf("explicit values were overwritten: (%q, %q)", c.MakcKeyboardBackend, c.MakcMouseBackend)
	}
}
