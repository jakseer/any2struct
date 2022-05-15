# Any2struct
Convert anything to Go Struct

## Feature
Parse string with certain encoding format into Go Struct with specified tag.

Support parsing encoding format
- JSON
- SQL
- XML

Support output Go Struct Tag
- JSON
- Gorm
- XML

## Usage
```
func (*Convertor) Convert(input string, decodeType string, encodeTypes []string) (string, error)
```
Parse `input` which is `decodeType` type into Go Struct with `encodeTypes` tags.

### Example
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