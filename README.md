# dalg
一个Golang 的DATA ACCESS LAYER代码生成器.
这是个粗陋的版本，下一个版会结合PHAROSNET/DALC进行重构优化。

特性：

- 使用代码生成的方式代替反射，这样在运行时的损耗相比小很多。

- 对于查询，使用回调的方式减少了多余的循环次数。在通常类似的框架中，当查询某一个列表时，会返回一个Data List，而从rows转化成Data List的循环是被封装了的，往往之后的业务代码中，还会在循环Data List进行其他操作。而dalg是在从rows转化成Data List的循环中织入了函数体，这个函数体是执行业务逻辑代码的，那么就可以省去一个不必要的循环消耗。
- 支持Mac，Windows（需要admin权限，且请关闭360这类东西，因为要取系统环境变量），Linux（待测）。
- 支持额外第三方包类型，且会自动识别，但是要写全包名，如 `mapType="xx/xx.NullXxx"`。
- 增加了NullTime类型。

## 安装

```
go get -u github.com/pharosnet/dalg
```

## 使用

```
dalg [描述文件] [输出目录]
```

案例

```
dalg ./def.xml $GOPATH/src/project
```

 def.xml 指定了输出的包名，默认是dal，所以生成的代码在$GOPATH/src/project/dal中。

描述文件是一个xml文件，根是db，db包含package（输出的包名）和dialect（方言，当前支持postgres），在db下只有interface节点，interface是描述对象的，也就是对应输出的struct。interface的属性有class（类别，可以有table，enum，json，view。table是表，会生成CRUD代码，enum和json是扩张类型，具体见下面案例，view是查询时图，不会生成CRUD代码），schema，name（表名），mapName（生成代码的名称）和type（数据库类型）。在xml中，以map开头的属性都是输出用，不带map的都是描述数据库用。

注意，如果查询结果是计数类型时，类似count，那么按照逻辑，需要建立一个view类型的interface，因query的结果元类型为interface，且使用call back方式，所以query并不需要指定返回list还是one。

def.xml的案例如下：

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<db
        package="dal"
        dialect="postgres"
        tablespace=""
        owner=""
        >
    <interface class="table" schema="public" name="user" mapName="User" >
        <column name="id" type="varchar(64)" mapName="Id" mapType="sql.NullString" pk="true" increment="false" />
        <column name="name" type="varchar(255)" mapName="Name" mapType="sql.NullString" />
        <column name="age" type="numeric(18)" mapName="Age" mapType="sql.NullInt64" />
        <column name="sex" type="boolean" mapName="Sex" mapType="Sex" />
        <column name="money" type="numeric(18,2)" mapName="Money" mapType="sql.NullFloat64" />
        <column name="info" type="jsonb" mapName="Info" mapType="UserInfo" json="true" />
        <column name="create_by" type="varchar(255)" mapName="CreateBy" mapType="sql.NullString"  />
        <column name="create_time" type="timestamp" mapName="CreateTime" mapType="NullTime"  />
        <column name="modify_by" type="varchar(255)" mapName="ModifyBy" mapType="sql.NullString"  />
        <column name="modify_time" type="timestamp" mapName="ModifyTime" mapType="NullTime" />
        <column name="version" type="bigint" mapName="Version" mapType="sql.NullInt64" version="true" />
        <index name="name_age" type="btree" columns="name, age" sortOrder="desc" unique="false" ops=""/>
        <index name="create_time" type="btree" columns="create_time" sortOrder="desc" unique="false" ops=""/>
        <query mapName="List" result="" > <!-- result 支持 one, list, int64, string, bool, float64, 默认为list -->
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
        <field mapName="Id" mapType="string" />
        <field mapName="Age" mapType="int64" />
    </interface>
    <interface class="enum" mapName="Sex" mapType="string" type="bool"  >
        <option value="true" mapValue="MALE" default="true" />
        <option value="false" mapValue="FEMALE" />
    </interface>
    <interface class="view" mapName="UserCount">
    	<column name="count" type="numeric(18)" mapName="Count" mapType="sql.NullInt64" />
        <query mapName="Normal" >
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

### TODO

- [ ] 支持mysql
- [ ] 支持oracle
- [x] 支持DDL

## License

GNU GENERAL PUBLIC LICENSE
