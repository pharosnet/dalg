# Database access layer generator
This is a dal(database access layer) code generator for Golang.

Features:

- Using code generating instead of reflect way, so there is no reflect costs.

- In query function code, using function frame instead of directed returns. so we can write an O(n) costed code to achieve goals.
- Supports Linux, Windows and Mac.
- Supports customize type.
- DDL 

## Install

```
// in your project
go get -u github.com/pharosnet/dalc
// install dalg program
go get -u github.com/pharosnet/dalg
```

## Usage

```
dalg [definition file, it is a xml type file] [generated code file dir]
```

## Example:

command: 

```
dalg ./def.xml $GOPATH/src/project
```

def.xml: 

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<db
        package="dal" // optional, default value is dal. go code package name, e.g.: $GOPATH/src/project/dal
        dialect="postgres" // required, postgres, mysql, oracle
        tablespace="" // postgres tablespace
        owner="" // postgres owner
    	ddl="true" // when true, it will generate ddl file.
        >
    <!-- interface class can be table, view, enum and json -->
    <!-- schema is optional -->
    <!-- name is required when class is table, and it is table name  -->
    <!-- mapName is required and uniqued, it is go struct name,  -->
    <interface class="table" schema="public" name="user" mapName="User" >
        <!-- when class is table or view, column is required -->
        <column name="id" type="varchar(64)" mapName="Id" mapType="sql.NullString" pk="true" increment="false" />
        <column name="name" type="varchar(255)" mapName="Name" mapType="sql.NullString" />
        <column name="age" type="numeric(18)" mapName="Age" mapType="sql.NullInt64" />
        <column name="sex" type="boolean" mapName="Sex" mapType="Sex" />
        <column name="money" type="numeric(18,2)" mapName="Money" mapType="sql.NullFloat64" />
        <!-- mapType can be interface mapName -->
        <column name="info" type="jsonb" mapName="Info" mapType="UserInfo" json="true" />
        <column name="create_by" type="varchar(255)" mapName="CreateBy" mapType="sql.NullString"  />
        <column name="create_time" type="timestamp" mapName="CreateTime" mapType="NullTime"  />
        <column name="modify_by" type="varchar(255)" mapName="ModifyBy" mapType="sql.NullString"  />
        <column name="modify_time" type="timestamp" mapName="ModifyTime" mapType="NullTime" />
        <column name="version" type="bigint" mapName="Version" mapType="sql.NullInt64" version="true" />
        <index name="name_age" type="btree" columns="name, age" sortOrder="desc" unique="false" ops=""/>
        <index name="create_time" type="btree" columns="create_time" sortOrder="desc" unique="false" ops=""/>
        <query mapName="List" result="list" > <!-- result can be one, list, int64, string, bool, float64, default is list -->
            <arg mapName="limit" mapType="int64" />
            <arg mapName="offset" mapType="int64" />
            <sql>
                <![CDATA[
                SELECT "id", "name", "age", "sex", "money", "info", "create_by", "create_time", "modify_by", "modify_time", "version"
                FROM "user"
                LIMIT $1 OFFSET $2
                ]]>
            </sql>
        </query>
    </interface>
    <interface class="json" mapName="UserInfo" >
        <!-- when class is table or view, field is required -->
        <field mapName="Id" mapType="string" />
        <field mapName="Age" mapType="int64" />
    </interface>
    <interface class="enum" mapName="Sex" mapType="string" type="bool"  >
        <!-- when class is table or view, option is required -->
        <option value="true" mapValue="MALE" default="true" />
        <option value="false" mapValue="FEMALE" />
    </interface>
    <interface class="view" mapName="UserCount">
    	<column name="count" type="numeric(18)" mapName="Count" mapType="sql.NullInt64" />
        <query mapName="Normal" result="int64" >
            <arg mapName="limit" mapType="int64" />
            <arg mapName="offset" mapType="int64" />
            <sql>
                <![CDATA[
                SELECT COUNT("id") AS "COUNT"
                FROM "user"
                LIMIT $1 OFFSET $2
                ]]>
            </sql>
        </query>
    </interface>

</db>
```

To see full example at the example folder.

### TODO

- [ ] support mysql
- [ ] support oracle
- [x] 支持DDL

## License

GNU GENERAL PUBLIC LICENSE
