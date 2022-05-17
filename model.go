package miio

type ids struct {
	Did  string `json:"did"`
	Siid int    `json:"siid"`
	Piid int    `json:"piid"`
}

type value struct {
	Value interface{} `json:"value"`
}

type resp struct {
	Code int `json:"code"`
}

type GetPropertyReq struct {
	ids
}

type GetPropertyResp struct {
	ids
	value
	resp
}

type SetPropertyReq struct {
	ids
	value
}

type SetPropertyResp struct {
	ids
	resp
}
