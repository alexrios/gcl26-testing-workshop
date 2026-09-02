package fuzzing

import "testing"

func TestEntryRoundTrip(t *testing.T) {
	want := Entry{Key: "conference", Value: "GopherCon LATAM"}
	data, err := Encode(want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Decode(Encode(%#v)) = %#v", want, got)
	}
}
