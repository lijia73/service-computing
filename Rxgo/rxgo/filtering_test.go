package rxgo

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestDebounce(t *testing.T) {
	res := []int{}
	Just(1,2,3,4,5,6).Map(func(x int) int {
		switch x {
		case 1:
			time.Sleep(0 * time.Millisecond)
		case 2:
			time.Sleep(260 * time.Millisecond)
		case 3:
			time.Sleep(300 * time.Millisecond)
		case 4:
			time.Sleep(100 * time.Millisecond)
		case 5:
			time.Sleep(260 * time.Millisecond)
		case 6:
			time.Sleep(50 * time.Millisecond)
		}
		return x
	}).Debounce(250 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{1,2,4}, res, "Debounce Test Error!")
}

func TestDistinct(t *testing.T) {
	res := []int{}
	Just(1, 2, 1, 1, 2, 3, 4, 4).Distinct().Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{1, 2, 3, 4}, res, "Distinct Test Error!")
}

func TestElementAt(t *testing.T) {
	res := []int{}
	for i:=0;i<6;i++{
		Just(18,12,21,33,15,66).ElementAt(i).Subscribe(func(x int) {
			res = append(res, x)
		})
	}

	assert.Equal(t, []int{18,12,21,33,15,66}, res, "Distinct Test Error!")
}

func TestFirst(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).First().Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{18}, res, "Distinct Test Error!")
}

func TestIgnoreElements(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).IgnoreElements().Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{}, res, "IgnoreElements Test Error!")
}

func TestLast(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).Last().Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{66}, res, "Last Test Error!")
}

func TestSample(t *testing.T) {
	res := []int{}
	Just(1,2,3,4,5,6).Map(func(x int) int {
		switch x {
		case 1:
			time.Sleep(0 * time.Millisecond)
		case 2:
			time.Sleep(10 * time.Millisecond)
		case 3:
			time.Sleep(5 * time.Millisecond)
		case 4:
			time.Sleep(20 * time.Millisecond)
		case 5:
			time.Sleep(20 * time.Millisecond)
		case 6:
			time.Sleep(50 * time.Millisecond)
		}
		return x
	}).Sample(25 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{3,4,5,6}, res, "Sample Test Error!")
}


func TestSkip(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).Skip(3).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{33,15,66}, res, "Skip Test Error!")
}

func TestSkipLast(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).SkipLast(3).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{18,12,21}, res, "SkipLast Test Error!")
}

func TestTake(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).Take(2).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{18,12}, res, "Take Test Error!")
}

func TestTakeLast(t *testing.T) {
	res := []int{}
	Just(18,12,21,33,15,66).TakeLast(2).Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{15,66}, res, "TakeLast Test Error!")
}



