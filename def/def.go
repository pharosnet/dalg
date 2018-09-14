package def

import (
	"encoding/xml"
)

type Db struct {
	XMLName          xml.Name    `xml:"db"`
	Package          string      `xml:"package,attr"`
	Dialect          string      `xml:"dialect,attr"`
	//Driver           string      `xml:"driver,attr"`
	//ConnMaxLifetime  string      `xml:"connMaxLifetime,attr,omitempty"`
	//ConnMaxLifetimeV time.Duration      `xml:"-"`
	//MaxIdleConns     int64       `xml:"maxIdleConns,attr,omitempty"`
	//MaxOpenConns     int64       `xml:"maxOpenConns,attr,omitempty"`
	EnableNil bool `xml:"enableNil"`
	Interfaces       []Interface `xml:"interfaces"`
}

type Interface struct {
	Type      string    `xml:"type,attr"`
	Schema    string    `xml:"schema,attr,omitempty"`
	Name      string    `xml:"name,attr"`
	MapName   string    `xml:"mapName,attr"`
	Columns   []Column  `xml:"columns"`
	ExtraType ExtraType `xml:"extraType,omitempty"`
	Queries   []Query   `xml:"queries,omitempty"`
	Package   string 	`xml:"-"`
	EnableNil bool `xml:"-"`
	Dialect          string `xml:"-"`
	Pks []Column `xml:"-"`
	PkNum int64 `xml:"-"`
	CommonColumns []Column `xml:"-"`
	Version Column `xml:"-"`
	InsertSql string `xml:"-"`
	UpdateSql string `xml:"-"`
	DeleteSql string `xml:"-"`
	GetOneSql string `xml:"-"`
}

type Column struct {
	Name       string `xml:"name,attr"`
	Type       string `xml:"type,attr"`
	MapName    string `xml:"mapName,attr"`
	MapType    string `xml:"mapType,attr"`
	Pk         bool   `xml:"pk,attr,omitempty"`
	DbIncrement bool `xml:"dbIncrement"`
	CreateBy   bool   `xml:"createBy,attr,omitempty"`
	CreateTime bool   `xml:"createTime,attr,omitempty"`
	ModifyBy   bool   `xml:"modifyBy,attr,omitempty"`
	ModifyTime bool   `xml:"modifyTime,attr,omitempty"`
	Version    bool   `xml:"version,attr,omitempty"`
	EnableNil  bool   `xml:"enableNil,attr,omitempty"`
}

type ExtraType struct {
	EnumInterfaces []Enum `xml:"enum,omitempty"`
	JsonInterfaces []Json `xml:"json,omitempty"`
	Packages []string `xml:"package,omitempty"`
}

type Enum struct {
	Id      string       `xml:"id,attr"`
	Type 	string `xml:"-"`
	MapType string       `xml:"mapType,attr"`
	Options []EnumOption `xml:"options"`
}

type EnumOption struct {
	Value    string `xml:"value,attr"`
	MapValue string `xml:"mapValue,attr"`
}

type Json struct {
	Id     string      `xml:"id,attr"`
	Fields []JsonField `xml:"fields"`
}

type JsonField struct {
	MapName string `xml:"mapName,attr"`
	MapType string `xml:"mapType,attr"`
}

type Query struct {
	Args   []QueryArg `xml:"args"`
	Result string     `xml:"result,attr"`
	Sql    Sql        `xml:"sql"`
}

type Sql struct {
	Value string `xml:",cdata"`
}

type QueryArg struct {
	MapName string `xml:"mapName,attr"`
	MapType string `xml:"mapType,attr"`
}
