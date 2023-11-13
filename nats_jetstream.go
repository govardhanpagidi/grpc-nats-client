package main

type ConversionReq struct {
	TenantID       int32   `protobuf:"varint,1,opt,name=tenantId" json:"tenantId,omitempty"`
	BankID         int32   `protobuf:"varint,2,opt,name=bankId" json:"bankId,omitempty"`
	BaseCurrency   string  `protobuf:"bytes,3,opt,name=baseCurrency" json:"baseCurrency,omitempty"`
	TargetCurrency string  `protobuf:"bytes,4,opt,name=targetCurrency" json:"targetCurrency,omitempty"`
	Tier           string  `protobuf:"bytes,5,opt,name=tier" json:"tier,omitempty"`
	Amount         float64 `protobuf:"fixed64,6,opt,name=amount" json:"amount,omitempty"`
	RequestedOn    string  `protobuf:"bytes,7,opt,name=requestedOn" json:"requestedOn,omitempty"`
}
