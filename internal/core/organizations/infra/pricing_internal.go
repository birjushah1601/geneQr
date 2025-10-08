package infra

// internal helpers to validate precedence logic via unit tests

type priceCandidate struct {
    hasOrg     bool
    hasChannel bool
    price      float64
    currency   string
}

func precedence(hasOrg, hasChannel bool) int {
    switch {
    case hasOrg && hasChannel:
        return 4
    case !hasOrg && hasChannel:
        return 3
    case hasOrg && !hasChannel:
        return 2
    default:
        return 1
    }
}

func pickBest(cands []priceCandidate) *priceCandidate {
    if len(cands) == 0 { return nil }
    best := cands[0]
    bestScore := precedence(best.hasOrg, best.hasChannel)
    for i := 1; i < len(cands); i++ {
        s := precedence(cands[i].hasOrg, cands[i].hasChannel)
        if s > bestScore {
            best = cands[i]
            bestScore = s
        }
    }
    return &best
}
