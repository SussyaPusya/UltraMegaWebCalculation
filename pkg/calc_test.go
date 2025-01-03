package pkg_test

import (
	"testing"

	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg"
)

// Ebeishiy test for the calc func
func TestCalc(t *testing.T) {

	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "*/",
			expression:     "2*6/2",
			expectedResult: 6,
		},
		{
			name:           "*(/)",
			expression:     "3*(6/2)",
			expectedResult: 9,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := pkg.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("succses test %s returned error ", testCase.expression)

			}
			if val != testCase.expectedResult {
				t.Fatalf("incorrect number should be %f, has been %f", testCase.expectedResult, val)
			}
		})
	}
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "simple",
			expression: "1+1*",
		},
		{
			name:       "priority",
			expression: "2+2**2",
		},
		{
			name:       "priority",
			expression: "((2+2-*(2",
		},
		{
			name:       "/",
			expression: "",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := pkg.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
		})
	}
}
