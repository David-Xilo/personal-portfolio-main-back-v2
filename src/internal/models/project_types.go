package models

type ProjectType string

const (
	ProjectTypeUndefined ProjectType = "undefined"
	ProjectTypeTech      ProjectType = "tech"
	ProjectTypeGame      ProjectType = "game"
	ProjectTypeFinance   ProjectType = "finance"
)
