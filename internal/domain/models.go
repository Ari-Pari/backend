package domain

import "time"

type Translation struct {
	EngName string `json:"engName"`
	RuName  string `json:"ruName"`
	ArmName string `json:"armName"`
}

type Genre string

const (
	War         Genre = "WAR"
	Road        Genre = "ROAD"
	Cult        Genre = "CULT"
	Lyrical     Genre = "LYRICAL"
	Reverse     Genre = "REVERSE"
	Ritual      Genre = "RITUAL"
	Community   Genre = "COMMUNITY"
	Hunting     Genre = "HUNTING"
	Pilgrimage  Genre = "PILGRIMAGE"
	Memorable   Genre = "MEMORABLE"
	Memorial    Genre = "MEMORIAL"
	Funeral     Genre = "FUNERAL"
	Festive     Genre = "FESTIVE"
	Wedding     Genre = "WEDDING"
	Matchmakers Genre = "MATCHMAKERS"
	Labor       Genre = "LABOR"
	Amulet      Genre = "AMULET"
)

type HoldingType string

const (
	Free         HoldingType = "FREE"
	LittleFinger HoldingType = "LITTLE_FINGER"
	Palm         HoldingType = "PALM"
	Crossed      HoldingType = "CROSSED"
	Back         HoldingType = "BACK"
	Belt         HoldingType = "BELT"
	Shoulder     HoldingType = "SHOULDER"
	Dagger       HoldingType = "DAGGER"
	Whip         HoldingType = "WHIP"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
	Multi  Gender = "MULTI"
)

type Region struct {
	Id   int64
	Name Translation
}

type DanceShort struct {
	Id           int64
	Name         Translation
	NameKey      string
	Paces        []int32
	HoldingTypes []HoldingType
	Gender       Gender
	Complexity   *int32
	Genres       []Genre
	RegionIds    []int64
	DeletedAt    *time.Time
}

type SongShort struct {
	Id        int64
	Name      Translation
	NameKey   string
	DanceIds  []int64
	ArtistIds []int64
}

type VideoType string

const (
	Lesson VideoType = "LESSON"
	Video  VideoType = "VIDEO"
	Source VideoType = "SOURCE"
)

type VideoShort struct {
	Id       *int64
	Name     Translation
	Type     VideoType
	NameKey  string
	Link     string
	DanceIds []int64
}

type ArtistShort struct {
	Id        int64
	Name      Translation
	NameKey   string
	Url       string
	DeletedAt *time.Time
}
