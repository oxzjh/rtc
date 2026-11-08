package rtc

import "time"

type Option func(*rtc)

func WithTimeout(timeout time.Duration) Option {
	return func(r *rtc) {
		r.timeout = timeout
	}
}

func WithOrdered(ordered bool) Option {
	return func(r *rtc) {
		r.ordered = ordered
	}
}

func WithOnConnected(onConnected func()) Option {
	return func(r *rtc) {
		r.onConnected = onConnected
	}
}

func WithOnTimeout(onTimeout func()) Option {
	return func(r *rtc) {
		r.onTimeout = onTimeout
	}
}

func WithOnClosed(onClosed func()) Option {
	return func(r *rtc) {
		r.onClosed = onClosed
	}
}
