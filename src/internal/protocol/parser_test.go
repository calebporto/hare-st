package protocol

import "testing"

func TestParseRawRequest(t *testing.T) {
	p := &Parser{}

	t.Run("empty request", func(t *testing.T) {
		_, err := p.ParseRawRequest("")
		if err == nil {
			t.Fatal("expected error for empty request, got nil")
		}
	})

	t.Run("valid get command", func(t *testing.T) {
		result, err := p.ParseRawRequest("get\nminhachave")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := ParsedRawRequest{Verb: "get", Key: "minhachave", DataType: "", Data: ""}
		if result != want {
			t.Fatalf("expected %+v, got %+v", want, result)
		}
	})

	t.Run("valid put command", func(t *testing.T) {
		result, err := p.ParseRawRequest("put\nminhachave\nint\n42")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := ParsedRawRequest{Verb: "put", Key: "minhachave", DataType: "int", Data: "42"}
		if result != want {
			t.Fatalf("expected %+v, got %+v", want, result)
		}
	})

	t.Run("invalid number of commands", func(t *testing.T) {
		_, err := p.ParseRawRequest("put\nminhachave\nint")
		if err == nil {
			t.Fatal("expected error for 3 commands, got nil")
		}
	})

	t.Run("too few commands", func(t *testing.T) {
		_, err := p.ParseRawRequest("get")
		if err == nil {
			t.Fatal("expected error for single command, got nil")
		}
	})

	t.Run("too many commands", func(t *testing.T) {
		_, err := p.ParseRawRequest("put\nchave\nint\n42\nextra")
		if err == nil {
			t.Fatal("expected error for too many commands, got nil")
		}
	})
}

func TestIsValidVerb(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		verb   string
		wantOk bool
	}{
		{"get", true},
		{"put", true},
		{"clean", true},
		{"GET", false},
		{"delete", false},
		{"", false},
	}

	for _, c := range cases {
		ok, err := p.isValidVerb(c.verb)
		if ok != c.wantOk {
			t.Errorf("isValidVerb(%q) = %v, want %v (err: %v)", c.verb, ok, c.wantOk, err)
		}
		if ok && err != nil {
			t.Errorf("isValidVerb(%q) returned ok=true but err=%v", c.verb, err)
		}
		if !ok && err == nil {
			t.Errorf("isValidVerb(%q) returned ok=false but err=nil", c.verb)
		}
	}
}

func TestIsValidKey(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		name   string
		key    string
		wantOk bool
	}{
		{"valid alphanumeric", "chave123", true},
		{"valid with underscore", "minha_chave", true},
		{"valid with hyphen", "minha-chave", true},
		{"valid with colon", "user:1000", true},
		{"valid with cedilla", "chaveç", true},
		{"invalid accented char", "chavé", false},
		{"invalid space", "minha chave", false},
		{"invalid empty", "", false},
		{"too long key", generateString("a", 257), false},
		{"max length key", generateString("a", 256), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := p.isValidKey(c.key)
			if ok != c.wantOk {
				t.Errorf("isValidKey(%q) = %v, want %v (err: %v)", c.key, ok, c.wantOk, err)
			}
		})
	}
}

func TestIsValidType(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		dataType string
		wantOk   bool
	}{
		{"char", true},
		{"string", true},
		{"int", true},
		{"float", true},
		{"boolean", true},
		{"INT", false},
		{"unknown", false},
		{"", false},
	}

	for _, c := range cases {
		ok, err := p.isValidType(c.dataType)
		if ok != c.wantOk {
			t.Errorf("isValidType(%q) = %v, want %v (err: %v)", c.dataType, ok, c.wantOk, err)
		}
	}
}

func TestIsValidData(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		name     string
		data     string
		dataType string
		wantOk   bool
	}{
		{"valid int", "42", "int", true},
		{"valid negative int", "-42", "int", true},
		{"invalid int", "abc", "int", false},
		{"int overflow", "99999999999999999999", "int", false},

		{"valid float", "3.14", "float", true},
		{"invalid float", "abc", "float", false},

		{"valid boolean true", "true", "boolean", true},
		{"valid boolean false", "false", "boolean", true},
		{"invalid boolean", "True", "boolean", false},
		{"invalid boolean yes", "yes", "boolean", false},

		{"valid char", "a", "char", true},
		{"invalid char multiple", "ab", "char", false},
		{"invalid char empty", "", "char", false},

		{"valid string", "hello world", "string", true},
		{"too long string", generateString("a", 1025), "string", false},

		{"unknown data type", "value", "unknown", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := p.isValidData(c.data, c.dataType)
			if ok != c.wantOk {
				t.Errorf("isValidData(%q, %q) = %v, want %v (err: %v)", c.data, c.dataType, ok, c.wantOk, err)
			}
		})
	}
}

func generateString(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}
