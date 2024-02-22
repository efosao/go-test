package models

type Cookie struct {
	Theme string `cookie:"theme"`
}

type key string

const ThemeOptionsKey = key("themeOptions")
const ThemeKey = key("theme")
