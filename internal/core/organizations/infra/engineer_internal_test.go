package infra

import "testing"

func TestHasAllSkills(t *testing.T) {
    if !hasAllSkills(nil, []string{"x"}) { t.Fatalf("nil required should be true") }
    if hasAllSkills([]string{"a"}, nil) { t.Fatalf("missing skills should be false") }
    if !hasAllSkills([]string{"a","b"}, []string{"a","b","c"}) { t.Fatalf("superset should be true") }
    if hasAllSkills([]string{"a","d"}, []string{"a","b","c"}) { t.Fatalf("missing d should be false") }
}

func TestEligible_DirectEngineer(t *testing.T) {
    if !eligible([]string{"s1"}, "r1", []string{"s1","s2"}, "r1", nil) {
        t.Fatalf("expected eligible with direct skills and region match")
    }
    if eligible([]string{"s1","s3"}, "r1", []string{"s1","s2"}, "r2", nil) {
        t.Fatalf("should be ineligible: region mismatch and missing skill")
    }
}

func TestEligible_ViaCoverage(t *testing.T) {
    cover := []covEntry{
        {region:"r2", skills:[]string{"s1","s3"}},
        {region:"r3", skills:[]string{"s2"}},
    }
    // Engineer lacks direct region/skills, but coverage matches
    if !eligible([]string{"s1","s3"}, "r2", []string{"s1"}, "r1", cover) {
        t.Fatalf("expected eligible via coverage entry r2 with s1,s3")
    }
    // No coverage region match
    if eligible([]string{"s2"}, "r1", []string{"s1"}, "rX", cover) {
        t.Fatalf("should be ineligible: no region match and coverage not r1")
    }
}
