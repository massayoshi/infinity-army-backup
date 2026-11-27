package main

import "testing"

func TestParseRedirect(t *testing.T) {
	cases := []struct {
		content string
		target  string
		ok      bool
	}{
		{"#REDIRECT [[Combat_Module#Burst_.28B.29]]", "Combat_Module#Burst_.28B.29", true},
		{"#redirect [[Skills_and_Equipment_Module#Common]]\n", "Skills_and_Equipment_Module#Common", true},
		{"  #REDIRECT [[ Line of Fire | LoF ]]", "Line of Fire", true},
		{"{{n5-section-basic-rules}}\nCoherency is the game mechanic...", "", false},
		{"See #redirect [[Elsewhere]] in the text", "", false},
	}
	for _, c := range cases {
		target, ok := parseRedirect(c.content)
		if target != c.target || ok != c.ok {
			t.Errorf("parseRedirect(%q) = %q, %v; want %q, %v", c.content, target, ok, c.target, c.ok)
		}
	}
}

func TestPageFolder(t *testing.T) {
	if got := pageFolder(nil); got != "Uncategorized" {
		t.Errorf("pageFolder(nil) = %q, want Uncategorized", got)
	}
	if got := pageFolder([]string{"Category:Ammunition & Weaponry/N5"}); got != "Ammunition_and_WeaponryN5" {
		t.Errorf("pageFolder = %q, want Ammunition_and_WeaponryN5", got)
	}
	if got := pageFolder([]string{"Category:N5.2 Update", "Category:Optional", "Category:Special Skills"}); got != "Optional" {
		t.Errorf("pageFolder skips release tags: got %q, want Optional", got)
	}
	if got := pageFolder([]string{"Category:N5 FAQ 0.1", "Category:N5.3 Update"}); got != "N5_FAQ_0.1" {
		t.Errorf("pageFolder only release tags: got %q, want N5_FAQ_0.1", got)
	}
}

func TestFactionFolderName(t *testing.T) {
	slugCount := map[string]int{"contracted-back-up": 2, "panoceania": 1}
	if got := factionFolderName("panoceania", 101, slugCount); got != "panoceania" {
		t.Errorf("unique slug: got %q", got)
	}
	if got := factionFolderName("contracted-back-up", 999, slugCount); got != "contracted-back-up-999" {
		t.Errorf("shared slug: got %q", got)
	}
}

func TestNamespaceListURL(t *testing.T) {
	got := namespaceListURL("https://example.com/api.php?action=query&list=allpages", namespaceTemplate)
	want := "https://example.com/api.php?action=query&list=allpages&apnamespace=10"
	if got != want {
		t.Errorf("namespaceListURL = %q, want %q", got, want)
	}
}

func TestStaleKey(t *testing.T) {
	if staleKey("wiki/Optional/Doctor.wiki") != staleKey("wiki/N5.2_Update/Doctor.wiki") {
		t.Error("same page in different folders should share a key")
	}
	if staleKey("wiki/Templates/Activation.wiki") == staleKey("wiki/Skills/Activation.wiki") {
		t.Error("a template must not share a key with a regular page")
	}
}
