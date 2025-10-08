package infra

import "testing"

func TestPrecedenceRanking(t *testing.T) {
    cases := []struct{
        org, ch bool
        want int
    }{
        {false, false, 1},
        {true, false, 2},
        {false, true, 3},
        {true, true, 4},
    }
    for _, c := range cases {
        if got := precedence(c.org, c.ch); got != c.want {
            t.Fatalf("precedence(%v,%v)=%d want %d", c.org, c.ch, got, c.want)
        }
    }
}

func TestPickBestCandidate(t *testing.T) {
    cands := []priceCandidate{
        {hasOrg:false, hasChannel:false, price:1000, currency:"USD"}, // global
        {hasOrg:true, hasChannel:false, price:950, currency:"USD"},   // org
        {hasOrg:false, hasChannel:true, price:980, currency:"USD"},   // channel
        {hasOrg:true, hasChannel:true, price:940, currency:"USD"},    // org_channel
    }
    best := pickBest(cands)
    if best == nil || !(best.hasOrg && best.hasChannel) || best.price != 940 {
        t.Fatalf("expected org_channel candidate with price 940; got %#v", best)
    }
}

func TestPickBestWhenMissingSomeScopes(t *testing.T) {
    cands := []priceCandidate{
        {hasOrg:false, hasChannel:false, price:1000, currency:"USD"}, // global
        {hasOrg:true, hasChannel:false, price:970, currency:"USD"},   // org
    }
    best := pickBest(cands)
    if best == nil || !best.hasOrg || best.hasChannel || best.price != 970 {
        t.Fatalf("expected org candidate; got %#v", best)
    }
}
