package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"

	sclock "snowflake/clock"
)

type ID int64

type Snowflake struct {
	mu sync.Mutex

	epoch         int64
	datacenterID  int64
	machineID     int64
	sequence      int64
	prevTimestamp int64
}

const (
	sequenceBits     = 12
	machineIDBits    = 5
	datacenterIDBits = 5
	timeStampBits    = 41

	maxDatacenterID = 31   // 5 bits
	maxMachineID    = 31   // 5 bits
	maxSequence     = 4095 // 12 bits

	machineShift    = sequenceBits                       // 12
	datacenterShift = machineShift + machineIDBits       // 17
	timeStampShift  = datacenterShift + datacenterIDBits //
)

var overMaxSequenceError = errors.New("reached the maximum number of IDs during the same milliseconds")

var clock sclock.Clock = sclock.NewRealClock()

func NewSnowflake(epochTime time.Time, datacenterID, machineID int64) (*Snowflake, error) {
	if datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("datacenterID must be between 0 and %d", maxDatacenterID)
	}

	if machineID > maxMachineID {
		return nil, fmt.Errorf("machineID must be between 0 and %d", maxMachineID)
	}

	return &Snowflake{
		mu:           sync.Mutex{},
		epoch:        epochTime.UnixMilli(),
		datacenterID: datacenterID,
		machineID:    machineID,
	}, nil
}

func (s *Snowflake) Generate() ID {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := clock.Now().UnixMilli() - s.epoch
	if timestamp == s.prevTimestamp { // ミリ秒単位で同時刻なら連番をインクリメント
		s.sequence++
		if s.sequence > maxSequence {
			panic(overMaxSequenceError) // 連番が足りなくなったらpanic
		}
	} else { // ミリ秒ズレたら連番をリセット
		s.sequence = 0
		s.prevTimestamp = timestamp
	}

	return ID((timestamp << timeStampShift) | (s.datacenterID << datacenterShift) | (s.machineID << machineShift) | s.sequence)
}
