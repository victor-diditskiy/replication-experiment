package entity

import (
	"math/rand"
	"time"
)

var names = []string{
	"Moscow",
	"Saint-Petersburg",
	"London",
	"New York",
	"Vladimir",
	"Krasnodar",
	"Orel",
	"Velikiy Novgorod",
	"Kazan",
	"Nalchik",
	"Sochi",
	"Yalta",
	"Simferopol",
	"Voronezh",
	"Tumen",
}

type Data struct {
	ID        int64
	Name      string
	Value     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func RandomData() Data {
	l := len(names) - 1
	i := rand.Intn(l)
	return Data{
		Name:  names[i],
		Value: rand.Int63n(10 ^ 9),
	}
}
