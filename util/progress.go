package util

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// progress bar
func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{total: total, current: 0, percent: 0}
}

type ProgressBar struct {
	total   int64
	current int64
	percent int
}

func (pb *ProgressBar) ProgressChanged(event *oss.ProgressEvent) {
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
func (pb *ProgressBar) draw() {
	num := pb.percent / 5
	fmt.Printf("\r[%v%v] %v / %v %v%%", multiString("=", num), multiString(" ", 20-num), pb.current, pb.total, pb.percent)
}
func multiString(s string, num int) string {
	var buffer bytes.Buffer
	for i := 0; i < num; i++ {
		buffer.WriteString(s)
	}

	return buffer.String()
}
