package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"math/rand"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel/imdraw"

	_ "image/png"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

var mopos (*check)
var AvgCol (*colision)
var mode string = "menu"
var stop bool = false

var szerokosc2 float64 = 400 * -1

var (
	PositionOfPlayer1 = pixel.ZV
	PositionOfPlayer2 = pixel.ZV
	GameFreeze        = true
	Speed             = 0.0
	BallPos           = pixel.ZV
	Xm                = true
	reflection        = "player1"
	Yminus            = true
	average           = 0.0
)

func reset() {
	reflection = "player1"
	Yminus = true
	average = 0.0
	Speed /= 2
	BallPos.X, BallPos.Y = 0, 0
	PositionOfPlayer1.X, PositionOfPlayer1.Y, PositionOfPlayer2.X, PositionOfPlayer2.Y = 0, 180, 0, -180
	GameFreeze = true
}

func run() {
	icon, err := loadPicture("img/icon.png")
	if err != nil {
		panic(err)
	}
	cfg := pixelgl.WindowConfig{
		Title:     "Kuleczka the game!",
		Bounds:    pixel.R(0, 0, 400, 400),
		VSync:     true,
		Icon:      []pixel.Picture{icon},
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	f, err := os.Open("sounds/beep.mp3")
	if err != nil {
		log.Fatal(err)
	}
	k, err := os.Open("sounds/music.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer1, format, err := mp3.Decode(k)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	defer streamer1.Close()
	streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/39))
	Lines := []pixel.Line{}
	Slice := []pixel.Vec{}

	player1, err := loadPicture("img/platforma1.png")
	if err != nil {
		panic(err)
	}
	PositionOfPlayer1.Y = 180
	PositionOfPlayer1.X = 0

	player2, err := loadPicture("img/platforma2.png")
	if err != nil {
		panic(err)
	}
	spritesheet, err := loadPicture("img/wspomagacz.png")
	if err != nil {
		panic(err)
	}
	var Frames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 25 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 25 {
			Frames = append(Frames, pixel.R(x, y, x+25, y+25))
		}
	}

	modes, err := loadPicture("img/mode.png")
	if err != nil {
		panic(err)
	}

	control, err := loadPicture("img/controls.png")
	if err != nil {
		panic(err)
	}
	var modeFrames []pixel.Rect
	for x := modes.Bounds().Min.X; x < modes.Bounds().Max.X; x += 260 {
		for y := modes.Bounds().Min.Y; y < modes.Bounds().Max.Y; y += 100 {
			modeFrames = append(modeFrames, pixel.R(x, y, x+260, y+100))
		}
	}
	modesel, modesel2, modesel3, modesel4, modesel5, modesel6 := pixel.NewSprite(modes, modeFrames[0]), pixel.NewSprite(modes, modeFrames[1]), pixel.NewSprite(modes, modeFrames[2]), pixel.NewSprite(modes, modeFrames[3]), pixel.NewSprite(modes, modeFrames[4]), pixel.NewSprite(modes, modeFrames[5])
	modepos := []pixel.Vec{pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV}
	modepos[0].Y, modepos[1].Y, modepos[2].Y, modepos[3].Y, modepos[4].Y, modepos[5].Y = 30, -30, 175, 175, -175, -175
	modepos[2].X, modepos[3].X, modepos[4].X, modepos[5].X = -100, 142, -85, 142

	PosEditMode, Pos2pMode, Pos1pMode := pixel.ZV, pixel.ZV, pixel.ZV
	PosEditMode.Y, PosEditMode.X, Pos2pMode.Y, Pos2pMode.X, Pos1pMode.Y, Pos1pMode.X = -150, 0, -75, 0, 0, 0
	heart, err := loadPicture("img/serce.png")
	if err != nil {
		panic(err)
	}
	player1hp, player2hp := 3, 3
	Health := []float64{-140.0, -150.0, -160.0, 160.0, 150.0, 140.0, 180.0, -180.0}
	HP := []pixel.Vec{pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV, pixel.ZV}
	for h := 0; h < 6; h++ {
		HP[h].X = Health[h]
		if h < 3 {
			HP[h].Y = Health[6]
		} else {
			HP[h].Y = Health[7]
		}
	}
	ball, err := loadPicture("img/kula.png")
	if err != nil {
		panic(err)
	}
	reset()
	var MovingSpeed float64 = 100.0
	last := time.Now()
	boostPos, boostwait, boost, Player1BoostSlot, Player2BoostSlot, BoostNumber := pixel.ZV, 0, pixel.NewSprite(spritesheet, Frames[rand.Intn(len(Frames))]), -1, -1, 1
	ctrls := &beep.Ctrl{Streamer: beep.Loop(-1, streamer1), Paused: false}
	Player1spdBoost, Player2spdBoost := 0.0, 0.0
	reversemovp1, ToStopRevMovP1, reversemovp2, ToStopRevMovP2 := false, 0, false, 0
	BegginingOfLineSet, EndOfLineSet := false, false
	help := false
	imd := imdraw.New(nil)
	BegginingOfLine, EndOfLine := pixel.ZV, pixel.ZV
	rand.Seed(time.Now().UnixNano())
	volume1 := &effects.Volume{
		Streamer: ctrls,
		Base:     2,
		Volume:   -1.5,
		Silent:   false,
	}
	speaker.Play(volume1)
	var size [5]float64
	size[0], size[1], size[2], size[3], size[4] = 0.8, 0.8, 0.8, 0.8, 0.8
	LastReflection := "player2"
	fontsize := [7]float64{0.8, 0.8, 0.8, 0.7, 0.8, 0.7, 0.8}
	repeat := 0
	for !win.Closed() {
		prostokatg1 := pixel.R(PositionOfPlayer1.X-50, PositionOfPlayer1.Y-10, PositionOfPlayer1.X+50, PositionOfPlayer1.Y+15)
		prostokatg2 := pixel.R(PositionOfPlayer2.X-50, PositionOfPlayer2.Y-10, PositionOfPlayer2.X+50, PositionOfPlayer2.Y+5)
		circle := pixel.C(BallPos, 3)
		dt := time.Since(last).Seconds()
		last = time.Now()
		if GameFreeze == false {
			if BallPos.Y > 190 {
				fmt.Println("Player 1 just lost 1 Health")
				reset()
				player1hp--
				reflection = "player2"
				Yminus = true
			} else if BallPos.Y < -190 {
				fmt.Println("Player 2 just lost 1 Health")
				reset()
				reflection = "player1"
				player2hp--
			}
			if prostokatg1.IntersectCircle(circle) != pixel.ZV {
				if BallPos.X > PositionOfPlayer1.X {
					Xm = false
				} else {
					Xm = true
				}
				average = (*colision).Average(AvgCol, Xm, BallPos, PositionOfPlayer1, true)
				reflection = "player1"
				Yminus = true
			}
			if prostokatg2.IntersectCircle(circle) != pixel.ZV {
				if BallPos.X > PositionOfPlayer2.X {
					Xm = false
				} else if BallPos.X < PositionOfPlayer2.X {
					Xm = true
				} else {
					average = 0
				}
				average = (*colision).Average(AvgCol, Xm, BallPos, PositionOfPlayer2, false)
				reflection = "player2"
				Yminus = false
			}
			a := 0
			for range Lines {
				if circle.IntersectLine(Lines[a]) != pixel.ZV {
					// if BallPos.X < Slice[a].X {
					// 	Xm = false
					// } else {
					// 	Xm = true
					// }
					Xm = (*colision).GoesXMinus(AvgCol, BallPos, Lines[a])
					if repeat == 0 {
						if reflection == "obstacle" && BallPos.Y < Slice[a].Y {
							Yminus = true
						} else if reflection == "obstacle" {
							Yminus = false
						}
						repeat++
					}
					average = (*colision).Average(AvgCol, Xm, BallPos, Slice[a], Yminus)
					if reflection != "obstacle" {
						reflection = "obstacle"
					} else {
						reflection = "obstacle2"
					}
				}
				a++
			}
			if average == 0 {
				BallPos.X += average
			} else if average < 0 {
				BallPos.X += average - Speed
			} else {
				BallPos.X += average + Speed
			}

			if BallPos.X < -180 {
				reflection = "sciana1"
			}
			if BallPos.X > 180 {
				reflection = "sciana2"
			}
			switch reflection {
			case "player1":
				BallPos.Y -= 1 + Speed
			case "player2":
				BallPos.Y += 1 + Speed
			case "sciana1":
				if Yminus == false {
					BallPos.Y += 1 + Speed
				} else {
					BallPos.Y -= 1 + Speed
				}
				average = 2
			case "sciana2":
				if Yminus == false {
					BallPos.Y += 1 + Speed
				} else {
					BallPos.Y -= 1 + Speed
				}
				average = -2
			default:
				if Yminus == false {
					BallPos.Y -= 1 + Speed
				} else {
					BallPos.Y += 1 + Speed
				}
			}
			if reflection != LastReflection {
				repeat = 0
				ToStopRevMovP1--
				ToStopRevMovP2--
				shot := buffer.Streamer(1, buffer.Len())
				volume := &effects.Volume{
					Streamer: shot,
					Base:     2,
					Volume:   -3,
					Silent:   false,
				}
				speaker.Play(volume)
				boostwait++
				if boostwait == 4 {
					stop = false
					boostwait = 0
				}
				LastReflection = reflection
				Speed += 0.15
			}
			if ToStopRevMovP1 == 0 {
				reversemovp1 = false
			} else if ToStopRevMovP2 == 0 {
				reversemovp2 = false
			}
			if PositionOfPlayer1.X < 155.0 {
				if mode == "multiplayer" {
					if win.Pressed(pixelgl.KeyRight) && reversemovp1 == false || win.Pressed(pixelgl.KeyLeft) && reversemovp1 == true {
						PositionOfPlayer1.X += MovingSpeed*dt + Player1spdBoost
					}
				} else if mode == "singleplayer" {
					if BallPos.X > PositionOfPlayer1.X {
						PositionOfPlayer1.X += MovingSpeed*dt + 1
					}
				}
			}
			if PositionOfPlayer1.X > -155.0 {
				if mode == "multiplayer" {
					if win.Pressed(pixelgl.KeyLeft) && reversemovp1 == false || win.Pressed(pixelgl.KeyLeft) && reversemovp1 == true {
						PositionOfPlayer1.X -= MovingSpeed*dt + Player1spdBoost
					}
				} else if mode == "singleplayer" {
					if PositionOfPlayer1.X > BallPos.X {
						PositionOfPlayer1.X -= MovingSpeed*dt + 1 + Player1spdBoost
					}
				}
			}
			if PositionOfPlayer2.X < 155.0 {
				if win.Pressed(pixelgl.KeyD) && reversemovp2 == false || win.Pressed(pixelgl.KeyA) && reversemovp2 == true {
					PositionOfPlayer2.X += MovingSpeed*dt + Player2spdBoost
				}
			}
			if PositionOfPlayer2.X > -155.0 {
				if win.Pressed(pixelgl.KeyA) && reversemovp2 == false || win.Pressed(pixelgl.KeyD) && reversemovp2 == true {
					PositionOfPlayer2.X -= MovingSpeed*dt + Player2spdBoost
				}
			}
		} // freeze end
		if win.JustPressed(pixelgl.KeySpace) {
			if mode == "singleplayer" || mode == "multiplayer" {
				GameFreeze = false
			}
		}
		mousepos := pixel.ZV
		mousepos2 := pixel.IM.Moved(win.Bounds().Center().Sub(mousepos))
		mouse := mousepos2.Unproject(win.MousePosition())
		if win.JustPressed(pixelgl.KeyB) {
			mode = "menu"
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			os.Exit(5)
		}
		if win.JustPressed(pixelgl.KeyM) || mode == "singleplayer" {
			switch Player1BoostSlot {
			case 0:
				player1hp++
			case 1:
				Player1spdBoost++
			case 2:
				Speed += 5.0
			case 3:
				reversemovp2 = true
				ToStopRevMovP2 = 3
			default:
			}
			Player1BoostSlot = -1
		}
		if win.JustPressed(pixelgl.KeyE) {
			switch Player2BoostSlot {
			case 0:
				player2hp++
			case 1:
				Player2spdBoost++
			case 2:
				Speed += 5.0
			case 3:
				reversemovp1 = true
				ToStopRevMovP1 = 3
			default:
			}
			Player2BoostSlot--
		}
		i := 0
		for range modepos {
			if (*check).IsSame(mopos, mouse, modepos[i]) == true && fontsize[i] < 1 && mode == "menu" {
				shot := buffer.Streamer(1, buffer.Len())
				volume := &effects.Volume{
					Streamer: shot,
					Base:     2,
					Volume:   -4,
					Silent:   false,
				}
				speaker.Play(volume)
				fontsize[i] *= 1.25
			} else if (*check).IsSame(mopos, mouse, modepos[i]) == false {
				fontsize[i] = 0.8
			}
			i++
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if mode == "menu" {
				if (*check).IsSame(mopos, mouse, modepos[0]) == true {
					mode = "singleplayer"
					GameFreeze = true
				} else if (*check).IsSame(mopos, mouse, modepos[1]) == true {
					mode = "multiplayer"
					GameFreeze = true
				} else if (*check).IsSame(mopos, mouse, modepos[2]) == true {
					mode = "edytor"
					GameFreeze = true
					help = true
				} else if (*check).IsSame(mopos, mouse, modepos[3]) == true {
					os.Exit(3)
				} else if (*check).IsSame(mopos, mouse, modepos[5]) == true {
					mode = "controls"
				}
			}
		}
		if mode == "edytor" && win.JustPressed(pixelgl.MouseButtonLeft) && help == false {
			BegginingOfLine = mouse
			BegginingOfLineSet = true
		} else if mode == "edytor" && win.JustPressed(pixelgl.MouseButtonRight) && help == false {
			EndOfLine = mouse
			EndOfLineSet = true
		}
		help = false
		if (*check).IsSame2(mopos, BallPos, boostPos) == true {
			boostPos.X = 6969
			boostPos.Y = 6969
			if GameFreeze != true || reflection == "player1" {
				Player1BoostSlot = BoostNumber
			} else {
				Player2BoostSlot = BoostNumber
			}
		}
		switch mode {
		case "menu":
			win.Clear(colornames.Mediumaquamarine)
		case "edytor":
			win.Clear(colornames.Lightpink)
		default:
			win.Clear(colornames.Lightcoral)
		}
		if mode != "menu" && mode != "controls" {
			Healthg1g2 := pixel.NewSprite(heart, heart.Bounds())
			if player1hp > 2 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[0])))
			}
			if player1hp > 1 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[1])))
			}
			if player1hp > 0 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[2])))
			} else {
				fmt.Println("Player 1 lost")
				os.Exit(3)
			}
			if player2hp > 2 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[3])))
			}
			if player2hp > 1 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[4])))
			}
			if player2hp > 0 {
				Healthg1g2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(HP[5])))
			} else {
				fmt.Println("Player 2 lost")
				os.Exit(3)
			}
			if stop == false {
				BoostNumber = rand.Intn(len(Frames))
				boost = pixel.NewSprite(spritesheet, Frames[BoostNumber])
				boostPos.X = -150.0 + rand.Float64()*(150.0-(-150.0))
				boostPos.Y = -150.0 + rand.Float64()*(150.0-(-150.0))
				stop = true
			}
			boost.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(boostPos)))
			player1 := pixel.NewSprite(player1, player1.Bounds())
			player1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(PositionOfPlayer1)))
			player2 := pixel.NewSprite(player2, player2.Bounds())
			player2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(win.Bounds().Center().Add(PositionOfPlayer2)))
			ball := pixel.NewSprite(ball, ball.Bounds())
			ball.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.3).Moved(win.Bounds().Center().Add(BallPos)))
		} else if mode == "controls" {
			controls := pixel.NewSprite(control, control.Bounds())
			controls.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(win.Bounds().Center()))
		} else {
			modesel.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[0]).Moved(win.Bounds().Center().Add(modepos[0])))
			modesel2.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[1]).Moved(win.Bounds().Center().Add(modepos[1])))
			modesel3.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[2]).Moved(win.Bounds().Center().Add(modepos[2])))
			modesel4.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[3]).Moved(win.Bounds().Center().Add(modepos[3])))
			modesel5.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[4]).Moved(win.Bounds().Center().Add(modepos[4])))
			modesel6.Draw(win, pixel.IM.Scaled(pixel.ZV, fontsize[5]).Moved(win.Bounds().Center().Add(modepos[5])))
		}

		if mode == "edytor" && BegginingOfLineSet == true && EndOfLineSet == true {
			imd.Color = colornames.Red
			imd.Push(pixel.V(BegginingOfLine.X+200, BegginingOfLine.Y+200))
			imd.Push(pixel.V(EndOfLine.X+200, EndOfLine.Y+200))
			imd.Line(5)
			BegginingOfLineSet = false
			EndOfLineSet = false
			Line := pixel.L(BegginingOfLine, EndOfLine)
			vector := Line.Center()
			Slice = append(Slice, vector)
			fmt.Println(vector)
			Lines = append(Lines, Line)

		}
		if mode != "menu" {
			imd.Draw(win)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
