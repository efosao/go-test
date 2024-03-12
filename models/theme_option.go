package models

type ThemeOption struct {
	Value    string
	Label    string
	Selected bool
}

type Config struct {
	Path         string
	Theme        string
	ThemeOptions []ThemeOption
}
