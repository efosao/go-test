package models

type ThemeOption struct {
	Value    string
	Label    string
	Selected bool
}

type Config struct {
	Path         string
	ShowNav      bool
	Theme        string
	ThemeOptions []ThemeOption
}
