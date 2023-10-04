package snowflake

import (
	"testing"
	"time"

	sclock "github.com/nomuyoshi/snowflake/clock"
)

// 同一ミリ秒で発行できるID数超えてしまうため回数指定して実行 -benchtime=4096x
func BenchmarkGenerateID(b *testing.B) {
	epoch := time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC)
	snowflake, _ := NewSnowflake(epoch, 5, 1)
	for i := 0; i < b.N; i++ {
		snowflake.Generate()
	}
}

func TestGenerateSequence(t *testing.T) {
	defer setRealClock()

	epoch := time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC)

	var datacenterID int64 = 5
	var machineID int64 = 0
	snowflake, _ := NewSnowflake(epoch, datacenterID, machineID)

	t.Run("同一ミリ秒内では連番がカウントアップ", func(t *testing.T) {
		now := time.Date(2023, 9, 28, 10, 10, 10, 0, time.UTC)
		clock = sclock.NewFakeClock(now)

		wants := make([]ID, 0, maxSequence)
		wantBegin := 153553470095360
		wantEnd := 153553470099455
		for i := wantBegin; i < wantEnd; i++ {
			wants = append(wants, ID(i))
		}

		for _, want := range wants {
			got := snowflake.Generate()
			if got != want {
				t.Errorf("failed. got = %d, want = %d", got, want)
			}
		}
	})

	t.Run("同一ミリ秒で連番が12bitを超えたらpanic", func(t *testing.T) {
		now := time.Date(2023, 9, 28, 10, 10, 10, 0, time.UTC)
		clock = sclock.NewFakeClock(now)

		defer func() {
			err := recover()
			if err == overMaxSequenceError {
				return
			}

			t.Errorf("unexpected error. got = %s, want = %s", err, overMaxSequenceError)
		}()

		for i := 0; i < (maxSequence + 2); i++ {
			snowflake.Generate()
		}
	})

	t.Run("ミリ秒ずれると連番がリセット", func(t *testing.T) {
		now := time.Date(2023, 9, 28, 10, 10, 10, 1000000, time.UTC)
		clock = sclock.NewFakeClock(now)

		want := ID(153553474289664)
		got := snowflake.Generate()
		if got != want {
			t.Errorf("failed. got = %d, want = %d", got, want)
		}
	})
}

func setRealClock() {
	clock = sclock.NewRealClock()
}
