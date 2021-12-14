package hackflow

type sqlmap struct {
	baseTool
}

func newSqlmap() Tool {
	return &sqlmap{
		baseTool: baseTool{
			name: SQLMAP,
			desp: "自动化sql注入工具",
		},
	}
}

func GetSqlmap() *sqlmap {
	return container.Get(SQLMAP).(*sqlmap)
}

type SqlmapConfig struct {
	TargetURL   string
	Proxy       string
	BulkFile    string
	RandomAgent bool
	Batch       bool
}

//SetDebug 是否开启Debug
func (u *sqlmap) SetDebug(isDebug bool) {

}
func (s *sqlmap) Run(config *SqlmapConfig) (resultCh chan string, err error) {
	//todo
	return nil, nil
}

func (s *sqlmap) GetDbs(config *SqlmapConfig) (dbs []string, err error) {
	//todo
	return nil, nil
}

func (s *sqlmap) GetTables() (tables []string, err error) {
	//todo
	return nil, nil
}
