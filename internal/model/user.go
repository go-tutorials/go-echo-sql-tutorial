package model

import (
	"github.com/core-go/core"
	"time"
)

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"required,max=40" match:"equal"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" avro:"username" validate:"required,username,max=100" match:"prefix"`
	Email       string     `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" avro:"email" validate:"email,max=100" match:"prefix"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" avro:"phone" validate:"required,phone,max=18"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth" avro:"dateOfBirth"`
}


type ReqHeader struct {
	RequestId   string     `json:"requestId"`
	CorrId      string     `json:"corrId"`
	MobileNo    string     `json:"mobileNo"`
	CreditCard  string     `json:"creditCard"`
	RequestTime *time.Time `json:"requestTime"`
}

type UpdateReq struct {
	Header ReqHeader `json:"header"`
	Body   *User     `json:"body"`
}

type ResHeader struct {
	RequestId    string             `json:"requestId"`
	ResponseId   string             `json:"responseId"`
	RequestTime  *time.Time         `json:"requestTime"`
	ResponseTime *time.Time         `json:"responseTime"`
	CorrId       string             `json:"corrId"`
	Errors       *[]core.ErrorMessage `json:"errors,omitempty"`
}

type UpdateRes struct {
	Header ResHeader `json:"header"`
	Body   *ResBody  `json:"body"`
}
type ResBody struct {
	Res int64 `json:"count"`
}

func BuildResponseHeader(header ReqHeader, errors *[]core.ErrorMessage) ResHeader {
	now := time.Now()
	res := ResHeader{
		RequestId:    header.RequestId,
		CorrId:       header.CorrId,
		ResponseId:   "ResponseId",
		RequestTime:  header.RequestTime,
		ResponseTime: &now,
	}
	if errors != nil && len(*errors) > 0 {
		res.Errors = errors
	}
	return res
}
