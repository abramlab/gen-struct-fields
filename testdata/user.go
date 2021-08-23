package testdata

type User struct {
	Name         string `custom_tag:"name"`
	Login        string `custom_tag:"username"`
	Password     string `custom_tag:"-"`
	CustomField  string `custom_tag:""`
	CustomField1 string `json:"field"`
	AuthType     int    `custom_tag:"auth_type"`
}
