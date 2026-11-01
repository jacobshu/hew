package forestfox

import (
	_ "embed"
	"encoding/json"
	stdcolor "image/color"
	"log"

	"charm.land/lipgloss/v2"
)

//go:embed forestfox.json
var ffData []byte

type color struct {
	Value    string
	Ghostty  []string
	Prose    string
	Lipgloss stdcolor.Color
}

type Forestfox struct {
	BgDim         color `json:"bg_dim"`
	Bg0           color `json:"bg0"`
	Bg1           color `json:"bg1"`
	Bg2           color `json:"bg2"`
	Bg3           color `json:"bg3"`
	Bg4           color `json:"bg4"`
	Bg5           color `json:"bg5"`
	Grey0         color `json:"grey0"`
	Grey1         color `json:"grey1"`
	Grey2         color `json:"grey2"`
	Grey3         color `json:"grey3"`
	Fg0           color `json:"fg0"`
	Fg1           color `json:"fg1"`
	Fg2           color `json:"fg2"`
	BgRed         color `json:"bg_red"`
	DimRed        color `json:"dim_red"`
	Red           color `json:"red"`
	BrightRed     color `json:"bright_red"`
	BgGreen       color `json:"bg_green"`
	DimGreen      color `json:"dim_green"`
	Green         color `json:"green"`
	BrightGreen   color `json:"bright_green"`
	BgYellow      color `json:"bg_yellow"`
	DimYellow     color `json:"dim_yellow"`
	Yellow        color `json:"yellow"`
	BrightYellow  color `json:"bright_yellow"`
	BgBlue        color `json:"bg_blue"`
	DimBlue       color `json:"dim_blue"`
	Blue          color `json:"blue"`
	BrightBlue    color `json:"bright_blue"`
	BgMagenta     color `json:"bg_magenta"`
	DimMagenta    color `json:"dim_magenta"`
	Magenta       color `json:"magenta"`
	BrightMagenta color `json:"bright_magenta"`
	BgAqua        color `json:"bg_aqua"`
	DimAqua       color `json:"dim_aqua"`
	Aqua          color `json:"aqua"`
	BrightAqua    color `json:"bright_aqua"`
	BgOrange      color `json:"bg_orange"`
	DimOrange     color `json:"dim_orange"`
	Orange        color `json:"orange"`
	BrightOrange  color `json:"bright_orange"`
}

func GetTheme() Forestfox {
	var cfg Forestfox
	if err := json.Unmarshal(ffData, &cfg); err != nil {
		log.Fatalf("error unmarshalling forestfox.json: %v", err)
	}

	var ff Forestfox

	ff.BgDim.Lipgloss = lipgloss.Color(cfg.BgDim.Value)
	ff.Bg0.Lipgloss = lipgloss.Color(cfg.Bg0.Value)
	ff.Bg1.Lipgloss = lipgloss.Color(cfg.Bg1.Value)
	ff.Bg2.Lipgloss = lipgloss.Color(cfg.Bg2.Value)
	ff.Bg3.Lipgloss = lipgloss.Color(cfg.Bg3.Value)
	ff.Bg4.Lipgloss = lipgloss.Color(cfg.Bg4.Value)
	ff.Bg5.Lipgloss = lipgloss.Color(cfg.Bg5.Value)
	ff.Grey0.Lipgloss = lipgloss.Color(cfg.Grey0.Value)
	ff.Grey1.Lipgloss = lipgloss.Color(cfg.Grey1.Value)
	ff.Grey2.Lipgloss = lipgloss.Color(cfg.Grey2.Value)
	ff.Grey3.Lipgloss = lipgloss.Color(cfg.Grey3.Value)
	ff.Fg0.Lipgloss = lipgloss.Color(cfg.Fg0.Value)
	ff.Fg1.Lipgloss = lipgloss.Color(cfg.Fg1.Value)
	ff.Fg2.Lipgloss = lipgloss.Color(cfg.Fg2.Value)
	ff.BgRed.Lipgloss = lipgloss.Color(cfg.BgRed.Value)
	ff.DimRed.Lipgloss = lipgloss.Color(cfg.DimRed.Value)
	ff.Red.Lipgloss = lipgloss.Color(cfg.Red.Value)
	ff.BrightRed.Lipgloss = lipgloss.Color(cfg.BrightRed.Value)
	ff.BgGreen.Lipgloss = lipgloss.Color(cfg.BgGreen.Value)
	ff.DimGreen.Lipgloss = lipgloss.Color(cfg.DimGreen.Value)
	ff.Green.Lipgloss = lipgloss.Color(cfg.Green.Value)
	ff.BrightGreen.Lipgloss = lipgloss.Color(cfg.BrightGreen.Value)
	ff.BgYellow.Lipgloss = lipgloss.Color(cfg.BgYellow.Value)
	ff.DimYellow.Lipgloss = lipgloss.Color(cfg.DimYellow.Value)
	ff.Yellow.Lipgloss = lipgloss.Color(cfg.Yellow.Value)
	ff.BrightYellow.Lipgloss = lipgloss.Color(cfg.BrightYellow.Value)
	ff.BgBlue.Lipgloss = lipgloss.Color(cfg.BgBlue.Value)
	ff.DimBlue.Lipgloss = lipgloss.Color(cfg.DimBlue.Value)
	ff.Blue.Lipgloss = lipgloss.Color(cfg.Blue.Value)
	ff.BrightBlue.Lipgloss = lipgloss.Color(cfg.BrightBlue.Value)
	ff.BgMagenta.Lipgloss = lipgloss.Color(cfg.BgMagenta.Value)
	ff.DimMagenta.Lipgloss = lipgloss.Color(cfg.DimMagenta.Value)
	ff.Magenta.Lipgloss = lipgloss.Color(cfg.Magenta.Value)
	ff.BrightMagenta.Lipgloss = lipgloss.Color(cfg.BrightMagenta.Value)
	ff.BgAqua.Lipgloss = lipgloss.Color(cfg.BgAqua.Value)
	ff.DimAqua.Lipgloss = lipgloss.Color(cfg.DimAqua.Value)
	ff.Aqua.Lipgloss = lipgloss.Color(cfg.Aqua.Value)
	ff.BrightAqua.Lipgloss = lipgloss.Color(cfg.BrightAqua.Value)
	ff.BgOrange.Lipgloss = lipgloss.Color(cfg.BgOrange.Value)
	ff.DimOrange.Lipgloss = lipgloss.Color(cfg.DimOrange.Value)
	ff.Orange.Lipgloss = lipgloss.Color(cfg.Orange.Value)
	ff.BrightOrange.Lipgloss = lipgloss.Color(cfg.BrightOrange.Value)

	return ff
}
