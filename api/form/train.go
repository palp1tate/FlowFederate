package form

type TrainForm struct {
	ID   string `json:"id" form:"id" binding:"required"`
	Conf string `json:"conf" form:"conf" binding:"required"`
}
