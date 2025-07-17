package main

import "testing"

func TestUsedCar_Price(t *testing.T) {
	tests := []struct {
		name   string
		rules  []PriceRule
		want   int
	}{
		{
			name:  "SUV・高級・3万km・5年",
			rules: []PriceRule{
				TypeRule{CarType: "SUV"},
				ModelRule{Model: "高級"},
				MileageRule{Mileage: 30000},
				YearRule{Years: 5},
			},
			want: 200 + 100 + 80 - 30 - 75, // 275
		},
		{
			name:  "セダン・中級・5万km・10年",
			rules: []PriceRule{
				TypeRule{CarType: "セダン"},
				ModelRule{Model: "中級"},
				MileageRule{Mileage: 50000},
				YearRule{Years: 10},
			},
			want: 200 + 50 + 50 - 50 - 150, // 100
		},
		{
			name:  "軽・低級・0km・0年",
			rules: []PriceRule{
				TypeRule{CarType: "軽"},
				ModelRule{Model: "低級"},
				MileageRule{Mileage: 0},
				YearRule{Years: 0},
			},
			want: 200 + 20 + 20 - 0 - 0, // 240
		},
		{
			name:  "価格が0未満になるケース",
			rules: []PriceRule{
				TypeRule{CarType: "セダン"},
				ModelRule{Model: "低級"},
				MileageRule{Mileage: 200000},
				YearRule{Years: 30},
			},
			want: 0, // マイナスは0に補正
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			car := NewUsedCar(tt.rules...)
			got := car.Price()
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTypeRule_Apply(t *testing.T) {
	cases := []struct {
		typeName string
		want    int
	}{
		{"SUV", 300},
		{"セダン", 250},
		{"軽", 220},
		{"不明", 200},
	}
	for _, c := range cases {
		got := TypeRule{CarType: c.typeName}.Apply(200)
		if got != c.want {
			t.Errorf("TypeRule(%s): got %d, want %d", c.typeName, got, c.want)
		}
	}
}

func TestModelRule_Apply(t *testing.T) {
	cases := []struct {
		model string
		want  int
	}{
		{"高級", 280},
		{"中級", 250},
		{"低級", 220},
		{"不明", 200},
	}
	for _, c := range cases {
		got := ModelRule{Model: c.model}.Apply(200)
		if got != c.want {
			t.Errorf("ModelRule(%s): got %d, want %d", c.model, got, c.want)
		}
	}
}

func TestMileageRule_Apply(t *testing.T) {
	cases := []struct {
		mileage int
		want    int
	}{
		{0, 200},
		{10000, 190},
		{50000, 150},
		{200000, 0},
	}
	for _, c := range cases {
		got := MileageRule{Mileage: c.mileage}.Apply(200)
		if got != c.want {
			t.Errorf("MileageRule(%d): got %d, want %d", c.mileage, got, c.want)
		}
	}
}

func TestYearRule_Apply(t *testing.T) {
	cases := []struct {
		years int
		want  int
	}{
		{0, 200},
		{1, 185},
		{10, 50},
		{20, -100},
	}
	for _, c := range cases {
		got := YearRule{Years: c.years}.Apply(200)
		if got != c.want {
			t.Errorf("YearRule(%d): got %d, want %d", c.years, got, c.want)
		}
	}
} 