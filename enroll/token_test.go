package enroll

import "testing"

func TestRoundTrip(t *testing.T) {
	fp := Fingerprint([]byte("x"))
	tok := Token{CAFingerprint: fp, Secret: "abc"}
	got, err := ParseToken(tok.String())
	if err != nil || got != tok {
		t.Fatalf("got %v %v", got, err)
	}
	for _, bad := range []string{"", "tb1.x.y", "tb2." + fp + ".s", "tb1." + fp + "."} {
		if _, err := ParseToken(bad); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}
