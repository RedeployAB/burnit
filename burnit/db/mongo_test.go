package db

import (
	"testing"
)

func TestToURI(t *testing.T) {
	var tests = []struct {
		name  string
		input *MongoClientOptions
		want  string
	}{
		{
			input: &MongoClientOptions{
				Address: "localhost",
			},
			want: "mongodb://localhost",
		},
		{
			input: &MongoClientOptions{
				Address: "localhost:27017",
			},
			want: "mongodb://localhost:27017",
		},
		{
			input: &MongoClientOptions{
				Address:  "localhost:27017",
				Username: "user",
			},
			want: "mongodb://user@localhost:27017",
		},
		{
			input: &MongoClientOptions{
				Address:  "localhost:27017",
				Username: "user",
				Password: "1234",
			},
			want: "mongodb://user:1234@localhost:27017",
		},
		{
			input: &MongoClientOptions{
				Address:  "localhost:27017",
				Username: "user",
				Password: "1234",
				Database: "db",
			},
			want: "mongodb://user:1234@localhost:27017/db",
		},
		{
			input: &MongoClientOptions{
				Address:  "localhost:27017",
				Username: "user",
				Password: "1234",
				Database: "db",
				SSL:      true,
			},
			want: "mongodb://user:1234@localhost:27017/db?ssl=true",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := toURI(test.input)
			if got != test.want {
				t.Errorf("incorrect value, want: %s, got: %s\n", test.want, got)
			}
		})

	}
}
