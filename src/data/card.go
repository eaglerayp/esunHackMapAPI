package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Activity struct {
	Id              bson.ObjectId `bson:"_id"`
	Card            string        `bson:"Card"`
	DiscountType    string        `bson:"DiscountType"`
	ActivityName    string        `bson:"ActivityName"`
	Content         string        `bson:"Content"`
	Remark          string        `json:"-" bson:"Remark"`
	StoreBranchName string        `json:"StoreBranchName,omitempty" bson:"StoreBranchName"`
	Address         string        `json:"Address,omitempty" bson:"Address"`
	Tel             string        `json:"Tel,omitempty" bson:"Tel"`
	GoogleAddress   string        `json:"GoogleAddress,omitempty" bson:"GoogleAddress"`
	Time            string        `bson:"Time"`
	Longutitude     float64       `json:"Longutitude,omitempty" bson:"Longutitude"`
	Latitude        float64       `json:"Latitude,omitempty" bson:"Latitude"`
	IsToday         bool          `bson:"-"`
}

const (
	DB         = "esun"
	Collection = "Activity"
)

var (
	mongoSession *mgo.Session
)

func InitMongo(url string) (err error) {
	mongoSession, err = mgo.Dial(url)
	return
}

func GetJob() (results []Activity, err error) {
	results = []Activity{}
	// Card == card AND Discount == discount
	selector := bson.M{"$or": []bson.M{bson.M{"GoogleAddress": bson.M{"$exists": false}}, bson.M{"GoogleAddress": nil}}}
	err = mongoSession.DB(DB).C(Collection).Find(selector).All(&results)
	return
}

func Set(doc Activity) error {
	selector := bson.M{"_id": doc.Id}
	_, err := mongoSession.DB(DB).C(Collection).Upsert(selector, doc)
	return err
}

func GetActivity(card, discount []string) (results []Activity, err error) {
	results = []Activity{}
	// Card == card AND Discount == discount
	selector := bson.M{"Card": bson.M{"$in": card}, "DiscountType": bson.M{"$in": discount}}
	err = mongoSession.DB(DB).C(Collection).Find(selector).All(&results)
	return
}
