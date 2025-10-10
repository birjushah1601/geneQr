package app

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestDeliver_SendsHeadersAndSignature(t *testing.T) {
    var gotEvent, gotSig, gotTS string
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gotEvent = r.Header.Get("X-Webhook-Event")
        gotSig = r.Header.Get("X-Webhook-Signature")
        gotTS = r.Header.Get("X-Webhook-Timestamp")
        _, _ = io.ReadAll(r.Body)
        w.WriteHeader(200)
    }))
    defer srv.Close()

    d := &WebhookDispatcher{client: srv.Client()}
    secret := "test_secret"
    payload := []byte(`{"ok":true}`)
    j := deliveryJob{EventType: "ticket.created", Payload: payload, EndpointURL: srv.URL, Secret: &secret}
    if err := d.deliver(context.Background(), j); err != nil {
        t.Fatalf("deliver error: %v", err)
    }
    if gotEvent != "ticket.created" { t.Fatalf("wrong event header: %s", gotEvent) }
    if gotTS == "" { t.Fatalf("missing timestamp header") }
    if !strings.HasPrefix(gotSig, "t=") || !strings.Contains(gotSig, ",v1=") { t.Fatalf("bad signature format: %s", gotSig) }
    // verify HMAC
    parts := strings.Split(gotSig, ",")
    var ts, v1 string
    for _, p := range parts {
        if strings.HasPrefix(p, "t=") { ts = strings.TrimPrefix(p, "t=") }
        if strings.HasPrefix(p, "v1=") { v1 = strings.TrimPrefix(p, "v1=") }
    }
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(ts))
    mac.Write([]byte("."))
    mac.Write(payload)
    want := hex.EncodeToString(mac.Sum(nil))
    if v1 != want { t.Fatalf("signature mismatch: got %s want %s", v1, want) }
}

func TestDeliver_Non2xxReturnsError(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(500)
    }))
    defer srv.Close()
    d := &WebhookDispatcher{client: srv.Client()}
    j := deliveryJob{EventType: "x", Payload: []byte("{}"), EndpointURL: srv.URL}
    if err := d.deliver(context.Background(), j); err == nil {
        t.Fatalf("expected error on non-2xx response")
    }
}
