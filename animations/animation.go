package animations

type Animation interface {
	Update()
	Frame() int
}

type LoopAnimation struct {
	first        int
	last         int
	step         int     // how many indices do we move per frame
	speedInTps   float32 // how many ticks before next frame
	frameCounter float32
	frame        int
}

func (a *LoopAnimation) Update() {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame > a.last {
			a.frame = a.first
		}
	}
}
func (a *LoopAnimation) Frame() int {
	return a.frame
}

func NewLoopAnimation(first, last, step int, speed float32) Animation {
	return &LoopAnimation{
		first,
		last,
		step,
		speed,
		speed,
		first,
	}
}

type OneTimeAnimation struct {
	first        int
	last         int
	step         int
	speedInTps   float32
	frameCounter float32
	frame        int
	stopped      bool
}

func (a *OneTimeAnimation) Update() {
	if a.stopped {
		return
	}
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame >= a.last {
			a.stopped = true
		}
	}
}
func (a *OneTimeAnimation) Frame() int {
	return a.frame
}

func NewOneTimeAnimation(first, last, step int, speed float32) Animation {
	return &OneTimeAnimation{
		first,
		last,
		step,
		speed,
		speed,
		first,
		false,
	}
}
