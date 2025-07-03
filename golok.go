package golok

import (
	"context"
	"fmt"
	"sync"
)

type profileHandle struct {
	index  int
	ch     <-chan *[]*string
	ctx    context.Context
	cancel context.CancelFunc
	cache  *[]*string
}

type Golok struct {
	mtx         sync.Mutex
	cursor      int
	profileList []*profileHandle
}

func NewGolok() *Golok {
	g := &Golok{
		mtx:         sync.Mutex{},
		cursor:      0,
		profileList: []*profileHandle{},
	}
	return g
}

func (g *Golok) NewProfile(index int) *Profile {
	if index < 0 {
		return nil
	}

	ch := make(chan *[]*string)
	ctx, cancel := context.WithCancel(context.Background())

	if index >= len(g.profileList) {
		ph := &profileHandle{
			index:  len(g.profileList),
			ch:     ch,
			ctx:    ctx,
			cancel: cancel,
			cache:  nil,
		}
		g.profileList = append(g.profileList, ph)
		go g.profileWorker(ph)
	} else {
		ph := &profileHandle{
			index:  index,
			ch:     ch,
			ctx:    ctx,
			cancel: cancel,
			cache:  nil,
		}
		g.profileList = insert(g.profileList, index, ph)

		g.updateIndices()
		g.clearLines(len(g.profileList) - 1)
		g.printCaches()

		go g.profileWorker(ph)
	}

	return &Profile{
		ch:         ch,
		cancel:     cancel,
		components: []Component{},
	}
}

func (g *Golok) Close() {
	for _, v := range g.profileList {
		v.cancel()
	}
	g.move(len(g.profileList))
}

func (g *Golok) updateIndices() {
	for i, v := range g.profileList {
		v.index = i
	}
}

func (g *Golok) move(index int) {
	diff := g.cursor - index
	dist := absInt(diff)
	g.cursor = index

	if diff > 0 {
		for range dist {
			fmt.Print("\033[A")
		}
	} else if diff < 0 {
		for range dist {
			fmt.Print("\033[B")
		}
	}

	fmt.Print("\r")
}

func (g *Golok) clear() {
	fmt.Print("\033[2K")
}

func (g *Golok) clearLines(n int) {
	for i := n; i >= 0; i-- {
		g.move(i)
		g.clear()
	}
}

func (g *Golok) print(s *[]*string) {
	if s != nil {
		for _, log := range *s {
			if log != nil {
				fmt.Print(*log)
			}
		}
	}
}

func (g *Golok) printCaches() {
	for i, v := range g.profileList {
		g.move(i)
		g.clear()
		g.print(v.cache)
	}
}

func (g *Golok) profileWorker(ph *profileHandle) {
	for {
		select {
		case <-ph.ctx.Done():
			g.mtx.Lock()

			g.profileList = remove(g.profileList, ph.index)
			g.updateIndices()
			g.clearLines(len(g.profileList) + 1)
			g.printCaches()

			g.mtx.Unlock()
			return

		case lok := <-ph.ch:
			if lok != nil {
				ph.cache = lok
			} else {
				continue
			}

			g.mtx.Lock()

			g.move(ph.index)
			g.clear()
			g.print(lok)

			g.mtx.Unlock()
		}
	}
}
