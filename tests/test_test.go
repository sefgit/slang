package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"slang"
	"io/ioutil"
)

func TestTestOperator__TrivialTests(t *testing.T) {
	a := assert.New(t)
	succs, fails, err := slang.TestOperator("test_data/voidOp_test.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(1, succs)
	a.Equal(0, fails)
}

func TestTestOperator__SimpleFail(t *testing.T) {
	a := assert.New(t)
	succs, fails, err := slang.TestOperator("test_data/voidOp_corruptTest.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(0, succs)
	a.Equal(1, fails)
}

func TestTestOperator__ComplexTests(t *testing.T) {
	a := assert.New(t)
	succs, fails, err := slang.TestOperator("test_data/nested_op/usingSubCustomOpDouble_test.json", ioutil.Discard, true)
	a.Nil(err)
	a.Equal(2, succs)
	a.Equal(0, fails)
}