package structs

type UserDetails struct {
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	HandleName   string `json:"handlename,omitempty"`
	Email        string `json:"email,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Number       int    `json:"phoneNumber,omitempty"`
	SocialHandle string `json:"socialHandle,omitempty"`
	BankDetails  string `json:"bankDetails,omitempty"`
	Image        string `json:"image,omitempty"`
	Password     string `json:"password,omitempty"`
	Resume       string `json:"resume,omitempty"`
	Pincode      int    `json:"pinCode,omitempty"`
	Type         int    `json:"type,omitempty"`
}
