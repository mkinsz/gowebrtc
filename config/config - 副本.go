package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/deepch/vdk/av"
)

// Cfg config
var Cfg *Configs

// Configs configs
type Configs struct {
	Server  Server            `json:"server"`
	Streams map[string]Stream `json:"streams"`
}

// Server struct name ip port with tag
type Server struct {
	HTTPPort string `json:"http_port"`
}

// Stream stream struct
type Stream struct {
	URL    string `json:"url"`
	Status bool   `json:"status"`
	Codecs []av.CodecData
	Cl     map[string]viwer
}

type viwer struct {
	c chan av.Packet
}

func init() {
	CFG = load()
}

func load() *Configs {
	t := Configs{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("read config file err: ", err.Error())
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatalln(err)
	}

	for i, v := range t.Streams {
		v.Cl = make(map[string]viwer)
		t.Streams[i] = v
	}

	return &t
}

// T for test
type T struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

func serialize() {
	t := T{}
	t.Name = "Demo-Json"
	t.IP = "127.0.0.1"
	t.Port = 8080

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Marshal err: ", err.Error())
		return
	}
	fmt.Println("Marshal Json: ", string(b))
}

func serializeMap() {
	t := make(map[string]interface{})
	t["Name"] = "Json-Demo"
	t["Ip"] = "127.0.0.1"
	t["Port"] = 8080

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Marshal err: ", err.Error())
		return
	}
	fmt.Println("Marshal Json: ", string(b))
}

func testReflect() {
	t := T{}
	t.Name = "Demo-Json"
	t.IP = "127.0.0.1"
	t.Port = 8080

	getType := reflect.TypeOf(t)
	getValue := reflect.ValueOf(t)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	s.Field(0).SetInt(77)
	s.Field(1).SetString("Sunset Strip")
	fmt.Println("t is now", t)
}
