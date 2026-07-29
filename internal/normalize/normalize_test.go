package normalize

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "strict empty",
			profile: ProfileStrict,
			input:   "",
			want:    "",
		},
		{
			name:    "strict line endings",
			profile: ProfileStrict,
			input:   "one\r\ntwo\rthree",
			want:    "one\ntwo\nthree",
		},
		{
			name:    "strict trims line endings",
			profile: ProfileStrict,
			input:   "\n\nabc  \nxyz   \n\n",
			want:    "abc\nxyz",
		},
		{
			name:    "strict keeps punctuation",
			profile: ProfileStrict,
			input:   "Привет, Мир!",
			want:    "Привет, Мир!",
		},
		{
			name:    "plain empty",
			profile: ProfilePlainTextRU,
			input:   "",
			want:    "",
		},
		{
			name:    "plain lowercase",
			profile: ProfilePlainTextRU,
			input:   "ПрИвЕт",
			want:    "привет",
		},
		{
			name:    "plain yo",
			profile: ProfilePlainTextRU,
			input:   "Ёжик ёлка",
			want:    "ежик елка",
		},
		{
			name:    "plain punctuation",
			profile: ProfilePlainTextRU,
			input:   "Привет, мир!",
			want:    "привет мир",
		},
		{
			name:    "plain newlines",
			profile: ProfilePlainTextRU,
			input:   "один\r\nдва\rтри\nчетыре",
			want:    "один два три четыре",
		},
		{
			name:    "plain spaces",
			profile: ProfilePlainTextRU,
			input:   "  один   два\t\tтри  ",
			want:    "один два три",
		},
		{
			name:    "mixed languages",
			profile: ProfilePlainTextRU,
			input:   "Hello, Мир!",
			want:    "hello мир",
		},
		{
			name:    "unknown profile",
			profile: Profile("unknown"),
			input:   "text",
			wantErr: true,
		},
		{
			name:    "unicode NFC",
			profile: ProfileStrict,
			input:   "е\u0308",
			want:    "ё",
		},
		{
			name:    "plain removes symbols",
			profile: ProfilePlainTextRU,
			input:   "Счёт № 12345",
			want:    "счет 12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.input, tt.profile)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Normalize returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
