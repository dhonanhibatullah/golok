package golok

import (
	"fmt"
	"strings"
)

type Styling struct {
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	EnableFg      bool
	EnableBg      bool
	Fg            uint32
	Bg            uint32
	Format        string
}

func (s *Styling) Apply(text string) *string {
	var (
		unstyled  bool
		ansiSeq   strings.Builder
		formatted string
	)
	unstyled = true

	if s.Bold || s.Italic || s.Underline || s.Strikethrough || s.EnableFg || s.EnableBg {
		unstyled = false
		ansiSeq.WriteString("\033[")
		if s.Bold {
			ansiSeq.WriteString("1;")
		}
		if s.Italic {
			ansiSeq.WriteString("3;")
		}
		if s.Underline {
			ansiSeq.WriteString("4;")
		}
		if s.Strikethrough {
			ansiSeq.WriteString("9;")
		}
		if s.EnableFg {
			ansiSeq.WriteString(fmt.Sprintf("38;2;%d;%d;%d;", (s.Fg>>16)&0xFF, (s.Fg>>8)&0xFF, s.Fg&0xFF))
		}
		if s.EnableBg {
			ansiSeq.WriteString(fmt.Sprintf("48;2;%d;%d;%d;", (s.Bg>>16)&0xFF, (s.Bg>>8)&0xFF, s.Bg&0xFF))
		}
	}

	if s.Format != "" {
		formatted = fmt.Sprintf(s.Format, text)
	} else {
		formatted = text
	}

	if !unstyled {
		formatted = fmt.Sprintf("%s%s\033[0m", ansiSeq.String()[:ansiSeq.Len()-1]+"m", formatted)
	}

	return &formatted
}
