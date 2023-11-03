package model

type restRequest struct {
	RequestContext restRequestContext `json:"requestContext"`
	Req            interface{}        `json:"req"`
}

type restRequestContext struct {
	PartnerId string `json:"partnerId"`
	Serial    string `json:"serial"`
	Timestamp string `json:"timestamp"`
}

type restResponse struct {
	Result_code string      `json:"result_code"`
	Result_info string      `json:"result_info"`
	Rsp         interface{} `json:"rsp"`
}
