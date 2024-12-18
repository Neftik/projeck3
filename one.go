package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if !regexp.MustCompile(`^[\d+\-*/().]+$`).MatchString(expression) {
		return 0, errors.New("Знаки невалидны!!!")
	}

	tokens, err := tokenize(expression)
	if err != nil {
		return 0, fmt.Errorf("tokenization error: %w", err)
	}

	output, err := shuntingYard(tokens)
	if err != nil {
		return 0, fmt.Errorf("shunting yard error: %w", err)
	}

	result, err := evaluateRPN(output)
	if err != nil {
		return 0, fmt.Errorf("evaluation error: %w", err)
	}

	return result, nil
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	var current strings.Builder

	for _, char := range expression {
		if strings.ContainsRune("0123456789.", char) {
			current.WriteRune(char)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			if strings.ContainsRune("+-*/()", char) {
				tokens = append(tokens, string(char))
			} else {
				return nil, fmt.Errorf("invalid character '%c'", char)
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

func shuntingYard(tokens []string) ([]string, error) {
	precedence := map[string]int{
		"+": 1, "-": 1, "*": 2, "/": 2,
	}

	var output []string
	var operators []string

	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
		case token == "(":
			operators = append(operators, token)
		case token == ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 || operators[len(operators)-1] != "(" {
				return nil, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		case precedence[token] > 0:
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		default:
			return nil, fmt.Errorf("invalid token '%s'", token)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number '%s': %w", token, err)
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("insufficient operands")
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				result = a / b
			default:
				return 0, fmt.Errorf("invalid operator '%s'", token)
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

type RequestBody struct {
	Expression string `json:"expression"`
}

type ResponseBody struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")

	if r.Method != http.MethodPost {
		log.Println("Invalid method")
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if reqBody.Expression == "" {
		log.Println("Empty expression provided")
		http.Error(w, `{"error": "Expression cannot be empty"}`, http.StatusBadRequest)
		return
	}

	result, err := Calc(reqBody.Expression)
	resBody := ResponseBody{}

	if err != nil {
		log.Println("Calculation error:", err)
		resBody.Error = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		if result == math.Trunc(result) {
			resBody.Result = fmt.Sprintf("%g", result)
		} else {
			resBody.Result = fmt.Sprintf("%f", result)
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		log.Println("Failed to encode response:", err)
	}
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}

//http://localhost:8080/api/v1/calculate
