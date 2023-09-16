package models

type DirectoryStat struct {
	IsDirectory bool
	Name        string // name of directory to show user
	Uri         string // the uri (link) to the directory
	Image       string // preview image of directory
}
