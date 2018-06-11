package semver

import (
	"testing"
)

func Test_trimDistribution(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			input:    "0.15.0~rc.1~git20180507.28967e3+ds",
			expected: "0.15.0",
		},
		{
			input:    "foobar",
			expected: "foobar",
		},
	}

	for i, tt := range tests {
		output := trimDistribution(tt.input)

		if tt.expected != output {
			t.Errorf("[%d] expected %q got %q", i, tt.expected, output)
		}
	}
}
