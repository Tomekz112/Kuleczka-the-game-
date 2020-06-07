package main

import (
	"github.com/faiface/pixel"
)

type check struct {
	mousepos  pixel.Vec
	objectpos pixel.Vec
	ballpos   pixel.Vec
}

func (m *check) IsSame(mousepos, objectpos pixel.Vec) bool {
	return mousepos.Y < objectpos.Y+37 &&
		mousepos.Y > objectpos.Y-25 &&
		mousepos.X > objectpos.X-128 &&
		mousepos.X < objectpos.X+130
}

func (m *check) IsSame2(ballpos, objectpos pixel.Vec) bool {
	return ballpos.X > objectpos.X-28 &&
		ballpos.X < objectpos.X+28 &&
		ballpos.Y > objectpos.Y-28 &&
		ballpos.Y < objectpos.Y+28
}

type colision struct {
	Xminus    bool //X is below 0
	BallPos   pixel.Vec
	ObjectPos pixel.Vec
	Yminus    bool //ball goes from down or up
	Line      pixel.Line
}

func (c *colision) Average(Xminus bool, BallPos pixel.Vec, ObjectPos pixel.Vec, Yminus bool) float64 {
	line := pixel.L(ObjectPos, BallPos)
	a, b := pixel.ZV, pixel.ZV
	need := false
	point := pixel.ZV
	var GameMap pixel.Line
	a.Y, b.Y = ObjectPos.Y, ObjectPos.Y
	for need == false {
		if a.Y > BallPos.Y {
			a.Y--
			b.Y--
		} else {
			a.Y++
			b.Y++
		}
		a.X, b.X = 155, -155
		GameMap = pixel.L(a, b)
		point, need = line.Intersect(GameMap)
	}
	avgline := pixel.L(ObjectPos, point)
	avg := 0.0
	if Xminus == true {
		avg = avgline.Len() * -1
		avg++
	} else {
		avg = avgline.Len()
		avg--
	}
	return avg
}

func (c *colision) GoesXMinus(ObjectPos pixel.Vec, Line pixel.Line) bool {
	rectangle := Line.Bounds()
	vectors := rectangle.Vertices()
	end := false
	for i := 0; i == 3; i++ {
		if ObjectPos.X > vectors[i].X {
			end = true
		}
		i++
	}
	if end != false {
		return true
	} else {
		return false
	}
}
