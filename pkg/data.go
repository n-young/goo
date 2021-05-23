package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Data interface - built from YAML data files.
type Data interface {
	getChild(key string) (Data, error)
	getValue() (string, error)
	getList() ([]Data, error)
}

// DataNode - means there is more data below.
type DataNode map[string]Data
func (n DataNode) getChild(key string) (Data, error) {
	ret, ok := n[key]
	if ok {
		return ret, nil
	}
	return nil, GenericError{"Unknown data key: " + key}
}
func (n DataNode) getValue() (string, error) {
	return "", GenericError{"Called getValue on a Node."}
}
func (n DataNode) getList() ([]Data, error) {
	return nil, GenericError{"Called getList on a Node."}
}
func (n DataNode) setTitle(title string) {
	n["title"] = DataLeaf(title)
}
func (n DataNode) setGlobal(data Data) {
	n["global"] = data
}

// DataLeaf - means there is no more data below.
type DataLeaf string
func (l DataLeaf) getChild(key string) (Data, error) {
	return nil, GenericError{"Called getChild on a Leaf."}
}
func (l DataLeaf) getValue() (string, error) {
	return string(l), nil
}
func (l DataLeaf) getList() ([]Data, error) {
	return nil, GenericError{"Called getList on a Leaf."}
}

// DataLeaf - means there is no more data below.
type DataList []Data
func (l DataList) getChild(key string) (Data, error) {
	return nil, GenericError{"Called getChild on a List."}
}
func (l DataList) getValue() (string, error) {
	return "", GenericError{"Called getValue on a List."}
}
func (l DataList) getList() ([]Data, error) {
	return []Data(l), nil
}


// Cast a generic unmarshal to the Data format above.
func CastData(m map[interface{}]interface{}) (Data, error) {
	var err error
	data := make(map[string]Data)
	for k, v := range m {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			data[k.(string)], err = CastData(v)
			Check(err)
		case string:
			data[k.(string)] = DataLeaf(v)
		case []interface{}:
			newList := DataList(make([]Data, len(v)))
			for kk, vv := range v {
				newList[kk], err = CastData(vv.(map[interface{}]interface{}))
				Check(err)
			}
			data[k.(string)] = newList
		default:
			fmt.Printf("map = %+v\n", m)
			fmt.Printf("data = %T\n", v)
			return nil, GenericError{"Malformed Data."}
		}
	}
	return DataNode(data), nil
}

// Gets data from a file in the format above.
func GetDataFromFile(filename string) Data {
	// Reads data out into a string.
	bytes, io_err := ioutil.ReadFile(filename)
	Check(io_err)

	// Unmarshal into m.
	m := make(map[interface{}]interface{})
	y_err := yaml.Unmarshal(bytes, &m)
	Check(y_err)

	// Build Data object.
	data, d_err := CastData(m)
	Check(d_err)
	return data
}

// Converts a map of data sources to the format above.
func GetData(m map[string]string) Data {
	data := make(DataNode)
	for k, v := range m {
		data[k] = GetDataFromFile(v)
	}
	return data
}

// Gets a specific datapoint (home.title) from data.
func ExtractData(payload string, data Data) (string, error) {
	fields := strings.Split(payload, ".")
	curr := data
	var c_err error
	for _, field := range fields {
		curr, c_err = curr.getChild(field)
		Check(c_err)
	}
	return curr.getValue()
}

// Gets a specific DataList from data.
func ExtractList(payload string, data Data) ([]Data, error) {
	fields := strings.Split(payload, ".")
	curr := data
	var c_err error
	for _, field := range fields {
		curr, c_err = curr.getChild(field)
		Check(c_err)
	}
	return curr.getList()
}
