package protocol

import (
	"testing"

	"github.com/calebportodev/golang/hare-st/src/internal/storage"
)

func TestParseRawRequest(t *testing.T) {
	p := &Parser{}

	t.Run("empty request", func(t *testing.T) {
		_, err := p.ParseRawRequest("")
		if err == nil {
			t.Fatal("expected error for empty request, got nil")
		}
	})

	t.Run("valid get command", func(t *testing.T) {
		result, err := p.ParseRawRequest("get minhachave")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := ParsedRawRequest{Verb: "get", Key: "minhachave", DataType: "", Data: ""}
		if result != want {
			t.Fatalf("expected %+v, got %+v", want, result)
		}
	})

	t.Run("valid put command", func(t *testing.T) {
		result, err := p.ParseRawRequest("put minhachave int 42")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := ParsedRawRequest{Verb: "put", Key: "minhachave", DataType: "int", Data: "42"}
		if result != want {
			t.Fatalf("expected %+v, got %+v", want, result)
		}
	})

	t.Run("invalid number of commands", func(t *testing.T) {
		_, err := p.ParseRawRequest("put minhachave int")
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
		_, err := p.ParseRawRequest("put chave int 42 extra")
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

func TestCleanRawRequest(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"trailing newline", "get chave\n", "get chave"},
		{"trailing carriage return", "get chave\r\n", "get chave"},
		{"surrounding spaces", "  get chave  \n", "get chave"},
		{"no trailing characters", "get chave", "get chave"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := p.CleanRawRequest(c.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != c.want {
				t.Errorf("CleanRawRequest(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

func TestParsePutResponse(t *testing.T) {
	p := &Parser{}
	got, err := p.ParsePutResponse(ParsedRawRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "status success" {
		t.Errorf("ParsePutResponse() = %q, want %q", got, "status success")
	}
}

func TestParseCleanResponse(t *testing.T) {
	p := &Parser{}
	got, err := p.ParseCleanResponse(ParsedRawRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "status success" {
		t.Errorf("ParseCleanResponse() = %q, want %q", got, "status success")
	}
}

func TestParseGetResponse(t *testing.T) {
	p := &Parser{}

	t.Run("key not found", func(t *testing.T) {
		got, err := p.ParseGetResponse(storage.StorageValue{}, false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "data: nil nil" {
			t.Errorf("ParseGetResponse() = %q, want %q", got, "data: nil nil")
		}
	})

	t.Run("key found", func(t *testing.T) {
		got, err := p.ParseGetResponse(storage.StorageValue{DataType: "int", Data: "42"}, true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := "data: int 42"
		if got != want {
			t.Errorf("ParseGetResponse() = %q, want %q", got, want)
		}
	})
}

func TestTrimFloatTrailingZeros(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		input string
		want  string
	}{
		{"1.230000", "1.23"},
		{"5.000000", "5"},
		{"42", "42"},
		{"-3.100000", "-3.1"},
		{"0.000000", "0"},
	}

	for _, c := range cases {
		got := p.trimFloatTrailingZeros(c.input)
		if got != c.want {
			t.Errorf("trimFloatTrailingZeros(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestIsValidRequestEdgeCases(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		name     string
		commands []string
		wantOk   bool
	}{
		{"too few", []string{"get"}, false},
		{"invalid number of commands", []string{"put", "key", "int"}, false},
		{"too many", []string{"put", "key", "int", "1", "extra"}, false},
		{"valid two commands", []string{"get", "key"}, true},
		{"invalid verb", []string{"unknown", "key"}, false},
		{"invalid key", []string{"get", "invalid key"}, false},
		{"invalid type", []string{"put", "key", "unknown", "1"}, false},
		{"invalid data", []string{"put", "key", "int", "abc"}, false},
		{"valid four commands", []string{"put", "key", "int", "1"}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := p.isValidRequest(c.commands)
			if ok != c.wantOk {
				t.Errorf("isValidRequest(%v) = %v, want %v (err: %v)", c.commands, ok, c.wantOk, err)
			}
		})
	}
}

func TestIsValidDataLengthRules(t *testing.T) {
	p := &Parser{}

	cases := []struct {
		name     string
		data     string
		dataType string
		wantOk   bool
	}{
		{"int max digits", "9223372036854775807", "int", true},
		{"int too many digits", generateString("9", 20), "int", false},
		{"negative int max digits", "-9223372036854775808", "int", true},
		{"only sign", "-", "int", false},
		{"too long data", generateString("a", 1025), "string", false},

		{"float max total digits", generateString("9", 9) + "." + generateString("9", 6), "float", true},
		{"float too many decimal digits", "1." + generateString("9", 7), "float", false},
		{"float too many total digits", generateString("9", 10) + "." + generateString("9", 6), "float", false},
		{"float missing decimal separator", "123", "float", false},
		{"float multiple dots", "1.2.3", "float", false},
		{"float missing integer part", ".5", "float", false},
		{"float missing decimal part", "5.", "float", false},

		{"char two-byte rune", "ç", "char", true},
		{"char three-byte rune", "€", "char", false},
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
