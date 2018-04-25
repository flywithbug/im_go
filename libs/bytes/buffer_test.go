package bytes
import (
	"testing"
)

func TestBuffer(t *testing.T) {
	p := NewPool(2, 10)
	b := p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	count := 0
	for   {
		if count  > 100{
			break
		}
		b = p.Get()
		if b.Bytes() == nil || len(b.Bytes()) == 0 {
			t.FailNow()
		}
		count++
	}
}
