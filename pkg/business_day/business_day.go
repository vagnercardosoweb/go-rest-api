package business_day

import (
	"time"
)

// Returns the date for Easter Sunday for a given year
func easterSunday(year int) time.Time {
	if year < 1583 {
		year = 1583
	}

	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	r := 22 + h + l - 7*m

	return time.Date(
		year,
		time.March,
		r,
		0,
		0,
		0,
		0,
		time.Local,
	)
}

// Map of national holidays
var holidays = map[string]string{
	"01/01": "Ano Novo",
	"21/04": "Tiradentes",
	"01/05": "Dia do Trabalho",
	"07/09": "Independência do Brasil",
	"12/10": "Nossa Senhora Aparecida",
	"02/11": "Finados",
	"15/11": "Proclamação da República",
	"25/12": "Natal",
}

func Is(date time.Time) bool {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return false
	}

	if _, ok := holidays[date.Format("02/01")]; ok {
		return false
	}

	return true
}

func Next(date time.Time, inclusive bool) time.Time {
	easterTime := easterSunday(date.Year())

	if _, ok := holidays[easterTime.Format("02/01")]; !ok {
		holidays[easterTime.Format("02/01")] = "Páscoa"

		goodFriday := easterTime.AddDate(0, 0, -2)
		easterMonday := easterTime.AddDate(0, 0, 1)
		carnavalMonday := easterTime.AddDate(0, 0, -48)
		carnavalTuesday := easterTime.AddDate(0, 0, -47)
		corpusChristi := easterTime.AddDate(0, 0, 60)

		holidays[goodFriday.Format("02/01")] = "Sexta-feira Santa"
		holidays[easterMonday.Format("02/01")] = "Segunda-feira de Páscoa"
		holidays[carnavalMonday.Format("02/01")] = "Segunda-feira de Carnaval"
		holidays[carnavalTuesday.Format("02/01")] = "Terça-feira de Carnaval"
		holidays[corpusChristi.Format("02/01")] = "Corpus Christi"
	}

	days := 1
	if inclusive {
		days = 0
	}

	for {
		date = date.AddDate(0, 0, days)
		if Is(date) {
			return date
		}
	}
}
