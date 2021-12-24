package param

type ParamStartURLDirScan struct {
	baseParam
	URL string `form:"url" json:"url" binding:"required"`
}

type ParamStopURLDirScan struct {
	baseParam
	URL string `form:"url" json:"url" binding:"required"`
}

type ParamDeleteURLSubDir struct {
	baseParam
	ParentURL string `form:"parent_url" json:"parent_url" binding:"required"`
	SubURL    string `form:"sub_url" json:"sub_url" binding:"required"`
}

type ParamGetSubDir struct {
	Page
	ParentURL string `form:"parent_url" json:"parent_url" binding:"required"`
}
