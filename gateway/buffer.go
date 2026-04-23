package main

type RingBuffer struct {
    data     []Metric
    size     int
    cursor   int 
    full     bool
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        data: make([]Metric, size),
        size: size,
    }
}

func (r *RingBuffer) Add(m Metric) {
    r.data[r.cursor] = m
    r.cursor = (r.cursor + 1) % r.size
    if r.cursor == 0 {
        r.full = true
    }
}

func (r *RingBuffer) GetAll() []Metric {
    if !r.full {
        return r.data[:r.cursor]
    }
    return r.data
}

func (r *RingBuffer) Reset() {
    r.cursor = 0
    r.full = false
    r.data = make([]Metric, r.size) 
}