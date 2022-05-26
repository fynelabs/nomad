package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.NRGBA{0x18, 0x0C, 0x27, 0xFF}
	}

	if name == theme.ColorNamePrimary {
		return color.NRGBA{0xFF, 0x85, 0x00, 0x00}
	}

	if name == theme.ColorNameInputBackground {
		return color.NRGBA{0x00, 0x00, 0x00, 0x00}
	}

	if name == theme.ColorNamePlaceHolder {
		return color.NRGBA{0xFF, 0xFF, 0xFF, 0x40}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {

	if style.Bold && style.Italic {
		//Heavy with italics
		return resourceWorkSansBlackItalicTtf
	} else if style.Bold {
		//heavy
		return resourceWorkSansBlackTtf
	} else if style.Monospace {
		//Spaced out smaller font
		return resourceWorkSansRegularTtf
	}
	//standard bold
	return resourceWorkSansBoldTtf

}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)

}
