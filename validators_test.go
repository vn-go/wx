package wx

import (
	"testing"
	"time"
)

// ---- Test struct ----
type DataTest struct {
	Name        string    `check:"range:(10,20);"`
	Email       string    `check:"regex:^.+@.+\\..+$;range:(40,);"`
	Date        time.Time `check:"range:(2010-10-11 12:23:44Z,2100-10-11)"`
	Description string    `check:"range:(20,)"`
}

func TestXxx(t *testing.T) {
	ret := Validators.Check(&DataTest{
		Email: "test@gmail.com",
		Date:  time.Date(2010, 10, 11, 12, 23, 44, 2, time.UTC), //sao van kg hop le
	})
	t.Log(ret)

}
func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validators.Check(&DataTest{
			Date: time.Now().UTC(),
		})

	}
}
