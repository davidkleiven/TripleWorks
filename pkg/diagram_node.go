package pkg

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func FirstLetterOrEmpty(s string) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[0])
}

type DiagramNode struct {
	id int64

	Name string
	Type string
}

func (d DiagramNode) ID() int64 {
	return d.id
}

// DrawGlyph implements the GlyphDrawer interface
func (d *DiagramNode) Draw(c *draw.Canvas, pt vg.Point, r vg.Length) {
	icon := DefaultIcon{Letter: FirstLetterOrEmpty(d.Type)}
	icon.Draw(c, pt, r)

	labelPt := vg.Point{
		X: pt.X + vg.Length(1),
		Y: pt.Y + vg.Length(1),
	}
	ts := draw.TextStyle{
		Font:    font.From(plot.DefaultFont, 8),
		Color:   color.Black,
		Handler: plot.DefaultTextHandler,
	}
	width, height, _ := ts.Handler.Box(d.Name, ts.Font)

	rect := vg.Rectangle{
		Min: vg.Point{
			X: labelPt.X,
			Y: labelPt.Y,
		},
		Max: vg.Point{
			X: labelPt.X + width,
			Y: labelPt.Y + height,
		},
	}
	c.SetColor(color.White)
	c.Fill(rect.Path())

	c.SetColor(color.Black)
	c.FillText(ts, labelPt, d.Name)
}

type DefaultIcon struct {
	Letter string
}

func (d *DefaultIcon) Draw(c *draw.Canvas, pt vg.Point, size vg.Length) {
	half := size / 2

	rect := vg.Rectangle{
		Min: vg.Point{
			X: pt.X - half,
			Y: pt.Y - half,
		},
		Max: vg.Point{
			X: pt.X + half,
			Y: pt.Y + half,
		},
	}

	c.SetColor(color.White)
	c.Fill(rect.Path())

	c.SetColor(color.Black)
	c.Stroke(rect.Path())

	ts := draw.TextStyle{
		Font:    font.From(plot.DefaultFont, 12),
		Color:   color.Black,
		XAlign:  draw.XCenter,
		YAlign:  draw.YCenter,
		Handler: plot.DefaultTextHandler,
	}
	c.FillText(ts, pt, d.Letter)

}
