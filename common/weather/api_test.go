package weather

import (
	"encoding/json"
	"testing"
)

func TestGetWeather(t *testing.T) {
	w, err := GetWeather("101280604")
	if err != nil || w.errCode != "" {
		t.Fatalf("err:%v\nerrCode:%v\nerrDesc:%v\n", err, w.errCode, w.errDesc)
	}
	j, _ := json.Marshal(w)
	t.Logf("%v", string(j))
}
