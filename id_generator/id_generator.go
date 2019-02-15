package id_generator

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"hash/adler32"
	"io"
	"os"
	"sync/atomic"
	"time"
)

var machineId = getMachineId()
var processId = getPid()

func getMachineId() []byte {
	id := make([]byte, 3)

	hostName, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic("无法读取计算机名")
		}

		return id
	}

	hash := adler32.New()
	hash.Write([]byte(hostName))

	copy(id, hash.Sum(nil))

	return id
}

func getPid() []byte {
	pid := make([]byte, 4)
	pidNo := os.Getpid()

	pid[0] = byte(pidNo >> 24)
	pid[1] = byte(pidNo >> 16)
	pid[2] = byte(pidNo >> 8)
	pid[3] = byte(pidNo)

	return pid
}

type IdGenerator struct {
	counter uint32
}

func (generator *IdGenerator) Generate() string {
	buffer := bytes.Buffer{}
	buffer.Write(generator.getUnix())
	buffer.Write(machineId)
	buffer.Write(processId)
	buffer.Write(generator.getCount())

	return hex.EncodeToString(buffer.Bytes())
}

func (generator *IdGenerator) getUnix() []byte {
	bytes := make([]byte, 8)

	binary.BigEndian.PutUint64(bytes, uint64(time.Now().Unix()))

	return bytes[4:]
}

func (generator *IdGenerator) getCount() []byte {
	count := atomic.AddUint32(&generator.counter, 1)
	countRes := make([]byte, 4)

	countRes[0] = byte(count >> 24)
	countRes[1] = byte(count >> 16)
	countRes[2] = byte(count >> 8)
	countRes[3] = byte(count)

	return countRes
}
