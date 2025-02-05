package models

import (
	"time"
)

type Person string
type Medicine string

type PosologyEntry struct {
	HeavierThan      int64         `sheet:"Minimum Weight"`
	OlderThan        time.Duration `sheet:"Minimum Age"`
	Dose             string        `sheet:"Dose"`
	DoseInterval     time.Duration `sheet:"Dose interval"`
	MaxDoses         int64         `sheet:"Max doses"`
	MaxDosesInterval time.Duration `sheet:"Interval"`
}

type PersonCfg struct {
	Name     Person    `sheet:"Name"`
	Birth    time.Time `sheet:"Birthdate"`
	Weight   int64     `sheet:"Weight"`
	PhotoUrl string    `sheet:"Photo"`
}

type MedicineCfg struct {
	Posology []PosologyEntry
}

type Dose struct {
	Who  Person    `sheet:"Person"`
	What Medicine  `sheet:"Medicine"`
	When time.Time `sheet:"When,2006-01-02 15:04:05"`
}
