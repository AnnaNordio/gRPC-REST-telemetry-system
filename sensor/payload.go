package main

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
)

func getJsonSize(v interface{}) float64 {
	b, _ := json.Marshal(v)
	return float64(len(b)) 
}

func getProtoSize(m proto.Message) float64 {
	b, _ := proto.Marshal(m)
	return float64(len(b))
}