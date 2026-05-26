package engine

import (
	"testing"
	"time"
)

func TestTopN_BasicRanking(t *testing.T) {
	e := NewEngine()

	e.Add("кроссовки")
	e.Add("кроссовки")
	e.Add("кроссовки")

	e.Add("куртка")
	e.Add("куртка")

	e.Add("футболка")

	e.recalculate()

	top := e.Top(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 items, got %d", len(top))
	}

	if top[0].Query != "кроссовки" {
		t.Errorf("expected 'кроссовки' at top, got '%s'", top[0].Query)
	}

	if top[1].Query != "куртка" {
		t.Errorf("expected 'куртка' at second, got '%s'", top[1].Query)
	}
}

func TestTopN_LimitRespected(t *testing.T) {
	e := NewEngine()

	for _, q := range []string{"a", "b", "c", "d", "e"} {
		e.Add(q)
	}

	e.recalculate()

	top := e.Top(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 items, got %d", len(top))
	}
}

func TestTopN_EmptyEngine(t *testing.T) {
	e := NewEngine()
	e.recalculate()

	top := e.Top(10)
	if len(top) != 0 {
		t.Fatalf("expected empty top, got %d items", len(top))
	}
}

func TestStoplist_FiltersFromTop(t *testing.T) {
	e := NewEngine()

	e.Add("плохой сёрч")
	e.Add("плохой сёрч")
	e.Add("плохой сёрч")

	e.Add("хороший сёрч")
	e.Add("хороший сёрч")

	e.AddToStoplist("плохой сёрч")
	e.recalculate()

	top := e.Top(10)
	for _, entry := range top {
		if entry.Query == "плохой сёрч" {
			t.Errorf("'плохой сёрч' should be filtered by stoplist")
		}
	}
}

func TestStoplist_RemoveRestoredWord(t *testing.T) {
	e := NewEngine()

	e.Add("сомнительный сёрч")
	e.Add("сомнительный сёрч")

	e.AddToStoplist("сомнительный сёрч")

	e.recalculate()

	top := e.Top(10)
	for _, entry := range top {
		if entry.Query == "сомнительный сёрч" {
			t.Error("'сомнительный сёрч' should be filtered by stoplist")
		}
	}

	e.RemoveFromStoplist("сомнительный сёрч")

	e.recalculate()

	top = e.Top(10)
	found := false
	for _, entry := range top {
		if entry.Query == "сомнительный сёрч" {
			found = true
		}
	}

	if !found {
		t.Error("'сомнительный сёрч' should be back in top after being removed from stoplist")
	}
}

func TestSlidingWindow_DropsOldBuckets(t *testing.T) {
	e := NewEngine()

	e.Add("асикс")
	e.recalculate()

	top := e.Top(10)
	if len(top) == 0 {
		t.Fatal("expected 'асикс' in top")
	}

	for i := 0; i < bucketCount; i++ {
		e.rotate()
	}
	e.recalculate()

	top = e.Top(10)
	if len(top) != 0 {
		t.Errorf("expected empty top after window expired, got %d items", len(top))
	}
}

func TestTopN_CountsCorrect(t *testing.T) {
	e := NewEngine()

	for i := 0; i < 50; i++ {
		e.Add("худи кэжуал")
	}

	for i := 0; i < 30; i++ {
		e.Add("худи нишевое")
	}

	e.recalculate()

	top := e.Top(2)

	if top[0].Count != 50 {
		t.Errorf("expected count 50, got %d", top[0].Count)
	}

	if top[1].Count != 30 {
		t.Errorf("expected count 30, got %d", top[1].Count)
	}
}

func TestEngine_ConcurrentAccess(t *testing.T) {
	e := NewEngine()

	done := make(chan struct{})

	go func() {
		for i := 0; i < 1000; i++ {
			e.Add("кроссовки")
		}
		close(done)
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			e.recalculate()
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			e.Top(10)
			time.Sleep(time.Millisecond)
		}
	}()

	<-done
}

func BenchmarkEngine_Add(b *testing.B) {
	e := NewEngine()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Add("кроссовки")
		}
	})
}

func BenchmarkEngine_Top(b *testing.B) {
	e := NewEngine()

	for i := 0; i < 10000; i++ {
		e.Add("асикс")
		e.Add("нью бэленс")
		e.Add("найк")
	}

	e.recalculate()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Top(10)
		}
	})
}

func BenchmarkEngine_AddAndTop(b *testing.B) {
	e := NewEngine()

	go func() {
		for {
			e.recalculate()
			time.Sleep(time.Second)
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Add("кроссовки")
			e.Top(10)
		}
	})
}
