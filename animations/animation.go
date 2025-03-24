package animations

type Animation interface {
	// returns true if finished
	Update() bool
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

func (a *LoopAnimation) Update() bool {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame > a.last {
			a.frame = a.first
		}
	}
	return false
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

type SingleFrameAnimation struct {
	frame int
}

// Frame implements Animation.
func (s *SingleFrameAnimation) Frame() int {
	return s.frame
}

// Update implements Animation.
func (s *SingleFrameAnimation) Update() bool {
	return false
}

func NewSingleFrameAnimation(frame int) Animation {
	return &SingleFrameAnimation{
		frame: frame,
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
	removeAfter         bool
}

func (a *OneTimeAnimation) Update() (remove bool) {
	if a.stopped {
		return a.removeAfter
	}
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame >= a.last {
			a.stopped = true
		}
	}
	return false
}
func (a *OneTimeAnimation) Frame() int {
	return a.frame
}

func NewOneTimeAnimation(first, last, step int, speed float32, removeAfter bool) Animation {
	return &OneTimeAnimation{
		first,
		last,
		step,
		speed,
		speed,
		first,
		false,
		removeAfter,
	}
}
