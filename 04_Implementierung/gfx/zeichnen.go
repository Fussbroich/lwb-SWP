package gfx

// zeichnen.go — Zeichenprimitiven: Linien, Kreise, Rechtecke, Sektoren, Dreiecke.

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whitePixel dient als Quellbild für DrawTriangles (Ebitengine-Pflicht).
var whitePixel *ebiten.Image

func init() {
	whitePixel = ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
}

func linie(x1, y1, x2, y2 uint16) {
	vector.StrokeLine(drawTarget,
		float32(x1), float32(y1), float32(x2), float32(y2),
		1, gibStiftfarbe(), false)
}

func kreis(x, y, r uint16) {
	vector.StrokeCircle(drawTarget,
		float32(x), float32(y), float32(r),
		1, gibStiftfarbe(), true)
}

func vollkreis(x, y, r uint16) {
	vector.DrawFilledCircle(drawTarget,
		float32(x), float32(y), float32(r),
		gibStiftfarbe(), true)
}

func rechteck(x1, y1, b, h uint16) {
	vector.StrokeRect(drawTarget,
		float32(x1), float32(y1), float32(b), float32(h),
		1, gibStiftfarbe(), false)
}

func vollrechteck(x1, y1, b, h uint16) {
	vector.DrawFilledRect(drawTarget,
		float32(x1), float32(y1), float32(b), float32(h),
		gibStiftfarbe(), false)
}

// ===================== Kreissektoren =====================
//
// Die alte gfx-Konvention: Winkelangaben in Grad, 0° = Osten,
// Winkel wachsen entgegen dem Uhrzeigersinn (mathematisch positiv).
// In Bildschirmkoordinaten (Y nach unten) wird Y gespiegelt.

// gradZuBildschirm rechnet einen gfx-Winkel (Grad, 0=Ost, CCW)
// in Bildschirmkoordinaten um und liefert (cos, sin) für diesen Winkel.
func gradAufBildschirm(grad float64) (float64, float64) {
	rad := grad * math.Pi / 180.0
	return math.Cos(rad), -math.Sin(rad) // Y-Achse gespiegelt
}

// sektorPfad baut einen vector.Path für einen Kreissektor.
// Segmente bestimmen die Glattheit des Bogens.
func sektorPfad(cx, cy, r float32, w1, w2 uint16) vector.Path {
	start := float64(w1)
	end := float64(w2)
	if end <= start {
		end += 360
	}
	segments := int((end-start)/5) + 1
	if segments < 4 {
		segments = 4
	}

	var p vector.Path
	p.MoveTo(cx, cy)
	for i := 0; i <= segments; i++ {
		t := start + (end-start)*float64(i)/float64(segments)
		cosT, sinT := gradAufBildschirm(t)
		p.LineTo(cx+r*float32(cosT), cy+r*float32(sinT))
	}
	p.Close()
	return p
}

func kreissektor(x, y, r, w1, w2 uint16) {
	p := sektorPfad(float32(x), float32(y), float32(r), w1, w2)
	strokeOpts := &vector.StrokeOptions{Width: 1}
	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, strokeOpts)
	farbeAufVertices(vs, gibStiftfarbe())
	drawTarget.DrawTriangles(vs, is, whitePixel, nil)
}

func vollkreissektor(x, y, r, w1, w2 uint16) {
	p := sektorPfad(float32(x), float32(y), float32(r), w1, w2)
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	farbeAufVertices(vs, gibStiftfarbe())
	drawTarget.DrawTriangles(vs, is, whitePixel, nil)
}

func volldreieck(x1, y1, x2, y2, x3, y3 uint16) {
	var p vector.Path
	p.MoveTo(float32(x1), float32(y1))
	p.LineTo(float32(x2), float32(y2))
	p.LineTo(float32(x3), float32(y3))
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	farbeAufVertices(vs, gibStiftfarbe())
	drawTarget.DrawTriangles(vs, is, whitePixel, nil)
}

// farbeAufVertices setzt die Farbe aller Vertices auf clr.
func farbeAufVertices(vs []ebiten.Vertex, clr color.NRGBA) {
	cr := float32(clr.R) / 255
	cg := float32(clr.G) / 255
	cb := float32(clr.B) / 255
	ca := float32(clr.A) / 255
	for i := range vs {
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = cr * ca
		vs[i].ColorG = cg * ca
		vs[i].ColorB = cb * ca
		vs[i].ColorA = ca
	}
}
