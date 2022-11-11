package covid

import (
	"log"
	"testing"
)

func TestGetCovidResponse(t *testing.T) {
	cs, err := GetCovidResponse("美国")
	if err != nil {
		t.Errorf("err : %v", err)
	}
	log.Println(cs)
	log.Println(PrintCovidSituation(cs))
}

func BenchmarkLen(b *testing.B) {
	a := struct {
		len   int
		temps []int
	}{len: 0, temps: []int{}}
	for i := 0; i < b.N; i++ {
		if a.len == 0 {
			continue
		}
	}
}

func BenchmarkNum(b *testing.B) {
	a := struct {
		len   int
		temps []int
	}{len: 0, temps: []int{}}
	for i := 0; i < b.N; i++ {
		if len(a.temps) == 0 {
			continue
		}
	}
}

func BenchmarkLenStrign(b *testing.B) {
	a := struct {
		len   string
		temps []int
	}{len: "0", temps: []int{}}
	for i := 0; i < b.N; i++ {
		if a.len == "0" {
			continue
		}
	}
}
