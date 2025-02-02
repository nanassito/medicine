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
	Birth  time.Time `sheet:"Birthdate"`
	Weight int64     `sheet:"Weight"`
}

type MedicineCfg struct {
	Posology []*PosologyEntry
}
