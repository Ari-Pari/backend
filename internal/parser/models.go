package parser

type NameDto struct {
	EngName string `json:"en"`
	RuName  string `json:"ru"`
	ArmName string `json:"hy"`
}

type StateDto struct {
	Id   int64   `json:"id"`
	Name NameDto `json:"name"`
}

type HoldingTypeDto string

const (
	Free         HoldingTypeDto = "AZAT"
	LittleFinger HoldingTypeDto = "CHKUYT"
	Palm         HoldingTypeDto = "AP"
	Crossed      HoldingTypeDto = "KHACH"
	Back         HoldingTypeDto = "MEJQ"
	Belt         HoldingTypeDto = "GOTI"
	Shoulder     HoldingTypeDto = "US"
	Dagger       HoldingTypeDto = "ZENQ"
	Whip         HoldingTypeDto = "MTRAK"
)

type GenderDto string

const (
	Male   GenderDto = "BOY"
	Female GenderDto = "GIRL"
	Multi  GenderDto = "MULTI"
)

type TypeDto string

const (
	Active TypeDto = "ACTIVE"
	Extra  TypeDto = "EXTRA"
)

type GenreDto string

const (
	War         GenreDto = "WAR"
	Road        GenreDto = "ROAD"
	Cult        GenreDto = "CULT"
	Lyrical     GenreDto = "LYRICAL"
	Reverse     GenreDto = "REVERSE"
	Ritual      GenreDto = "RITUAL"
	Community   GenreDto = "COMMUNITY"
	Hunting     GenreDto = "HUNTING"
	Pilgrimage  GenreDto = "PILGRIMAGE"
	Memorable   GenreDto = "MEMORABLE"
	Memorial    GenreDto = "MEMORIAL"
	Funeral     GenreDto = "FUNERAL"
	Festive     GenreDto = "FESTIVE"
	Wedding     GenreDto = "WEDDING"
	Matchmakers GenreDto = "MATCHMAKERS"
	Labor       GenreDto = "LABOR"
	Amulet      GenreDto = "AMULET"
)

type DanceDto struct {
	Id           int64            `json:"id"`
	Name         NameDto          `json:"name"`
	Type         TypeDto          `json:"type"`
	NameKey      string           `json:"nameKey"`
	Temps        []int32          `json:"temps"`
	HoldingTypes []HoldingTypeDto `json:"holdingTypes"`
	Gender       GenderDto        `json:"gender"`
	Difficult    *int32           `json:"difficult"`
	Genres       []GenreDto       `json:"genres"`
	StateIds     []int64          `json:"states"`
}

type MusicDto struct {
	Id       int64   `json:"id"`
	Name     NameDto `json:"name"`
	NameKey  string  `json:"nameKey"`
	DanceIds []int64 `json:"danceIds"`
	Type     TypeDto `json:"type"`
	Artists  []int64 `json:"groupIds"`
}

type VideoTypeDto string

const (
	VideoTypeLesson VideoTypeDto = "LESSON"
	VideoTypeVideo  VideoTypeDto = "VIDEO"
	VideoTypeSource VideoTypeDto = "SOURCE"
)

type VideoDto struct {
	Name     NameDto      `json:"name"`
	Url      string       `json:"url"`
	Type     VideoTypeDto `json:"type"`
	DanceIds []int64      `json:"danceIds"`
}

type ArtistDto struct {
	Id   int64   `json:"id"`
	Name NameDto `json:"name"`
	Type TypeDto `json:"type"`
	Url  string  `json:"urlInsta"`
}
