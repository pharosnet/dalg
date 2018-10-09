package def

import (
	"encoding/xml"
)

type Db struct {
	XMLName    xml.Name `xml:"db"`
	Package    string   `xml:"package,attr"`
	Dialect    string   `xml:"dialect,attr"`
	DDL        bool     `xml:"ddl,attr"`
	Tablespace string   `xml:"tablespace,attr,omitempty"`
	Owner      string   `xml:"owner,attr,omitempty"`
	//Driver           string      `xml:"driver,attr"`
	//ConnMaxLifetime  string      `xml:"connMaxLifetime,attr,omitempty"`
	//ConnMaxLifetimeV time.Duration      `xml:"-"`
	//MaxIdleConns     int64       `xml:"maxIdleConns,attr,omitempty"`
	//MaxOpenConns     int64       `xml:"maxOpenConns,attr,omitempty"`
	Interfaces []Interface `xml:"interface"`
}

type Interface struct {
	Class         string       `xml:"class,attr"`
	Schema        string       `xml:"schema,attr,omitempty"`
	Name          string       `xml:"name,attr"`
	MapName       string       `xml:"mapName,attr"`
	Type          string       `xml:"type,attr,omitempty"`
	MapType       string       `xml:"mapType,attr,omitempty"`
	Columns       []Column     `xml:"column,omitempty"`
	Indexes       []Index      `xml:"index,omitempty"`
	Options       []EnumOption `xml:"option,omitempty"`
	Fields        []Field      `xml:"field,omitempty"`
	Queries       []Query      `xml:"query,omitempty"`
	Package       string       `xml:"-"`
	Imports       []string     `xml:"-"`
	Dialect       string       `xml:"-"`
	Pks           []Column     `xml:"-"`
	Tablespace    string       `xml:"-"`
	Owner         string       `xml:"-"`
	CommonColumns []Column     `xml:"-"`
	Version       Column       `xml:"-"`
}

type Column struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	MapName   string `xml:"mapName,attr"`
	MapType   string `xml:"mapType,attr"`
	Pk        bool   `xml:"pk,attr,omitempty"`
	Increment bool   `xml:"increment,attr,omitempty"`
	Version   bool   `xml:"version,attr,omitempty"`
	Json      bool   `xml:"json,attr,omitempty"`
	Xml       bool   `xml:"xml,attr,omitempty"`
	NotNull   bool   `xml:"notNull,attr,omitempty"`
	Default   string `xml:"default,attr,omitempty"`
}

type Index struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	Unique    bool   `xml:"unique,attr,omitempty"`
	Columns   string `xml:"columns,attr"`
	SortOrder string `xml:"sortOrder,attr,omitempty"`
	Ops       string `xml:"ops,attr,omitempty"`
}

type EnumOption struct {
	Value    string `xml:"value,attr"`
	MapValue string `xml:"mapValue,attr"`
	Default  bool   `xml:"default,attr"`
}

type Field struct {
	Name    string `xml:"name,attr"`
	MapName string `xml:"mapName,attr"`
	MapType string `xml:"mapType,attr"`
}

type Query struct {
	MapName string     `xml:"mapName,attr"`
	Args    []QueryArg `xml:"arg"`
	Result  string     `xml:"result,attr"`
	Sql     Sql        `xml:"sql"`
}

type Sql struct {
	Value string `xml:",cdata"`
}

type QueryArg struct {
	MapName string `xml:"mapName,attr"`
	MapType string `xml:"mapType,attr"`
}
