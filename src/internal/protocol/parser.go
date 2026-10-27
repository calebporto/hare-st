package protocol

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/calebportodev/golang/hare-st/src/internal/storage"
)

var Delimiter string = " "
var AcceptedVerbs []string = []string{"put", "get", "clean"}
var ValidKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_:çÇ-]+$`)
var AcceptedTypes []string = []string{"char", "string", "int", "float", "boolean"}

type Parser struct{}

type ParsedRawRequest struct {
	Verb     string
	Key      string
	DataType string
	Data     string
}

func (p *Parser) ParseRawRequest(rawRequest string) (ParsedRawRequest, error) {
	if rawRequest == "" {
		return ParsedRawRequest{}, errors.New("Empty request")
	}

	commands := strings.Split(rawRequest, Delimiter)

	_, err := p.isValidRequest(commands)
	if err != nil {
		return ParsedRawRequest{}, err
	}

	verb, key := commands[0], commands[1]
	dataType, data := "", ""

	if len(commands) == 4 {
		dataType, data = commands[2], commands[3]
	}
	if dataType == "float" {
		data = p.trimFloatTrailingZeros(data)
	}
	
	return ParsedRawRequest{
		Verb:     verb,
		Key:      key,
		DataType: dataType,
		Data:     data,
	}, nil
}

func (p *Parser) CleanRawRequest(rawRequest string) (string, error) {
	rawRequest = strings.TrimRight(rawRequest, "\n")
	rawRequest = strings.TrimRight(rawRequest, "\r")
	rawRequest = strings.TrimSpace(rawRequest)
	return rawRequest, nil
}

func (p *Parser) ParsePutResponse(parsedRequest ParsedRawRequest) (string, error) {
	return "status success", nil
}

func (p *Parser) ParseCleanResponse(parsedRequest ParsedRawRequest) (string, error) {
	return "status success", nil
}

func (p *Parser) ParseGetResponse(data storage.StorageValue, dataOk bool) (string, error) {
	if !dataOk {
		return "data: nil nil", nil
	}

	return fmt.Sprintf("data: %s %s", data.DataType, data.Data), nil
}

func (p *Parser) isValidRequest(commands []string) (ok bool, err error) {
	if len(commands) <= 1 {
		return false, errors.New("Too few commands")
	}
	if len(commands) > 4 {
		return false, errors.New("Too many commands")
	}
	if len(commands) == 3 {
		return false, errors.New("Invalid number of commands")
	}

	verb, key := commands[0], commands[1]
	isValidVerb, verbErr := p.isValidVerb(verb)
	if !isValidVerb {
		return false, verbErr
	}
	isValidKey, keyErr := p.isValidKey(key)
	if !isValidKey {
		return false, keyErr
	}

	if len(commands) == 2 {
		return true, nil
	}

	dataType, data := commands[2], commands[3]
	isValidType, typeErr := p.isValidType(dataType)
	if !isValidType {
		return false, typeErr
	}
	isValidData, dataErr := p.isValidData(data, dataType)
	if !isValidData {
		return false, dataErr
	}

	return true, nil
}

func (p *Parser) isValidVerb(verb string) (ok bool, err error) {
	if !slices.Contains(AcceptedVerbs, verb) {
		return false, errors.New("Invalid verb")
	}
	return true, nil
}

func (p *Parser) isValidKey(key string) (ok bool, err error) {
	if !ValidKeyPattern.MatchString(key) {
		return false, errors.New("Invalid key")
	}
	// Check key length
	if utf8.RuneCountInString(key) > 256 {
		return false, errors.New("Too long key")
	}
	return true, nil
}

func (p *Parser) isValidType(dataType string) (ok bool, err error) {
	if !slices.Contains(AcceptedTypes, dataType) {
		return false, errors.New("Invalid data type")
	}
	return true, nil
}

func (p *Parser) isValidData(data string, dataType string) (ok bool, err error) {
	// Check data length
	if utf8.RuneCountInString(data) > 1024 {
		return false, errors.New("Too long data")
	}
	if dataType == "int" {
		intDigits := strings.TrimPrefix(data, "-")
		if len(intDigits) == 0 || len(intDigits) > 19 {
			return false, fmt.Errorf("Invalid int length. Max length is 19 digits.")
		}
		_, intErr := strconv.ParseInt(data, 10, 64)
		if intErr != nil {
			return false, fmt.Errorf("Invalid int")
		}
		return true, nil
	}
	if dataType == "float" {
		_, floatErr := strconv.ParseFloat(data, 64)
		if floatErr != nil {
			return false, fmt.Errorf("Invalid float")
		}

		floatDigits := strings.TrimPrefix(data, "-")
		floatParts := strings.Split(floatDigits, ".")
		if len(floatParts) != 2 {
			return false, fmt.Errorf("Invalid float format. Expected a single '.' as decimal separator.")
		}
		integerPart, decimalPart := floatParts[0], floatParts[1]
		if len(integerPart) == 0 || len(decimalPart) == 0 {
			return false, fmt.Errorf("Invalid float format. Missing integer or decimal part.")
		}
		if len(decimalPart) > 6 {
			return false, fmt.Errorf("Invalid float. Max length is 6 decimal digits.")
		}
		if len(integerPart)+len(decimalPart) > 15 {
			return false, fmt.Errorf("Invalid float length. Max length is 15 digits.")
		}
		return true, nil
	}
	if dataType == "boolean" {
		if data != "true" && data != "false" {
			return false, fmt.Errorf("Invalid boolean")
		}
		return true, nil
	}
	if dataType == "char" {
		if utf8.RuneCountInString(data) != 1 {
			return false, fmt.Errorf("Invalid char")
		}
		if len(data) > 2 {
			return false, fmt.Errorf("Invalid char. Max length is 2 bytes.")
		}
		return true, nil
	}
	if dataType == "string" {
		if utf8.RuneCountInString(data) > 1024 {
			return false, fmt.Errorf("String too long. Max length is 1024 runes.")
		}
		return true, nil
	}
	return false, fmt.Errorf("Invalid data type: %v", dataType)
}

func (p *Parser) trimFloatTrailingZeros(data string) string {
	if !strings.Contains(data, ".") {
		return data
	}
	trimmed := strings.TrimRight(data, "0")
	trimmed = strings.TrimSuffix(trimmed, ".")
	return trimmed
}