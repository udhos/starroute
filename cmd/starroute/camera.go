package main

type camera struct {
	x, y   int
	sc     *scene
	cyclic bool
}

const camPanStep = 5

func newCamera(sc *scene, cyclic, centralize bool) *camera {
	c := &camera{sc: sc, cyclic: cyclic}

	if centralize {
		c.centralize()
		c.clamp()
	}

	return c
}

func (c *camera) mid() (int, int) {
	camXmax := c.maxX()
	camYmax := c.maxY()
	return camXmax / 2, camYmax / 2
}

// centralize centers the camera on the current scene.
func (c *camera) centralize() {
	c.x, c.y = c.mid()
	c.clamp()
}

// clamp forces the camera to remain within the tilemap.
func (c *camera) clamp() {
	if c.cyclic {
		// cyclic camera wraps around tilemap edges
		if c.x < 0 {
			c.x += c.sc.tiles.tilePixelWidth()
		} else {
			c.x = c.x % c.sc.tiles.tilePixelWidth()
		}
		if c.y < 0 {
			c.y += c.sc.tiles.tilePixelHeight()
		} else {
			c.y = c.y % c.sc.tiles.tilePixelHeight()
		}
		return
	}
	// non-cyclic camera cannot cross tilemap edges
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

// maxX returns the maximum x coordinate the camera can reach.
// maxX restricts the non-cyclic camera within the tilemap.
func (c *camera) maxX() int {
	return c.sc.tiles.tilePixelWidth() - c.sc.g.screenWidth
}

// maxY returns the maximum y coordinate the camera can reach.
// maxY restricts the non-cyclic camera within the tilemap.
func (c *camera) maxY() int {
	return c.sc.tiles.tilePixelHeight() - c.sc.g.screenHeight
}

// maxX returns the maximum x coordinate the cyclic camera can reach before resetting to 0.
func (c *camera) cyclicMaxX() int {
	return c.sc.tiles.tilePixelWidth() - 1
}

// maxY returns the maximum y coordinate the cyclic camera can reach before resetting to 0.
func (c *camera) cyclicMaxY() int {
	return c.sc.tiles.tilePixelHeight() - 1
}
