package snowflake

import (
    "sync"
    "time"
)

type config struct {
    epoch time.Time // 开始时间

    // 左移字节数
    timeShift    uint8 // 时间
    machineShift uint8 // 设备号

    timeMax    int64
    machineMax int64
    stepMax    int64
}

type node struct {
    conf *config
    mu   sync.Mutex

    timeMsStamp int64 // ms时间戳
    machineId   int64
    step        int64
}

func NewNode(timeBits, machineBits, stepBits uint8, machineId int64, epoch time.Time) (node, error) {
    conf := &config{
        epoch: epoch,

        timeShift:    machineBits + stepBits,
        machineShift: stepBits,

        timeMax:    -1 ^ (-1 << timeBits),
        machineMax: -1 ^ (-1 << machineBits),
        stepMax:    -1 ^ (-1 << stepBits),
    }

    // TODO 2020/7/27 21:45 判断 bit 大小和 machineId

    return node{
            conf:      conf,
            machineId: machineId,
        },
        nil
}

func (n *node) NextId() int64 {
    n.mu.Lock()
    defer n.mu.Unlock()

    timeStamp := n.conf.getTimeStamp()

    // 相同时间
    if timeStamp == n.timeMsStamp {
        n.step++
        if n.step&n.conf.stepMax == 0 {
            n.step = 0
            for timeStamp > n.timeMsStamp {
                timeStamp = n.conf.getTimeStamp()
            }
        }
        n.timeMsStamp = timeStamp

        return n.generate()
    }

    if timeStamp < n.timeMsStamp { // 时间回滚
        return 0
    }

    n.timeMsStamp = timeStamp
    n.step = 0

    return n.generate()
}

// 时间戳
func (c config) getTimeStamp() int64 {
    return time.Since(c.epoch).Nanoseconds() / 1000000
}

func (n node) generate() int64 {
    return (n.timeMsStamp << (n.conf.timeShift)) |
        (n.machineId << n.conf.machineShift) |
        (n.step)
}
