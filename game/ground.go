package game

var (
	groundY = ScreenY - ScreenY/8

	groundWidth = 100
)

type ground struct {
	x int
	y int
}

func (g *ground) move(speed int) {
	g.x -= speed
	if g.x < -groundWidth {
		g.x = g.x + groundWidth
	}
}
