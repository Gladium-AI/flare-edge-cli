package logs

import "testing"

func TestNormalizeWorkerTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "worker name",
			input: "test-project",
			want:  "test-project",
		},
		{
			name:  "workers dev url",
			input: "https://test-project.paoloanzani.workers.dev",
			want:  "test-project",
		},
		{
			name:  "custom domain root url",
			input: "https://example.com",
			want:  "example.com",
		},
		{
			name:  "custom domain path url",
			input: "https://example.com/api/logs",
			want:  "example.com/api/logs",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeWorkerTarget(test.input); got != test.want {
				t.Fatalf("normalizeWorkerTarget(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
