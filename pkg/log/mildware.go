package middleware

import "strings"

var fieldConfig FieldConfig

func InitializeFieldConfig(c LogConfig) {
	if len(c.Duration) > 0 {
		fieldConfig.Duration = c.Duration
	} else {
		fieldConfig.Duration = "duration"
	}
	fieldConfig.Log = c.Log
	fieldConfig.Ip = c.Ip
	if c.Map != nil && len(c.Map) > 0 {
		fieldConfig.Map = c.Map
	}
	if c.Constants != nil && len(c.Constants) > 0 {
		fieldConfig.Constants = c.Constants
	}
	if len(c.Fields) > 0 {
		fields := strings.Split(c.Fields, ",")
		fieldConfig.Fields = fields
	}
	if len(c.Masks) > 0 {
		fields := strings.Split(c.Masks, ",")
		fieldConfig.Masks = fields
	}
	if len(c.Skips) > 0 {
		fields := strings.Split(c.Skips, ",")
		fieldConfig.Skips = fields
	}
}
