package param

type ParamURLDirScan struct {
	baseParam
	URL string `form:"url" json:"url" binding:"required"`
}
