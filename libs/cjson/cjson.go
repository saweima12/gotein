package cjson

import jsoniter "github.com/json-iterator/go"

var (
	c                   = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal             = c.Marshal
	MarshalToString     = c.MarshalToString
	Unmarshal           = c.Unmarshal
	UnmarshalFromString = c.UnmarshalFromString
)
