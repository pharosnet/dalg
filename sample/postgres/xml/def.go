package main

import (
	"encoding/xml"
	"fmt"
	"github.com/pharosnet/dalg/sample/postgres/xml/def"
)

func main()  {

	root := def.Db{}
	root.Dialect = "mysql"
	root.Driver = "xx"
	root.ConnMaxLifetime = "2s"
	root.MaxIdleConns = 2
	root.MaxOpenConns = 2
	interfaces := make([]def.Interface, 0, 1)
	i1 := def.Interface{}
	i1.Type = "table"
	i1.Schema = "dbo"
	i1.Name = "user"
	i1.MapName = "User"
	cols := make([]def.Column, 0, 1)
	cols = append(cols,
		def.Column{
			Name:"id",
			Type:"string",
			MapName:"Id",
			MapType:"string",
			Pk:true,
		},
		def.Column{
			Name:"name",
			Type:"string",
			MapName:"Name",
			MapType:"string",
		},
		def.Column{
			Name:"age",
			Type:"int",
			MapName:"Age",
			MapType:"int64",
		},
		def.Column{
			Name:"sex",
			Type:"bool",
			MapName:"Sex",
			MapType:"#Sex",
		},
		def.Column{
			Name:"create_by",
			Type:"string",
			MapName:"CreateBy",
			MapType:"string",
			CreateBy:true,
		},
		def.Column{
			Name:"create_time",
			Type:"datetime",
			MapName:"CreateTime",
			MapType:"time.Time",
			CreateTime:true,
		},
		def.Column{
			Name:"modify_by",
			Type:"string",
			MapName:"ModifyBy",
			MapType:"string",
			ModifyBy:true,
		},
		def.Column{
			Name:"modify_time",
			Type:"datetime",
			MapName:"ModifyTime",
			MapType:"time.Time",
			ModifyTime:true,
		},
		def.Column{
			Name:"version",
			Type:"int",
			MapName:"Version",
			MapType:"Int64",
			Version:true,
		},
	)
	i1.Columns = cols
	i1.ExtraType.EnumInterfaces = make([]def.Enum, 0, 1)

	i1.ExtraType.EnumInterfaces = append(i1.ExtraType.EnumInterfaces,
		def.Enum{
			Id:"e1",
			MapType:"string",
			Options:[]def.EnumOption{
				def.EnumOption{
					Value:"1",
					MapValue:"string",
				},
				def.EnumOption{
					Value:"2",
					MapValue:"string",
				},
			},
		},
	)

	i1.ExtraType.JsonInterfaces = []def.Json{
		{
			Id:"json",
			Fields:[]def.JsonField{
				{MapName:"id", MapType:"string"},
				{MapName:"id", MapType:"string"},
				{MapName:"id", MapType:"string"},
			},
		},
	}
	i1.ExtraType.Packages = []string{"11", "22"}


	i1.Queries = []def.Query{
		{
			Args:[]def.QueryArg{
				{MapName:"id", MapType:"string"},
			},
			Sql:def.Sql{"select * from user"},
		},
	}
	interfaces = append(interfaces, i1)
	root.Interfaces = interfaces


	bb, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bb))
}
