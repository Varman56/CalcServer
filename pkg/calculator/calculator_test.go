package calculator

import (
	"math"
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
		wantError      bool
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
			wantError:      false,
		},
		{
			name:           "priority_brackets",
			expression:     "(2+2)*2",
			expectedResult: 8,
			wantError:      false,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
			wantError:      false,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
			wantError:      false,
		},
		{
			name:           "/0",
			expression:     "24/0",
			expectedResult: 0,
			wantError:      true,
		},
		{
			name:           "brackets 1",
			expression:     "(0)-(1)",
			expectedResult: -1,
			wantError:      false,
		},
		{
			name:           "brackets 2",
			expression:     "(-1)-(1)",
			expectedResult: -2,
			wantError:      false,
		},
		{
			name:           "brackets 3",
			expression:     "(-1-(-1))",
			expectedResult: 0,
			wantError:      false,
		},
		{
			name:           "t1",
			expression:     "(-1)+(0)",
			expectedResult: -1,
			wantError:      false,
		},
		{
			name:           "t2",
			expression:     "-(-(-1)+0)",
			expectedResult: -1,
			wantError:      false,
		},
		{
			name:           "t3",
			expression:     "(-(-0)+1)",
			expectedResult: 1,
			wantError:      false,
		},
		{
			name:           "t4",
			expression:     "((3/2)*2)*(7/7*5)",
			expectedResult: 15.0,
			wantError:      false,
		},
		{
			name:           "t5",
			expression:     "((3/2)*2)*(7/7-7)",
			expectedResult: -18.0,
			wantError:      false,
		},
		{
			name:           "error test",
			expression:     "((3/2)*2)*(7/(7-7))",
			expectedResult: -18.0,
			wantError:      true,
		},
		{
			name:           "t6",
			expression:     "((22.2/2)*3)*(-7)",
			expectedResult: -233.1,
			wantError:      false,
		},
		{
			name:           "big expression",
			expression:     "15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))*3-(2+(1+1))*(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))",
			expectedResult: -30.0721649485,
			wantError:      false,
		},
		{
			name:           "added mult",
			expression:     "(7)(5)",
			expectedResult: 35,
			wantError:      false,
		},
		{
			name:           "big expression without additional mult",
			expression:     "15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))3-(2+(1+1))(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))",
			expectedResult: -30.0721649485,
			wantError:      false,
		},
	}
	const EPS = 1e-9
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil && !testCase.wantError {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if !testCase.wantError && math.Abs(val-testCase.expectedResult) > EPS {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}
}
