package redis

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTxResult_AllBytes(t *testing.T) {
	result := TxResult{
		b: [][]byte{
			testBytes1,
			testBytes2,
			testBytes3,
		},
	}

	want := [][]byte{testBytes1, testBytes2, testBytes3}
	got := result.AllBytes()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AllBytes() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestTxResult_FirstBytes(t *testing.T) {
	var tests = []struct {
		name  string
		input [][]byte
		want  []byte
	}{
		{
			name:  "multiple data",
			input: [][]byte{testBytes1, testBytes2, testBytes3},
			want:  testBytes1,
		},
		{
			name:  "one",
			input: [][]byte{testBytes1},
			want:  testBytes1,
		},
		{
			name:  "empty",
			input: [][]byte{},
			want:  nil,
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				b: test.input,
			}

			got := result.FirstBytes()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("FirstBytes() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTxResult_LastBytes(t *testing.T) {
	var tests = []struct {
		name  string
		input [][]byte
		want  []byte
	}{
		{
			name:  "multiple data",
			input: [][]byte{testBytes1, testBytes2, testBytes3},
			want:  testBytes3,
		},
		{
			name:  "one",
			input: [][]byte{testBytes1},
			want:  testBytes1,
		},
		{
			name:  "empty",
			input: [][]byte{},
			want:  nil,
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				b: test.input,
			}

			got := result.LastBytes()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("LastBytes() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTxResult_IndexBytes(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			b     [][]byte
			index int
		}
		want []byte
	}{
		{
			name: "index 0",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 0,
			},
			want: testBytes1,
		},
		{
			name: "index 1",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 1,
			},
			want: testBytes2,
		},
		{
			name: "index 2",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 2,
			},
			want: testBytes3,
		},
		{
			name: "index 3",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 3,
			},
			want: nil,
		},
		{
			name: "index 4",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 4,
			},
			want: nil,
		},
		{
			name: "index -1",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{testBytes1, testBytes2, testBytes3},
				index: 4,
			},
			want: nil,
		},
		{
			name: "empty",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     [][]byte{},
				index: 1,
			},
			want: nil,
		},
		{
			name: "nil",
			input: struct {
				b     [][]byte
				index int
			}{
				b:     nil,
				index: 1,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				b: test.input.b,
			}

			got := result.IndexBytes(test.input.index)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("IndexBytes() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTxResult_AllMaps(t *testing.T) {
	result := TxResult{
		m: []map[string]string{
			testMap1,
			testMap2,
			testMap3,
		},
	}

	want := []map[string]string{testMap1, testMap2, testMap3}
	got := result.AllMaps()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AllMaps() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestTxResult_FirstMap(t *testing.T) {
	var tests = []struct {
		name  string
		input []map[string]string
		want  map[string]string
	}{
		{
			name:  "multiple data",
			input: []map[string]string{testMap1, testMap2, testMap3},
			want:  testMap1,
		},
		{
			name:  "one",
			input: []map[string]string{testMap1},
			want:  testMap1,
		},
		{
			name:  "empty",
			input: []map[string]string{},
			want:  nil,
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				m: test.input,
			}

			got := result.FirstMap()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("FirstMap() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTxResult_LastMap(t *testing.T) {
	var tests = []struct {
		name  string
		input []map[string]string
		want  map[string]string
	}{
		{
			name:  "multiple data",
			input: []map[string]string{testMap1, testMap2, testMap3},
			want:  testMap3,
		},
		{
			name:  "one",
			input: []map[string]string{testMap1},
			want:  testMap1,
		},
		{
			name:  "empty",
			input: []map[string]string{},
			want:  nil,
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				m: test.input,
			}

			got := result.LastMap()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("LastMap() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTxResult_IndexMap(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			m     []map[string]string
			index int
		}
		want map[string]string
	}{
		{
			name: "index 0",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 0,
			},
			want: testMap1,
		},
		{
			name: "index 1",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 1,
			},
			want: testMap2,
		},
		{
			name: "index 2",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 2,
			},
			want: testMap3,
		},
		{
			name: "index 3",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 3,
			},
			want: nil,
		},
		{
			name: "index 4",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 4,
			},
			want: nil,
		},
		{
			name: "index -1",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{testMap1, testMap2, testMap3},
				index: 4,
			},
			want: nil,
		},
		{
			name: "empty",
			input: struct {
				m     []map[string]string
				index int
			}{
				m:     []map[string]string{},
				index: 1,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TxResult{
				m: test.input.m,
			}

			got := result.IndexMap(test.input.index)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("IndexMap() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

var (
	testBytes1 = []byte("data1")
	testBytes2 = []byte("data2")
	testBytes3 = []byte("data3")
)

var (
	testMap1 = map[string]string{
		"key1": "value1",
	}
	testMap2 = map[string]string{
		"key2": "value2",
	}
	testMap3 = map[string]string{
		"key3": "value3",
	}
)
