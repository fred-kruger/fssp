package fssp

import "strings"

type Physical struct {
	Region     int
	Firstname  string
	Secondname string
	Lastname   string
	Birthdate  Birthdate
}

type Birthdate struct {
	Day   string
	Month string
	Year  string
}

func (p *Birthdate) SetBirthDay(day string) {
	p.Day = day;
}

func (p *Birthdate) SetBirthMonth(month string) {
	p.Month = month
}

func (p *Birthdate) SetBirthYear(year string) {
	p.Year = year
}

func (p Birthdate) GetBirthdate() string {
	var birthdate strings.Builder
	birthdate.Grow(10)
	birthdate.WriteString(p.Day)
	birthdate.WriteString(".")
	birthdate.WriteString(p.Month)
	birthdate.WriteString(".")
	birthdate.WriteString(p.Year)

	return birthdate.String()
}
