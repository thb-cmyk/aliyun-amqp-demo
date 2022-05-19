package voltage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TableNode struct {
	ParentDir string
	TableName string

	Column map[string]string

	TimeStamp time.Time
}

type DataEntry struct {
	Id string

	Raw       []byte
	Structure interface{}

	key interface{}

	Location *TableNode

	Timestamp time.Time

	Next *DataEntry
}

const USERNAME string = ""
const PASSWORD string = ""

const MAXTABLE int = 10
const MAXDATAENTRY int = 20

var global_table_list [MAXTABLE]*TableNode

var global_data_entry [MAXDATAENTRY]*DataEntry
var global_data_read_pointer int = 0
var global_data_write_pointer int = 0
var global_data_used int = 0

func DataInformationGen(raw []byte, id string, try bool) (*DataEntry, bool) {
	var entry *DataEntry = new(DataEntry)

	if raw == nil {
		return nil, true
	}

	err := json.Unmarshal(raw, &entry.Structure)
	if err != nil {
		fmt.Printf("location 002 error, unmarshal!\n\r")
		return nil, false
	}

	entry.Id = id
	entry.Location = nil
	entry.Raw = raw
	entry.Timestamp = time.Now()
	entry.Next = nil
	entry.key = entry.Structure.(map[string]interface{})["table"]

	index := Table_select(entry, 1)
	if index == -1 {
		fmt.Printf("location 003 error, table select!\n\r")
		return nil, false
	}
	entry.Location = global_table_list[index]

	for {
		ok := Data_entry_add(entry, true)
		if ok {
			return entry, true
		} else if !try {
			return entry, true
		} else {
			if ok := Data_entry_remove(nil, 1); !ok {
				fmt.Printf("The fatal error is occurse, the datainformationgen function!\n\r")
				return nil, false
			}
		}
	}
}

func DataInsert(entry *DataEntry) bool {
	if entry == nil {
		fmt.Printf("The data is nil, please check!\n\r")
		return false
	}

	fmt.Printf("databasename:%s, tablename: %s\n\r", entry.Location.ParentDir, entry.Location.TableName)
	if entry.Structure == nil {
		return false
	}
	map_data := entry.Structure.(map[string]interface{})
	for k, v := range map_data {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	return true
}

func Data_entry_add(entry *DataEntry, try bool) bool {

	if global_data_read_pointer == global_data_write_pointer && global_data_used != 0 {
		if try {
			global_data_entry[global_data_write_pointer] = entry
			global_data_write_pointer = (global_data_write_pointer + 1) % MAXDATAENTRY
			global_data_read_pointer = global_data_write_pointer
			return true
		} else {
			return false
		}
	} else {
		global_data_entry[global_data_write_pointer] = entry
		global_data_write_pointer = (global_data_write_pointer + 1) % MAXDATAENTRY
		global_data_write_pointer++
		global_data_used++
		return true
	}
}

func Data_entry_remove(entry *DataEntry, try int) bool {
	if global_data_read_pointer == global_data_write_pointer && global_data_used == 0 {
		return false
	} else {
		global_data_entry[global_data_read_pointer] = nil
		global_data_read_pointer = (global_data_read_pointer + 1) % MAXDATAENTRY
		global_data_used--
		return true
	}
}

func Column_generate(entry *DataEntry) map[string]string {

	var column map[string]string = make(map[string]string, 100)

	for k, v := range entry.Structure.(map[string]interface{}) {
		switch vv := v.(type) {
		case string:
			column[k] = "string"
		case float64:
			column[k] = "float64"
		case []interface{}:
			column[k] = "arry"
		default:
			_ = vv
		}
	}

	return column
}

func Table_select_empty() int {
	var index int

	for index = 0; index < MAXTABLE && global_table_list[index] == nil; index++ {
		return index
	}
	return -1
}

func Table_create(parentdir string, tablename string, column map[string]string) int {
	index := Table_select_empty()
	if index == -1 {
		return -1
	}
	table := new(TableNode)
	table.ParentDir = parentdir
	table.TableName = tablename
	table.TimeStamp = time.Now()
	table.Column = make(map[string]string, len(column))
	table.Column = column

	/* The following code to create a database table to store the receveing data */
	os.Create(tablename)

	global_table_list[index] = table
	return index
}

func Table_select(entry *DataEntry, try int) int {

	if entry == nil {
		return -1
	}

	key := entry.key.(string)

	for i := 0; i < MAXTABLE && global_table_list[i] != nil; i++ {
		if key == global_table_list[i].TableName {
			return i
		}
	}
	if try == 1 {
		column := Column_generate(entry)
		if column == nil {
			return -1
		}
		index := Table_create("defautdatabase", key, column)
		return index
	}
	return -1
}

func Table_delete(parentdir string, tablename string) bool {
	for i := 0; i < MAXTABLE; i++ {
		if global_table_list[i] == nil {
			continue
		}
		if global_table_list[i].ParentDir == parentdir && global_table_list[i].TableName == tablename {
			global_table_list[i] = nil
			return true
		}
	}
	return false
}

func main() {

}
