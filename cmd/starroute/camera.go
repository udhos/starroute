package main

type camera struct {
	x, y int
	sc   *scene
}

const camPanStep = 5

func newCamera(sc *scene) *camera {
	return &camera{sc: sc}
}

// clamp forces the camera to remain within the tilemap.
func (c *camera) clamp() {
	c.x = max(min(c.x, c.maxX()), 0)
	c.y = max(min(c.y, c.maxY()), 0)
}

func (c *camera) stepUp() {
	c.y -= camPanStep
	c.clamp()
}

func (c *camera) stepDown() {
	c.y += camPanStep
	c.clamp()
}

func (c *camera) stepLeft() {
	c.x += camPanStep
	c.clamp()
}

func (c *camera) stepRight() {
	c.x -= camPanStep
	c.clamp()
}

func (c *camera) maxX() int {
	return c.sc.tiles.tilePixelWidth() - c.sc.g.screenWidth
}

func (c *camera) maxY() int {
	return c.sc.tiles.tilePixelHeight() - c.sc.g.screenHeight
}
