package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    command
		wantErr string
	}{
		{
			name:    "empty args",
			input:   []string{},
			want:    nil,
			wantErr: "no command was specified",
		},
		{
			name:    "unknown command",
			input:   []string{"idk", "1"},
			want:    nil,
			wantErr: "unknown command",
		},
		// add
		{
			name:    "correct add command with multiple words",
			input:   []string{"add", "Watch", "movie"},
			want:    addCommand{text: "Watch movie"},
			wantErr: "",
		},
		{
			name:    "empty text add command",
			input:   []string{"add", ""},
			want:    nil,
			wantErr: "todo text cannot be empty",
		},
		{
			name:    "space text add command",
			input:   []string{"add", "  "},
			want:    nil,
			wantErr: "todo text cannot be empty",
		},
		// list
		{
			name:    "correct list all",
			input:   []string{"list"},
			want:    listAllCommand{},
			wantErr: "",
		},
		{
			name:  "correct list one",
			input: []string{"list", "2"},
			want: listOneCommand{
				id: 2,
			},
			wantErr: "",
		},
		{
			name:    "list propagates id parsing error",
			input:   []string{"list", "-2"},
			want:    nil,
			wantErr: "list: id must not contain a sign",
		},
		// done
		{
			name:  "correct done command",
			input: []string{"done", "3"},
			want: doneCommand{
				id: 3,
			},
			wantErr: "",
		},
		{
			name:    "done propagates id parsing error",
			input:   []string{"done", "-3"},
			want:    nil,
			wantErr: "done: id must not contain a sign",
		},
		// undone
		{
			name:  "correct undone command",
			input: []string{"undone", "4"},
			want: undoneCommand{
				id: 4,
			},
			wantErr: "",
		},
		{
			name:    "undone propagates id parsing error",
			input:   []string{"undone", "-4"},
			want:    nil,
			wantErr: "undone: id must not contain a sign",
		},
		// delete
		{
			name:  "correct delete command",
			input: []string{"delete", "5"},
			want: deleteCommand{
				id: 5,
			},
			wantErr: "",
		},
		{
			name:    "delete propagates id parsing error",
			input:   []string{"delete", "-5"},
			want:    nil,
			wantErr: "delete: id must not contain a sign",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.input)

			checkErr(t, err, tt.wantErr)

			if tt.want != got {
				t.Errorf("несовпадение результатов:\nожидали: %+v\nполучили: %+v", tt.want, got)
			}
		})
	}
}

func TestParseIDArg(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    int
		wantErr string
	}{
		{
			name:    "correct order",
			input:   []string{"done", "4"},
			want:    4,
			wantErr: "",
		},
		{
			name:    "no id",
			input:   []string{"done"},
			want:    0,
			wantErr: "expected exactly 1 integer argument for id, but got 0",
		},
		{
			name:    "too many arguments",
			input:   []string{"done", "1", "2"},
			want:    0,
			wantErr: "expected exactly 1 integer argument for id, but got 2",
		},
		{
			name:    "id is not int",
			input:   []string{"done", "notint"},
			want:    0,
			wantErr: "id is not correct",
		},
		{
			name:    "no command at all",
			input:   []string{},
			want:    0,
			wantErr: "no command was provided, args: []",
		},
		{
			name:    "incorrect id with letter",
			input:   []string{"done", "4c"},
			want:    0,
			wantErr: "id is not correct",
		},
		{
			name:    "incorrect id with minus sign",
			input:   []string{"done", "-42"},
			want:    0,
			wantErr: "id must not contain a sign",
		},
		{
			name:    "incorrect id with plus sign",
			input:   []string{"done", "+42"},
			want:    0,
			wantErr: "id must not contain a sign",
		},
		{
			name:    "too big id",
			input:   []string{"delete", "2348948294823948948239428492849248928394283943289484"},
			want:    0,
			wantErr: "id is not correct",
		},
		{
			name:    "0 id correct too",
			input:   []string{"delete", "0"},
			want:    0,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIDArg(tt.input)

			checkErr(t, err, tt.wantErr)

			if tt.want != got {
				t.Errorf("несовпадение результатов:\nожидали: %+v\nполучили: %+v", tt.want, got)
			}
		})
	}
}
