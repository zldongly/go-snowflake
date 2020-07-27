package snowflake

import (
    "testing"
    "time"
)

var (
    epoch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestGetTime(t *testing.T) {
    now := time.Now()

    since := time.Since(epoch).Nanoseconds() / 1000000

    between := (now.UnixNano()  - epoch.UnixNano())/ 1000000

    t.Log(since)
    t.Log(between)
}

func TestMax(t *testing.T) {
    bitSize := 10
    numMax := -1 ^ (-1 << bitSize)
    t.Logf("%d", numMax)
}

func TestNode_NextId(t *testing.T) {
    node, _ := NewNode(41, 10, 12, 1, epoch)
    t.Logf("%#v", node.conf)
    t.Logf("%#v", node)
    id := node.NextId()
    t.Logf("%#v", node)
    t.Logf("id:%v", id)

    check := node.step + (node.machineId << node.conf.machineShift)
    check += node.timeMsStamp << node.conf.timeShift
    if id != check {
        t.Fatalf("\nid:\t\t%d\ncheck:\t%v", id, check)
    }
}


func BenchmarkNode_NextId(b *testing.B) {
    // 56.0 ns/op
    node, _ := NewNode(41, 10, 12, 1, epoch)
    //var id int64
    for i := 0; i < b.N; i++ {
        _ = node.NextId()
    }
}
