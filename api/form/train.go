package form

type TrainForm struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Conf     string `json:"conf" form:"conf" binding:"required"`
}
