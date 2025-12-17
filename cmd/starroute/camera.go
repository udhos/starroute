package main

type camera struct {
	x, y int
}

const camPanStep = 5

func (c *camera) stepUp() {
	c.y = max(c.y-camPanStep, 0)
}

func (c *camera) stepDown() {
	c.y += camPanStep
}

func (c *camera) stepLeft() {
	c.x += camPanStep
}

func (c *camera) stepRight() {
	c.x = max(c.x-camPanStep, 0)
}
