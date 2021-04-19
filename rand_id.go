package main

import (
    "crypto/md5"
    "encoding/binary"
    "encoding/hex"
    "fmt"
    "os"
    "sync/atomic"
    "time"
    "sync"
    "math/rand"
)


func GetRandId() string {
    var b [12]byte
    // 时间
    binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
    machineId := readMachineId()
    // 机器id
    b[4] = machineId[0]
    b[5] = machineId[1]
    b[6] = machineId[2]
    // 进程id
    // pid := os.Getpid()
    // 随机数
    pid := rand.Intn(100)
    b[7] = byte(pid >> 8)
    b[8] = byte(pid)
    var objectIdCounter uint32 = 0
    // 原子自增id 
    i := atomic.AddUint32(&objectIdCounter, 1)
    b[9] = byte(i >> 16)
    b[10] = byte(i >> 8)
    b[11] = byte(i)
    return hex.EncodeToString([]byte(b[:]))
}

func readMachineId() []byte {
    var sum [3]byte
    id := sum[:]
    hostname, err := os.Hostname()
    if err != nil {
        panic(fmt.Errorf("cannot get hostname: %v", err))
    }
    hw := md5.New()
    hw.Write([]byte(hostname))
    copy(id, hw.Sum(nil))
    return id
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        go func() {
            wg.Add(1)
            id := GetRandId()
            fmt.Println(id)
            wg.Done()
        }()
    }
    wg.Wait()
}