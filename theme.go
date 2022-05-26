package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.RGBA{0x18, 0x0C, 0x27, 0xFF}
	}

	if name == theme.ColorNamePrimary {
		return color.RGBA{0xFF, 0x85, 0x00, 0x40}
	}

	if name == theme.ColorNameInputBackground {
		return color.RGBA{0x00, 0x00, 0x00, 0x00}
	}

	if name == theme.ColorNamePlaceHolder {
		// return color.RGBA{0xFF, 0xFF, 0xFF, 0x40} //Text doesn't display with transparency?
		return color.RGBA{0x5B, 0x5B, 0x5B, 0xFF}
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
		//regular
	} else {
		//standard bold
		return resourceWorkSansBoldTtf

	}
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)

}
