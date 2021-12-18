package hackflow

type Installer interface {
	Install(link, dst string) (dirpath string, err error)
}
