package curriculum

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/yevhenshymotiuk/ap-curriculum-bot/helpers"
)

func TestWorkday(t *testing.T) {
	is := is.New(t)

	fn := os.Getenv("CURRICULUM_FILE")
	w, err := newWeekFromFile(fn)
	is.NoErr(err)

	l, err := helpers.LoadLocation()
	is.NoErr(err)

	d := NewSpecificDay(*w, time.Date(2020, 9, 21, 0, 0, 0, 0, l))
	dps1, dps2 := d[0], d[1]

	is.Equal(dps1[2].Name, "Методи штучного інтелекту")
	is.Equal(dps1[2].Name, dps2[2].Name)
}

func TestWeekend(t *testing.T) {
	is := is.New(t)

	fn := os.Getenv("CURRICULUM_FILE")
	w, err := newWeekFromFile(fn)
	is.NoErr(err)

	l, err := helpers.LoadLocation()
	is.NoErr(err)

	d := NewSpecificDay(*w, time.Date(2020, 9, 6, 0, 0, 0, 0, l))
	is.Equal(len(d[0]), 0)
}

func newWeekFromFile(fn string) (*Week, error) {
	f, err := os.Open(path.Join("..", "assets", fn))
	if err != nil {
		return nil, err
	}

	return NewWeek(io.Reader(f))
}
