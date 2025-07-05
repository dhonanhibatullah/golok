# **golok**

Want to have a logging like this?

Don't worry, golok has you covered!

## **A. Installation**

```bash
go get github.com/dhonanhibatullah/golok
```

## **B. Introduction**

Golok consists of 4 essential parts:

### **1. Golok Instance**

The Golok instance is the main part to initiate the Golok's logging mechanism. Golok instance is used to handle multiple profiles. We initiate it with:

```go
glk := golok.NewGolok()
defer glk.Close()
```

### **2. Profile**

A profile is **a line of log** which built by components. So in order to create a profile, we have to combine multiple components. Profile can be created from the Golok instance and provide the index to describe in which line of log we want to place the profile with:

```go
glkProfile := glk.NewProfile(index)
```

### **3. Component**

Currently, there are 4 types of component: `Text`, `ProgressBar`, `Timestamp`, and `Datetime`. A component requires a **bind**, which is a pointer to a variable we want the component value is (not required for `Timestamp` and `Datetime`). Here is the example of initiating each component and put it into the profile.

```go
textVal := "Hello Dhonan!"
progressVal := uint8(34)

myText := golok.NewText(&textVal, style1)
myProgressBar := golok.NewProgressBar(&progressVal, 40, style2)
myTimestamp := golok.NewTimestamp(golok.Second, 0, style3)
myDatetime := golok.NewDatetime(golok.YYYYnMMnDD_TimeMilli, style4)

glkProfile.AddComponent(0, myText)
glkProfile.AddComponent(1, myProgressBar)
glkProfile.AddComponent(2, myTimestamp)
glkProfile.AddComponent(3, myDatetime)
```

We will talk about creating a style right after this.

### **4. Styling**

A style consists of 9 properties, each has the effect based on its name.

```
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
```

This is the example of a style:

```go
style1 := &Styling{
    Bold:     true,
    EnableFg: true,
    EnableBg: true,
    Fg:       0xFFFFFF,
    Bg:       0x009E91,
    Format:   "%s",
}
```

Note: `Format` is how place the resulted string. You can put anything in the `Format` as long as there is one and only one format directive, which is the `%s`. If you do not want any format, leave it unassigned.

### **5. Complete Example**

Here is a fun complete example you could directly try.

```go
package main

import (
	"context"
	"time"

	"github.com/dhonanhibatullah/golok"
)

func createStyle(fg uint32, bg uint32, format string) *golok.Styling {
	return &golok.Styling{
		Bold:     true,
		EnableFg: true,
		EnableBg: true,
		Fg:       fg,
		Bg:       bg,
		Format:   format,
	}
}

func myRoutine(cancel context.CancelFunc, p *golok.Profile) {
	defer p.Close()
	defer cancel()

	textVal := ""
	progressVal := uint8(34)

	myText := golok.NewText(&textVal, createStyle(0xA15E00, 0xFFFFFF, " >> %s "))
	myProgressBar := golok.NewProgressBar(&progressVal, 40, createStyle(0xBDE300, 0x000000, " %s "))
	myTimestamp := golok.NewTimestamp(golok.Second, 0, createStyle(0x0AFF3B, 0x000000, " [ %s ] "))
	myDatetime := golok.NewDatetime(golok.YYYYnMMnDD_TimeMilli, createStyle(0x00CFCF, 0x000000, " < %s > "))

	p.AddComponent(3, myText)
	p.AddComponent(2, myProgressBar)
	p.AddComponent(0, myTimestamp)
	p.AddComponent(1, myDatetime)

	textVal = "Downloading..."
	progressVal = 25
	p.Render()
	time.Sleep(time.Duration(1) * time.Second)

	textVal = "Installing..."
	progressVal = 50
	p.Render()
	time.Sleep(time.Duration(1) * time.Second)

	textVal = "Starting..."
	progressVal = 75
	p.Render()
	time.Sleep(time.Duration(1) * time.Second)

	textVal = "Finishing..."
	progressVal = 100
	p.Render()
	time.Sleep(time.Duration(1) * time.Second)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	glk := golok.NewGolok()
	defer glk.Close()

	go myRoutine(cancel, glk.NewProfile(0))
	<-ctx.Done()
}
```

You should get something like this:
