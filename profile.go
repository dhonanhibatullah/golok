package golok

import "context"

type Profile struct {
	ch         chan<- *[]*string
	cancel     context.CancelFunc
	components []Component
}

func (p *Profile) AddComponent(index int, c Component) {
	if index < 0 {
		return
	} else if index >= len(p.components) {
		p.components = append(p.components, c)
	} else {
		p.components = insert(p.components, index, c)
	}
}

func (p *Profile) Render() {
	var res []*string
	for _, c := range p.components {
		res = append(res, c.render())
	}
	p.ch <- &res
}

func (p *Profile) Close() {
	p.cancel()
	close(p.ch)
}
