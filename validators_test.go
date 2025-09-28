package wx

import (
	"testing"
)

// ---- Test struct ----
type DataTest struct {
	Name string `json:"code" check:"range:[3:8]"`
	// Email       string    `check:"regex:^.+@.+\\..+$;range:(40,);"`
	// Date        time.Time `check:"range:(2010-10-11 12:23:44Z,2100-10-11)"`
	// Description string    `check:"range:(20:)"`
	// Salary      *float64  `check:"range:[1:)"`
	// Code        *string
}

func TestXxx(t *testing.T) {
	// f := -1.77
	ret := Validators.Check(&DataTest{
		Name: "ADM",
		// Salary: &f,
		// Email:  "test@gmail.com",
		// Date:   time.Date(2010, 10, 11, 12, 23, 44, 2, time.UTC), //sao van kg hop le
	})
	t.Log(ret)

}

// func BenchmarkXxx(b *testing.B) {
// 	f := -1.77
// 	for i := 0; i < b.N; i++ {
// 		Validators.Check(&DataTest{
// 			Salary: &f,
// 			Email:  "test@gmail.com",
// 			Date:   time.Date(2010, 10, 11, 12, 23, 44, 2, time.UTC), //sao van kg hop le
// 		})

// 	}
// }
