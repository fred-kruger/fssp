package fssp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBirthdate_GetBirthdate(t *testing.T) {
	assert.Equal(t,bd.GetBirthdate(),"23.03.1999")
}

func TestBirthdate_SetBirthDay(t *testing.T) {
	bd.SetBirthDay("01")

	assert.Equal(t,bd.GetBirthdate(),"01.03.1999")
}

func TestBirthdate_SetBirthMonth(t *testing.T) {
	bd.SetBirthMonth("01")

	assert.Equal(t,bd.GetBirthdate(),"01.01.1999")
}

func TestBirthdate_SetBirthYear(t *testing.T) {
	bd.SetBirthYear("2000")

	assert.Equal(t,bd.GetBirthdate(),"01.01.2000")
}

var bd = Birthdate{"23","03","1999"}
