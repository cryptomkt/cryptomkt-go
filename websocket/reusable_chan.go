package websocket

type reusableChan struct {
	channel   chan<- []byte
	callCount int
	closed    bool
}

func newReusableChan(channel chan<- []byte, callCount int) *reusableChan {
	return &reusableChan{
		channel:   channel,
		callCount: callCount,
		closed:    false,
	}
}

func (reusableChan *reusableChan) isDone() bool {
	return reusableChan.callCount < 1
}

func (reusableChan *reusableChan) send(data []byte) {
	if reusableChan.isDone() {
		return
	}
	reusableChan.callCount--
	reusableChan.channel <- data
}

func (reusableChan *reusableChan) close() {
	if reusableChan.closed {
		return
	}
	reusableChan.closed = true
	close(reusableChan.channel)
}
