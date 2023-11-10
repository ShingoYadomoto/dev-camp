package game

var (
	groundY = 400

	//groundHeight = 50
	groundWidth = 50
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
