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

	fn := os.Getenv("CURRICULUM_FILE1")
	w, err := newWeekFromFile(fn)
	is.NoErr(err)

	l, err := helpers.LoadLocation()
	is.NoErr(err)

	d := NewSpecificDay(*w, time.Date(2020, 9, 11, 0, 0, 0, 0, l))
	is.Equal(d[1].Name, "Технології розроблення програмних систем")
}

func TestWeekend(t *testing.T) {
	is := is.New(t)

	fn := os.Getenv("CURRICULUM_FILE1")
	w, err := newWeekFromFile(fn)
	is.NoErr(err)

	l, err := helpers.LoadLocation()
	is.NoErr(err)

	d := NewSpecificDay(*w, time.Date(2020, 9, 6, 0, 0, 0, 0, l))
	is.Equal(len(d), 0)
}

func newWeekFromFile(fn string) (*Week, error) {
	f, err := os.Open(path.Join("..", "assets", fn))
	if err != nil {
		return nil, err
	}

	return NewWeek(io.Reader(f))
}
