package timezone

import "testing"

func TestIsValid(t *testing.T) {
	cases := map[string]bool{
		"Asia/Jakarta":  true,
		"Asia/Makassar": true,
		"Asia/Jayapura": true,
		"Asia/Singapore": false,
		"":               false,
		"WIB":            false,
	}
	for tz, want := range cases {
		if got := IsValid(tz); got != want {
			t.Errorf("IsValid(%q) = %v, want %v", tz, got, want)
		}
	}
}

func TestResolve_ZoneOverrideWins(t *testing.T) {
	zoneTz := "Asia/Jayapura"
	loc, err := Resolve("Asia/Jakarta", &zoneTz)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jayapura" {
		t.Errorf("got %s, want Asia/Jayapura", loc.String())
	}
}

func TestResolve_FallsBackToCompanyWhenZoneNil(t *testing.T) {
	loc, err := Resolve("Asia/Makassar", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Makassar" {
		t.Errorf("got %s, want Asia/Makassar", loc.String())
	}
}

func TestResolve_FallsBackToJakartaWhenCompanyEmpty(t *testing.T) {
	loc, err := Resolve("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("got %s, want Asia/Jakarta", loc.String())
	}
}

func TestResolve_InvalidCompanyTimezone(t *testing.T) {
	_, err := Resolve("Not/AZone", nil)
	if err != ErrInvalidTimezone {
		t.Errorf("got %v, want ErrInvalidTimezone", err)
	}
}
