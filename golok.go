package golok

import (
	"context"
	"fmt"
	"sync"
)

type profileHandle struct {
	ch     <-chan *[]*string
	ctx    context.Context
	cancel context.CancelFunc
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
	idx := index
	ctx, cancel := context.WithCancel(context.Background())

	if idx >= len(g.profileList) {
		idx = len(g.profileList)
		g.profileList = append(
			g.profileList,
			&profileHandle{
				ch:     ch,
				ctx:    ctx,
				cancel: cancel,
			},
		)
	} else {
		g.profileList = insert(
			g.profileList,
			idx,
			&profileHandle{
				ch:     ch,
				ctx:    ctx,
				cancel: cancel,
			},
		)
	}

	go g.profileWorker(idx)

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
	g.moveCursorAndClear(len(g.profileList))
}

func (g *Golok) moveCursorAndClear(index int) {
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

	fmt.Print("\r\033[2K")
}

func (g *Golok) profileWorker(index int) {
	for {
		select {
		case <-g.profileList[index].ctx.Done():
			remove(g.profileList, index)
			return

		case lok := <-g.profileList[index].ch:
			g.mtx.Lock()
			g.moveCursorAndClear(index)
			if lok != nil {
				for _, s := range *lok {
					if s != nil {
						fmt.Print(*s)
					}
				}
			}
			g.mtx.Unlock()
		}
	}
}
