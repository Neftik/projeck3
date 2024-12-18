package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "unicode"
)

func Calc(expression string) (float64, error) {
    precedence := map[rune]int{
        '+': 1,
        '-': 1,
        '*': 2,
        '/': 2,
    }
    var numbers []float64
    var operators []rune
    var getNumber = func(c rune) (float64, error) {
        return strconv.ParseFloat(string(c), 64)
    }
    
    var applyOperator = func() error {
        if len(numbers) < 2 {
            return errors.New("недостаточно операндов")
        }
        b := numbers[len(numbers)-1]
        numbers = numbers[:len(numbers)-1]
        a := numbers[len(numbers)-1]
        numbers = numbers[:len(numbers)-1]
        op := operators[len(operators)-1]
        operators = operators[:len(operators)-1]
        var result float64
        switch op {
        case '+':
            result = a + b
        case '-':
            result = a - b
        case '*':
            result = a * b
        case '/':
            if b == 0 {
                return errors.New("деление на ноль")
            }
            result = a / b
        default:
            return errors.New("неизвестный оператор")
        }
        numbers = append(numbers, result)
        return nil
    }
    var processOperator = func(op rune) error {
        for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[op] {
            if err := applyOperator(); err != nil {
                return err
            }
        }
        operators = append(operators, op)
        return nil
    }

    for i, char := range expression {
        if unicode.IsDigit(char) {
            num, err := getNumber(char)
            if err != nil {
                return 0, err
            }
            numbers = append(numbers, num)
        } else if char == '(' {
            operators = append(operators, char)
        } else if char == ')' {
            for len(operators) > 0 && operators[len(operators)-1] != '(' {
                if err := applyOperator(); err != nil {
                    return 0, err
                }
            }
            if len(operators) == 0 || operators[len(operators)-1] != '(' {
                return 0, errors.New("несоответствие скобок")
            }
            operators = operators[:len(operators)-1]
        } else if char == '+' || char == '-' || char == '*' || char == '/' {
            if err := processOperator(char); err != nil {
                return 0, err
            }
        } else {
            return 0, fmt.Errorf("недопустимый символ в выражении на позиции %d: %c", i, char)
        }
    }
    for len(operators) > 0 {
        if operators[len(operators)-1] == '(' {
            return 0, errors.New("несоответствие скобок")
        }
        if err := applyOperator(); err != nil {
            return 0, err
        }
    }
    if len(numbers) != 1 {
        return 0, errors.New("некорректное выражение")
    }
    return numbers[0], nil
}

type RequestBody struct {
    Expression string `json:"expression"`
}

type ResponseBody struct {
    Result string `json:"result,omitempty"`
    Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var reqBody RequestBody
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusUnprocessableEntity)
        return
    }

    result, err := Calc(reqBody.Expression)
    var resBody ResponseBody

    if err != nil {
        if errors.Is(err, strconv.ErrSyntax) || errors.Is(err, strconv.ErrRange) {
            resBody.Error = "Expression is not valid"
            w.WriteHeader(http.StatusUnprocessableEntity)
        } else {
            resBody.Error = "Internal server error"
            w.WriteHeader(http.StatusInternalServerError)
        }
    } else {
        resBody.Result = fmt.Sprintf("%f", result)
        w.WriteHeader(http.StatusOK)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resBody)
}

func main() {
    http.HandleFunc("/api/v1/calculate", calculateHandler)
    fmt.Println("Server running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Server failed to start:", err)
    }
}


//http://localhost:8080/api/v1/calculate