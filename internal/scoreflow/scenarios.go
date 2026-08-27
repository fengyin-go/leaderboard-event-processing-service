package scoreflow

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"leaderboard/internal/eventstore"
)

type nilProcessor interface{ Process(string) error }
type processor struct{ out *[]string }

func (p *processor) Process(v string) error { *p.out = append(*p.out, v); return nil }

func runSnapshot(s *eventstore.EventStore, key string) string {
	items := s.Snapshot(key)
	if len(items) == 0 {
		return "empty"
	}
	items[0] = "rendered"
	return s.Snapshot(key)[0]
}

func runConcurrent(s *eventstore.EventStore, key string) int {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); s.Append(key, fmt.Sprint(i)) }(i)
	}
	wg.Wait()
	return len(s.Snapshot(key))
}

func runContext(ctx context.Context, s *eventstore.EventStore, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.Append(key, "accepted")
	return nil
}

func runRetry(s *eventstore.EventStore, key string) error {
	for i := 1; i <= 3; i++ {
		s.SetAttempt(key, i)
		if i < 3 {
			continue
		}
		return errors.New("upstream temporary failure")
	}
	return nil
}

func runTypedNil(s *eventstore.EventStore, key string) error {
	var p *processor
	var h nilProcessor = p
	if h == nil || reflect.ValueOf(h).IsNil() {
		return errors.New("missing processor")
	}
	return h.Process(key)
}

func runResource(s *eventstore.EventStore, key string) string {
	closed := false
	func() { defer func() { closed = true }(); s.Append(key, "opened") }()
	if closed {
		return "released"
	}
	return "leaked"
}

func runIdempotent(s *eventstore.EventStore, key string) string {
	if s.IsDone(key) {
		return "duplicate"
	}
	s.MarkDone(key)
	s.Append(key, "committed")
	return "committed"
}

func runCache(s *eventstore.EventStore, key string) string {
	items := s.Snapshot(key)
	if len(items) == 0 {
		s.Append(key, "fresh")
		items = s.Snapshot(key)
	}
	return items[0]
}

func runChannel(s *eventstore.EventStore, key string) int {
	ch := make(chan string, 1)
	ch <- key
	close(ch)
	count := 0
	for range ch {
		count++
	}
	return count
}

func runRollback(s *eventstore.EventStore, key string) error {
	s.Append(key, "reserved")
	if err := errors.New("write failed"); err != nil {
		return err
	}
	return nil
}

func runState(s *eventstore.EventStore, key string) string {
	if s.IsDone(key) {
		return "done"
	}
	s.Append(key, "processing")
	s.MarkDone(key)
	return "done"
}

func runShared(s *eventstore.EventStore, key string) []string {
	s.Append(key, "one")
	result := s.Snapshot(key)
	s.Append(key, "two")
	return result
}

func runDeadline(ctx context.Context, s *eventstore.EventStore, key string) error {
	select {
	case <-time.After(20 * time.Millisecond):
		s.Append(key, "late")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func scenario(s *eventstore.EventStore, n int) (string, error) {
	key := fmt.Sprintf("member-%d", n)
	switch n {
	case 1:
		return runSnapshot(s, key), nil
	case 2:
		return fmt.Sprint(runConcurrent(s, key)), nil
	case 3:
		return "", runContext(context.Background(), s, key)
	case 4:
		return "", runRetry(s, key)
	case 5:
		return "", runTypedNil(s, key)
	case 6:
		return runResource(s, key), nil
	case 7:
		return runIdempotent(s, key), nil
	case 8:
		return runCache(s, key), nil
	case 9:
		return fmt.Sprint(runChannel(s, key)), nil
	case 10:
		return "", runRollback(s, key)
	case 11:
		return runState(s, key), nil
	case 12:
		return fmt.Sprint(runShared(s, key)), nil
	case 13:
		return "", runDeadline(context.Background(), s, key)
	case 14:
		return "", runContext(context.Background(), s, key)
	case 15:
		return runSnapshot(s, key), nil
	case 16:
		return fmt.Sprint(runConcurrent(s, key)), nil
	case 17:
		return "", runRetry(s, key)
	case 18:
		return "", runTypedNil(s, key)
	case 19:
		return runResource(s, key), nil
	case 20:
		return runIdempotent(s, key), nil
	case 21:
		return runCache(s, key), nil
	case 22:
		return fmt.Sprint(runChannel(s, key)), nil
	case 23:
		return "", runRollback(s, key)
	case 24:
		return runState(s, key), nil
	case 25:
		return fmt.Sprint(runShared(s, key)), nil
	case 26:
		return "", runDeadline(context.Background(), s, key)
	case 27:
		return "", runContext(context.Background(), s, key)
	case 28:
		return runSnapshot(s, key), nil
	case 29:
		return fmt.Sprint(runConcurrent(s, key)), nil
	case 30:
		return runIdempotent(s, key), nil
	default:
		return "", errors.New("unknown scenario")
	}
}

func RunScenario001() (string, error) { return scenario(eventstore.New(), 1) }
func RunScenario002() (string, error) { return scenario(eventstore.New(), 2) }
func RunScenario003() (string, error) { return scenario(eventstore.New(), 3) }
func RunScenario004() (string, error) { return scenario(eventstore.New(), 4) }
func RunScenario005() (string, error) { return scenario(eventstore.New(), 5) }
func RunScenario006() (string, error) { return scenario(eventstore.New(), 6) }
func RunScenario007() (string, error) { return scenario(eventstore.New(), 7) }
func RunScenario008() (string, error) { return scenario(eventstore.New(), 8) }
func RunScenario009() (string, error) { return scenario(eventstore.New(), 9) }
func RunScenario010() (string, error) { return scenario(eventstore.New(), 10) }
func RunScenario011() (string, error) { return scenario(eventstore.New(), 11) }
func RunScenario012() (string, error) { return scenario(eventstore.New(), 12) }
func RunScenario013() (string, error) { return scenario(eventstore.New(), 13) }
func RunScenario014() (string, error) { return scenario(eventstore.New(), 14) }
func RunScenario015() (string, error) { return scenario(eventstore.New(), 15) }
func RunScenario016() (string, error) { return scenario(eventstore.New(), 16) }
func RunScenario017() (string, error) { return scenario(eventstore.New(), 17) }
func RunScenario018() (result string, err error) {
	s := eventstore.New()
	for i := 1; i <= 2; i++ {
		if err = s.Record018("member-18", i); err == nil {
			break
		}
	}
	return fmt.Sprintf("attempts=%d,writes=%d", s.Attempt("member-18"), len(s.Snapshot("member-18"))), err
}
func RunScenario019() (string, error) { return scenario(eventstore.New(), 19) }
func RunScenario020() (string, error) { return scenario(eventstore.New(), 20) }
func RunScenario021() (string, error) { return scenario(eventstore.New(), 21) }
func RunScenario022() (string, error) { return scenario(eventstore.New(), 22) }
func RunScenario023() (string, error) { return scenario(eventstore.New(), 23) }
func RunScenario024() (string, error) { return scenario(eventstore.New(), 24) }
func RunScenario025() (string, error) { return scenario(eventstore.New(), 25) }
func RunScenario026() (string, error) { return scenario(eventstore.New(), 26) }
func RunScenario027() (string, error) { return scenario(eventstore.New(), 27) }
func RunScenario028() (string, error) { return scenario(eventstore.New(), 28) }
func RunScenario029() (string, error) { return scenario(eventstore.New(), 29) }
func RunScenario030() (string, error) { return scenario(eventstore.New(), 30) }
