package voltage

import (
	"fmt"
	"testing"
)

func TestDataInformationGen(t *testing.T) {
	raw := []byte("{\"voltage\":4, \"table\":\"voltagetable\"}")

	dataentry, ok := DataInformationGen(raw, "voltage", 1)
	if !ok {
		fmt.Printf("datainformationgen is failed!\n\r")
	} else if dataentry == nil {
		fmt.Printf("dataentry is nil!\n\r")
	} else {
		fmt.Printf("id:%s\n\r", dataentry.Id)
		fmt.Printf("Raw:%s\n\r", string(dataentry.Raw))
		fmt.Printf("key:%s\n\r", dataentry.key.(string))
	}
}

func TestDataInsert(t *testing.T) {
	raw := []byte("{\"voltage\":4, \"table\":\"voltagetable\"}")

	dataentry, ok := DataInformationGen(raw, "voltage", 1)
	if !ok {
		fmt.Printf("datainformationgen is failed!\n\r")
	} else if dataentry == nil {
		fmt.Printf("dataentry is nil!\n\r")
	} else {
		fmt.Printf("id:%s\n\r", dataentry.Id)
		fmt.Printf("Raw:%s\n\r", string(dataentry.Raw))
		fmt.Printf("key:%s\n\r", dataentry.key.(string))
	}

	ok = DataInsert(dataentry)
	if !ok {
		fmt.Printf("datainsert false!\n\r")
	}
}

func TestTableCreate(t *testing.T) {
	raw := []byte("{\"voltage\":4, \"table\":\"voltagetable\"}")

	dataentry, ok := DataInformationGen(raw, "voltage", 1)
	if !ok {
		fmt.Printf("datainformationgen is failed!\n\r")
	} else if dataentry == nil {
		fmt.Printf("dataentry is nil!\n\r")
	} else {
		fmt.Printf("id:%s\n\r", dataentry.Id)
		fmt.Printf("Raw:%s\n\r", string(dataentry.Raw))
		fmt.Printf("key:%s\n\r", dataentry.key.(string))
	}

	ok = DataInsert(dataentry)
	if !ok {
		fmt.Printf("datainsert false!\n\r")
	}

	index := Table_select(dataentry, 1)
	if index == -1 {
		fmt.Printf("table_select failed!\n\r")
	} else {
		tablenode := global_table_list[index]
		fmt.Printf("ParentDir:%s\n\r", tablenode.ParentDir)
		fmt.Printf("TableName:%s\n\r", tablenode.TableName)
		for k, v := range tablenode.Column {
			fmt.Printf("key:%s, value:%s\n\r", k, v)
		}
	}

	if Table_delete("defautdatabase", dataentry.key.(string)) {
		if global_table_list[index] == nil {
			fmt.Printf("table_delete is success!\n\r")
		} else {
			fmt.Printf("table delete is failed!\n\r")

		}
	}
}
