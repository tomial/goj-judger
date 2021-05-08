package judger

import (
	"io/ioutil"
	"regexp"
	"strconv"
)

func measure() (ut, st, rt float64, ram uint64) {
	res, err := ioutil.ReadFile("./result")
	if err != nil {
		panic(err)
	}

	target := string(res)

	digit := regexp.MustCompile(`(\d*[.]\d+)+|(\d+kB)+`)
	result := digit.FindAllString(target, -1)

	ut, err = strconv.ParseFloat(result[0], 32)
	if err != nil {
		panic(err)
	}

	st, err = strconv.ParseFloat(result[1], 32)
	if err != nil {
		panic(err)
	}

	rt, err = strconv.ParseFloat(result[2], 32)
	if err != nil {
		panic(err)
	}

	ram, err = strconv.ParseUint(result[3][:len(result[3])-2], 10, 64)
	if err != nil {
		panic(err)
	}

	return
}
