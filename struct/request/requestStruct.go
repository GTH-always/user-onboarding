package requestStruct

type GetUserDetailsRequest struct {
	HandleName string `json:"handleName"`
}

type UserLogin struct {
	Password   string `json:"password"`
	HandleName string `json:"handleName"`
}
