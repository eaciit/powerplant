package models

import (
	. "github.com/eaciit/orm"
)

type AssetClass struct {
	ModelBase `bson:"-",json:"-"`
	OID       string `bson:"_id"`
	Id        string `bson:"code"`
	Title     string `bson:"classname"`
}

func (u *AssetClass) TableName() string {
	return "SampleAssetClass"
}

func (c *AssetClass) RecordID() interface{} {
	return c.OID
}

type AssetLevel struct {
	ModelBase `bson:"base"`
	OID       string `bson:"_id"`
	Id        string `bson:"code"`
	Title     string `bson:"levelname"`
}

func (u *AssetLevel) TableName() string {
	return "SampleAssetType"
}

func (c *AssetLevel) RecordID() interface{} {
	return c.OID
}

type AssetType struct {
	ModelBase `bson:"base"`
	PID       string `bson:"_id"`
	Id        string `bson:"code"`
	Title     string `bson:"typename"`
}

func (u *AssetType) TableName() string {
	return "SampleAssetLevel"
}

func (c *AssetType) RecordID() interface{} {
	return c.PID
}
