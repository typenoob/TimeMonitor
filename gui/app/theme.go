package app

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type themeWithNoPadding struct{}

var _ fyne.Theme = (*themeWithNoPadding)(nil)

func (m themeWithNoPadding) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m themeWithNoPadding) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m themeWithNoPadding) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m themeWithNoPadding) Size(name fyne.ThemeSizeName) float32 {
	log.Println(name)
	if name == theme.SizeNamePadding {
		log.Println("*")
		return 0
	} else {
		return theme.DefaultTheme().Size(name)
	}
}
