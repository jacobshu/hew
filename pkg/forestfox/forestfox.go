package forestfox

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
)

type color struct {
	Value    string
	Ghostty  []string
	Prose    string
	Lipgloss lipgloss.Color
}

type Forestfox struct {
	BgDim         color
	Bg0           color
	Bg1           color
	Bg2           color
	Bg3           color
	Bg4           color
	Bg5           color
	Grey0         color
	Grey1         color
	Grey2         color
	Grey3         color
	Fg0           color
	Fg1           color
	Fg2           color
	BgRed         color
	DimRed        color
	Red           color
	BrightRed     color
	BgGreen       color
	DimGreen      color
	Green         color
	BrightGreen   color
	BgYellow      color
	DimYellow     color
	Yellow        color
	BrightYellow  color
	BgBlue        color
	DimBlue       color
	Blue          color
	BrightBlue    color
	BgMagenta     color
	DimMagenta    color
	Magenta       color
	BrightMagenta color
	BgAqua        color
	DimAqua       color
	Aqua          color
	BrightAqua    color
	BgOrange      color
	DimOrange     color
	Orange        color
	BrightOrange  color
}

func GetTheme() Forestfox {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	os := runtime.GOOS
	var ffPath string
	if os == "windows" {
		ffPath = path.Join(homeDir, "\\source\\repos\\dotfiles\\config\\forestfox.toml")
	} else if os == "darwin" { // "darwin" is the value for macOS
		ffPath = path.Join(homeDir, "/dev/dotfiles/ff/forestfox.toml")
	}

	var ff Forestfox
	_, err = toml.DecodeFile(ffPath, &ff)
	if err != nil {
		log.Printf("error reading forestfox toml: %+v", err)
	}

	ff.BgDim.Lipgloss = lipgloss.Color(ff.BgDim.Value)
	ff.Bg0.Lipgloss = lipgloss.Color(ff.Bg0.Value)
	ff.Bg1.Lipgloss = lipgloss.Color(ff.Bg1.Value)
	ff.Bg2.Lipgloss = lipgloss.Color(ff.Bg2.Value)
	ff.Bg3.Lipgloss = lipgloss.Color(ff.Bg3.Value)
	ff.Bg4.Lipgloss = lipgloss.Color(ff.Bg4.Value)
	ff.Bg5.Lipgloss = lipgloss.Color(ff.Bg5.Value)
	ff.Grey0.Lipgloss = lipgloss.Color(ff.Grey0.Value)
	ff.Grey1.Lipgloss = lipgloss.Color(ff.Grey1.Value)
	ff.Grey2.Lipgloss = lipgloss.Color(ff.Grey2.Value)
	ff.Grey3.Lipgloss = lipgloss.Color(ff.Grey3.Value)
	ff.Fg0.Lipgloss = lipgloss.Color(ff.Fg0.Value)
	ff.Fg1.Lipgloss = lipgloss.Color(ff.Fg1.Value)
	ff.Fg2.Lipgloss = lipgloss.Color(ff.Fg2.Value)
	ff.BgRed.Lipgloss = lipgloss.Color(ff.BgRed.Value)
	ff.DimRed.Lipgloss = lipgloss.Color(ff.DimRed.Value)
	ff.Red.Lipgloss = lipgloss.Color(ff.Red.Value)
	ff.BrightRed.Lipgloss = lipgloss.Color(ff.BrightRed.Value)
	ff.BgGreen.Lipgloss = lipgloss.Color(ff.BgGreen.Value)
	ff.DimGreen.Lipgloss = lipgloss.Color(ff.DimGreen.Value)
	ff.Green.Lipgloss = lipgloss.Color(ff.Green.Value)
	ff.BrightGreen.Lipgloss = lipgloss.Color(ff.BrightGreen.Value)
	ff.BgYellow.Lipgloss = lipgloss.Color(ff.BgYellow.Value)
	ff.DimYellow.Lipgloss = lipgloss.Color(ff.DimYellow.Value)
	ff.Yellow.Lipgloss = lipgloss.Color(ff.Yellow.Value)
	ff.BrightYellow.Lipgloss = lipgloss.Color(ff.BrightYellow.Value)
	ff.BgBlue.Lipgloss = lipgloss.Color(ff.BgBlue.Value)
	ff.DimBlue.Lipgloss = lipgloss.Color(ff.DimBlue.Value)
	ff.Blue.Lipgloss = lipgloss.Color(ff.Blue.Value)
	ff.BrightBlue.Lipgloss = lipgloss.Color(ff.BrightBlue.Value)
	ff.BgMagenta.Lipgloss = lipgloss.Color(ff.BgMagenta.Value)
	ff.DimMagenta.Lipgloss = lipgloss.Color(ff.DimMagenta.Value)
	ff.Magenta.Lipgloss = lipgloss.Color(ff.Magenta.Value)
	ff.BrightMagenta.Lipgloss = lipgloss.Color(ff.BrightMagenta.Value)
	ff.BgAqua.Lipgloss = lipgloss.Color(ff.BgAqua.Value)
	ff.DimAqua.Lipgloss = lipgloss.Color(ff.DimAqua.Value)
	ff.Aqua.Lipgloss = lipgloss.Color(ff.Aqua.Value)
	ff.BrightAqua.Lipgloss = lipgloss.Color(ff.BrightAqua.Value)
	ff.BgOrange.Lipgloss = lipgloss.Color(ff.BgOrange.Value)
	ff.DimOrange.Lipgloss = lipgloss.Color(ff.DimOrange.Value)
	ff.Orange.Lipgloss = lipgloss.Color(ff.Orange.Value)
	ff.BrightOrange.Lipgloss = lipgloss.Color(ff.BrightOrange.Value)

	return ff
}
