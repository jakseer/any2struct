# Any2struct
Convert anything to Go Struct

## Feature
Parse string with certain encoding format into Go Struct with specified tag.

Support parsing encoding format
- JSON
- SQL
- YAML

Support output Go Struct Tag
- JSON
- Gorm
- YAML

## Usage
### cli mode
Parse and generate data in cli mode.

```text
Usage:
  any2struct generate [flags]

Flags:
  -d, --data-type string   input data type, support: yaml,json,sql
  -h, --help               help for generate
  -i, --input string       input file path
  -o, --output string      generated file path (default "./go-struct.go")
  -t, --tags string        generated struct tags, support: json,gorm,yaml
```
#### Example
```shell
any2struct generate -d json "{\"a\":1,\"b\":[1,2,\"b\"],\"c\":{\"c1\":1,\"c2\":[1,2]}}" -t=json,yaml
```
The stdout is:
```text
type Untitled struct { 
    A int64 `json:"a" yaml:"a"` 
    B []interface `json:"b" yaml:"b"` 
    C C `json:"c" yaml:"c"` 
}

type C struct { 
    C1 int64 `json:"c_1" yaml:"c_1"` 
    C2 []interface `json:"c_2" yaml:"c_2"` 
}
```

### package mode
Using functions in the package.

```
func (*Convertor) Convert(input string, decodeType string, encodeTypes []string) (string, error)
```
Parse `input` which is `decodeType` type into Go Struct with `encodeTypes` tags.

#### Example
```go
sql := "CREATE TABLE `USER`(" + "\n"
sql = sql + "`id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'primary key'," + "\n"
sql = sql + "`name`    VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'name'," + "\n"
sql = sql + "`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time'," + "\n"
sql = sql + "`deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'delete time'," + "\n"
sql = sql + "PRIMARY KEY(`id`)" + "\n"
sql = sql + ")ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='user table';" + "\n"

// parse json to go struct with tag gorm and json
out, err := NewConvertor().Convert(sql, DecodeTypeJson, []string{EncodeTypeJson, EncodeTypeGorm})
if err != nil {
    fmt.Println(out)
}
```

The stdout is :
```text
type USER struct { 
    Id int `json:"id" gorm:"column:id"` // primary key
    Name string `json:"name" gorm:"column:name"` // name
    CreatedAt unknown `json:"created_at" gorm:"column:created_at"` // create time
    DeletedAt unknown `json:"deleted_at" gorm:"column:deleted_at"` // delete time
}
```

## License
Released under the MIT license [MIT license](https://github.com/jakseer/any2struct/blob/master/LICENSE)