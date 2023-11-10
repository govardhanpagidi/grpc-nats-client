package fxconversion

import "github.com/golang/protobuf/proto"

type ConversionResponse struct {
	Amount          float64 `protobuf:"fixed64,1,opt,name=amount" json:"amount,omitempty"`
	ConvertedAmount float64 `protobuf:"fixed64,2,opt,name=convertedAmount" json:"convertedAmount,omitempty"`
	BaseCurrency    string  `protobuf:"bytes,3,opt,name=baseCurrency" json:"baseCurrency,omitempty"`
	TargetCurrency  string  `protobuf:"bytes,4,opt,name=targetCurrency" json:"targetCurrency,omitempty"`
	InitiatedOn     int64   `protobuf:"bytes,5,opt,name=initiatedOn" json:"initiatedOn,omitempty"`
	Rate            float64 `protobuf:"fixed64,6,opt,name=rate" json:"rate,omitempty"`
}

func (c ConversionResponse) Reset() {
	c = ConversionResponse{}
}

func (c ConversionResponse) String() string {
	return proto.CompactTextString(c)
}

func (c ConversionResponse) ProtoMessage() {

}
