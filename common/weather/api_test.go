package weather

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGetWeather(t *testing.T) {
	w, err := GetWeather("101280604")
	if err != nil || w.errCode != "" {
		t.Fatalf("err:%v\nerrCode:%v\nerrDesc:%v\n", err, w.errCode, w.errDesc)
	}
	j, _ := json.Marshal(w)
	t.Logf("%v", string(j))
	t.Logf("%v", w.GetAQIInfo())
	t.Logf("%v", w.GetCurrentWeatherInfo())
}

func TestGetCityLike(t *testing.T) {
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    *CityLikeResp
		wantErr bool
	}{
		{
			name: "ok",
			args: args{city: "南山"},
			want: &CityLikeResp{
				Status:  200,
				Message: "OK",
				Data: map[string]string{
					"101051206": "黑龙江, 鹤岗, 南山",
					"101280604": "广东, 深圳, 南山",
				},
			},
			wantErr: false,
		},
		{
			name:    "lack city",
			args:    args{city: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCityLike(tt.args.city)
			t.Logf("got: %v, err: %v\n", got, err)
			if got != nil {
				t.Logf("res: %v\n", got.GetCityLike())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCityLike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCityLike() = %v, want %v", got, tt.want)
			}
		})
	}
}
