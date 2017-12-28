package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Activity struct {
	Id            bson.ObjectId `bson:"_id"`
	Card          string        `bson:"Card"`
	DiscountType  string        `bson:"DiscountType"`
	ActivityName  string        `bson:"ActivityName"`
	Content       string        `bson:"Content"`
	Remark        string        `json:"-" bson:"Remark"`
	StoreName     string        `json:"StoreName,omitempty" bson:"StoreName"`
	Branch        string        `json:"Branch,omitempty" bson:"Branch"`
	Address       string        `json:"Address,omitempty" bson:"Address"`
	Tel           string        `json:"Tel,omitempty" bson:"Tel"`
	GoogleAdrress string        `json:"GoogleAdrress,omitempty" bson:"GoogleAdrress"`
	Time          string        `bson:"Time"`
	Longutitude   float64       `json:"Longutitude,omitempty" bson:"Longutitude"`
	Latitude      float64       `json:"Latitude,omitempty" bson:"Latitude"`
	IsToday       bool          `bson:"-"`
}

var (
	mongoSession *mgo.Session
)

func InitMongo(url string) (err error) {
	mongoSession, err = mgo.Dial(url)
	return
}

func GetActivity(card, discount []string) (results []Activity, err error) {
	results = []Activity{}
	// Card == card AND Discount == discount
	selector := bson.M{"Card": bson.M{"$in": card}, "DiscountType": bson.M{"$in": discount}}
	err = mongoSession.DB("esun").C("Activity").Find(selector).All(&results)
	return
}
