package main

import (
	"bytes"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func newProgressBar(total int64) *progressBar {
	return &progressBar{total: total, current: 0, percent: 0}
}

type progressBar struct {
	total   int64
	current int64
	percent int
}

func (pb *progressBar) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferDataEvent:
		if pb.total == event.TotalBytes {
			pb.current = event.ConsumedBytes
			percent := int(float64(pb.current) * 100 / float64(pb.total))
			if percent != pb.percent {
				pb.percent = percent
				pb.draw()
				if percent == 100 {
					fmt.Println()
				}
			}
		}
	default:
	}
}

func (pb *progressBar) draw() {
	num := pb.percent / 5
	fmt.Printf("\r[%v%v] %v / %v %v%%", mulitString("=", num), mulitString(" ", 20-num), pb.current, pb.total, pb.percent)
}

func mulitString(s string, num int) string {
	var buffer bytes.Buffer
	for i := 0; i < num; i++ {
		buffer.WriteString(s)
	}

	return buffer.String()
}
