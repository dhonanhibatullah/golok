package golok

import (
	"fmt"
	"strings"
	"time"
)

type Component interface {
	render() *string
}

/*
 * ______________________________[ TEXT ]______________________________
 */
type Text struct {
	bind  *string
	style *Styling
}

func NewText(bind *string, style *Styling) *Text {
	return &Text{
		bind:  bind,
		style: style,
	}
}

func (t *Text) render() *string {
	return t.style.Apply(*t.bind)
}

/*
 * ______________________________[ PROGRESS BAR ]______________________________
 */
type ProgressBar struct {
	bind   *uint8
	length uint8
	style  *Styling
}

func NewProgressBar(bind *uint8, length uint8, style *Styling) *ProgressBar {
	var len uint8

	if length < 4 {
		len = 4
	} else if length > 75 {
		len = 75
	} else {
		len = length
	}

	return &ProgressBar{
		bind:   bind,
		length: len,
		style:  style,
	}
}

func (pb *ProgressBar) render() *string {
	var res strings.Builder
	res.WriteString("|")

	prog := *pb.bind
	if prog > 100 {
		prog = 100
	}

	perc := float64(prog) / 100.0
	nbar := uint8(perc * float64(pb.length))
	for i := uint8(0); i < pb.length; i++ {
		if i < nbar {
			res.WriteString("\u2588")
		} else {
			res.WriteString(" ")
		}
	}
	res.WriteString("|")
	res.WriteString(fmt.Sprintf("%3d%%", prog))

	return pb.style.Apply(res.String())
}

/*
 * ______________________________[ TIMESTAMP ]______________________________
 */

type TimestampFmt uint8

const (
	Second TimestampFmt = iota
	Milli
	Micro
	Nano
)

type Timestamp struct {
	format TimestampFmt
	offset int64
	style  *Styling
}

func NewTimestamp(format TimestampFmt, unixOffset int64, style *Styling) *Timestamp {
	return &Timestamp{
		format: format,
		offset: unixOffset,
		style:  style,
	}
}

func (ts *Timestamp) render() *string {
	switch ts.format {
	case Second:
		timestamp := time.Now().Unix() - ts.offset
		return ts.style.Apply(fmt.Sprintf("%10d", timestamp))
	case Milli:
		timestamp := time.Now().UnixMilli() - ts.offset
		return ts.style.Apply(fmt.Sprintf("%10.3f", float64(timestamp)/1000.0))
	case Micro:
		timestamp := time.Now().UnixMicro() - ts.offset
		return ts.style.Apply(fmt.Sprintf("%10.6f", float64(timestamp)/1000000.0))
	case Nano:
		timestamp := time.Now().UnixNano() - ts.offset
		return ts.style.Apply(fmt.Sprintf("%10.9f", float64(timestamp)/1000000000.0))
	default:
		tmp := ""
		return &tmp
	}
}

/*
 * ______________________________[ DATETIME ]______________________________
 */

type DatetimeFmt uint8

const (
	YYYYnMMnDD_Time DatetimeFmt = iota
	YYYYnMMnDD_TimeMilli
	YYYYnMMnDD_TimeMicro
	YYYYnMMnDD_TimeNano
	YYYYnDDnMM_Time
	YYYYnDDnMM_TimeMilli
	YYYYnDDnMM_TimeMicro
	YYYYnDDnMM_TimeNano
	DDnMMnYYYY_Time
	DDnMMnYYYY_TimeMilli
	DDnMMnYYYY_TimeMicro
	DDnMMnYYYY_TimeNano
	MMnDDnYYYY_Time
	MMnDDnYYYY_TimeMilli
	MMnDDnYYYY_TimeMicro
	MMnDDnYYYY_TimeNano
	YYYY_Month_DD_Time
	YYYY_Month_DD_TimeMilli
	YYYY_Month_DD_TimeMicro
	YYYY_Month_DD_TimeNano
	YYYY_DD_Month_Time
	YYYY_DD_Month_TimeMilli
	YYYY_DD_Month_TimeMicro
	YYYY_DD_Month_TimeNano
	DD_Month_YYYY_Time
	DD_Month_YYYY_TimeMilli
	DD_Month_YYYY_TimeMicro
	DD_Month_YYYY_TimeNano
	Month_DD_YYYY_Time
	Month_DD_YYYY_TimeMilli
	Month_DD_YYYY_TimeMicro
	Month_DD_YYYY_TimeNano
)

var datetimeFmtStr = map[DatetimeFmt]string{
	YYYYnMMnDD_Time:         "2006-01-02 15:04:05",
	YYYYnMMnDD_TimeMilli:    "2006-01-02 15:04:05.000",
	YYYYnMMnDD_TimeMicro:    "2006-01-02 15:04:05.000000",
	YYYYnMMnDD_TimeNano:     "2006-01-02 15:04:05.000000000",
	YYYYnDDnMM_Time:         "2006-02-01 15:04:05",
	YYYYnDDnMM_TimeMilli:    "2006-02-01 15:04:05.000",
	YYYYnDDnMM_TimeMicro:    "2006-02-01 15:04:05.000000",
	YYYYnDDnMM_TimeNano:     "2006-02-01 15:04:05.000000000",
	DDnMMnYYYY_Time:         "01-02-2006 15:04:05",
	DDnMMnYYYY_TimeMilli:    "01-02-2006 15:04:05.000",
	DDnMMnYYYY_TimeMicro:    "01-02-2006 15:04:05.000000",
	DDnMMnYYYY_TimeNano:     "01-02-2006 15:04:05.000000000",
	MMnDDnYYYY_Time:         "02-01-2006 15:04:05",
	MMnDDnYYYY_TimeMilli:    "02-01-2006 15:04:05.000",
	MMnDDnYYYY_TimeMicro:    "02-01-2006 15:04:05.000000",
	MMnDDnYYYY_TimeNano:     "02-01-2006 15:04:05.000000000",
	YYYY_Month_DD_Time:      "2006 January 02 15:04:05",
	YYYY_Month_DD_TimeMilli: "2006 January 02 15:04:05.000",
	YYYY_Month_DD_TimeMicro: "2006 January 02 15:04:05.000000",
	YYYY_Month_DD_TimeNano:  "2006 January 02 15:04:05.000000000",
	YYYY_DD_Month_Time:      "2006 02 January 15:04:05",
	YYYY_DD_Month_TimeMilli: "2006 02 January 15:04:05.000",
	YYYY_DD_Month_TimeMicro: "2006 02 January 15:04:05.000000",
	YYYY_DD_Month_TimeNano:  "2006 02 January 15:04:05.000000000",
	DD_Month_YYYY_Time:      "02 January 2006 15:04:05",
	DD_Month_YYYY_TimeMilli: "02 January 2006 15:04:05.000",
	DD_Month_YYYY_TimeMicro: "02 January 2006 15:04:05.000000",
	DD_Month_YYYY_TimeNano:  "02 January 2006 15:04:05.000000000",
	Month_DD_YYYY_Time:      "January 02 2006 15:04:05",
	Month_DD_YYYY_TimeMilli: "January 02 2006 15:04:05.000",
	Month_DD_YYYY_TimeMicro: "January 02 2006 15:04:05.000000",
	Month_DD_YYYY_TimeNano:  "January 02 2006 15:04:05.000000000",
}

type Datetime struct {
	format string
	style  *Styling
}

func NewDatetime(format DatetimeFmt, style *Styling) *Datetime {
	return &Datetime{
		format: datetimeFmtStr[format],
		style:  style,
	}
}

func (dt *Datetime) render() *string {
	return dt.style.Apply(time.Now().Format(dt.format))
}
