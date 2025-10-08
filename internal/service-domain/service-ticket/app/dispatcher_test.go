package app

import "testing"

func TestSign(t *testing.T) {
    secret := "test_secret"
    ts := "1700000000"
    body := []byte(`{"hello":"world"}`)
    got := sign(secret, ts, body)
    // Expected for HMAC_SHA256(secret, ts+"."+body)
    want := "a5d9de1d23a233e086ce2cd0940311c91c3039e7d92c0d96a82fda4881746747"
    if got != want {
        t.Fatalf("unexpected signature: got %s, want %s", got, want)
    }
}
