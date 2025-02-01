package models

import "time"

type Person string

const (
	Aline  Person = "Aline"
	Dorian Person = "Dorian"
	Zaya   Person = "Zaya"
	Azel   Person = "Azel"
)

type Medicine string

const (
	ChildrenIbuprofen   Medicine = "ChildrenIbuprofen"
	InfantAcetaminophen Medicine = "InfantAcetaminophen"
)

type PosologyEntry struct {
	OlderThan time.Duration
	Interval  time.Duration
	Quantity  string
}

type PersonCfg struct {
	Birth    time.Time
	NextDose map[Medicine]time.Time
}

type MedicineCfg struct {
	Posology []PosologyEntry
}
