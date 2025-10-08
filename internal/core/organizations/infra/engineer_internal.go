package infra

// Helpers to validate engineer eligibility logic without DB

// covEntry represents a coverage row for an engineer
type covEntry struct {
    region string
    skills []string
}

func hasAllSkills(required, have []string) bool {
    if len(required) == 0 { return true }
    if len(have) == 0 { return false }
    set := make(map[string]struct{}, len(have))
    for _, s := range have { set[s] = struct{}{} }
    for _, r := range required {
        if _, ok := set[r]; !ok { return false }
    }
    return true
}

// eligible returns true if engineer is eligible per skills/region based on either
// direct engineer skills/region OR any coverage entry matching.
func eligible(requiredSkills []string, requiredRegion string, engineerSkills []string, engineerRegion string, cover []covEntry) bool {
    // Skills check
    skillsOK := len(requiredSkills) == 0 || hasAllSkills(requiredSkills, engineerSkills)
    // Region check
    regionOK := requiredRegion == "" || engineerRegion == requiredRegion
    if skillsOK && regionOK { return true }
    // Check coverage entries
    for _, c := range cover {
        if (requiredRegion == "" || c.region == requiredRegion) && hasAllSkills(requiredSkills, c.skills) {
            return true
        }
    }
    return false
}
