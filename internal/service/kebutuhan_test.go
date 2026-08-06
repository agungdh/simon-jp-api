package service

import "testing"

func TestKategoriJumlahMinimal(t *testing.T) {
	cases := []struct {
		kategori string
		want     int
	}{
		{"admin", 20},
		{"pejabat", 40},
		{"auditor", 40},
		{"", 0},
		{"unknown", 0},
	}
	for _, c := range cases {
		if got := KategoriJumlahMinimal(c.kategori); got != c.want {
			t.Fatalf("KategoriJumlahMinimal(%q) = %d, want %d", c.kategori, got, c.want)
		}
	}
}
