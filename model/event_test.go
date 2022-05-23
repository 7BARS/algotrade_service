package model

import (
	"fmt"
	"regexp"
	"testing"
)

func TestParseEvents1(t *testing.T) {
	// ([a-z]*\.[a-z_]*)\(([^)]+)\)([ |])(>|<|==|!=|>=|<=)\ ([0-9]*)|([ |])(and|or)
	rg, err := regexp.Compile(`(([a-z]*)\.([a-z_]*))\(([^)]+)\)[\s|\S](>|<|==|!=|>=|<=)[\s|\S]([0-9]*)[\s|\S](and|)`)
	if err != nil {
		t.Fatal(err)
	}
	str := "any.rsi(14, close, 1440) > 70 and any.rsi(4, close, 5) < 30 or any.rsi(3, close, 1) < 20 and aapl.change_percent(close) > 2"
	// res := rg.FindStringSubmatch(str)
	// res = rg.FindAllString(str, 10)
	res2 := rg.FindAllStringSubmatch(str, 10)
	for i := 0; i < len(res2); i++ {
		symbol := res2[i][2]
		typeOfIndicator := res2[i][3]
		operator := res2[i][5]
		value := res2[i][6]
		union := res2[i][7]
		fmt.Println(symbol)
		fmt.Println(typeOfIndicator)
		fmt.Println(operator)
		fmt.Println(value)
		fmt.Println(union)
	}

	// fmt.Println(res2)
	// eventsForSymbol := make(map[string]IndicatorEvent)
	// indexLast := strings.IndexByte(str, '.')
	// fmt.Println(indexLast)
	// symbol := str[:indexLast]
	// eventsForSymbol[symbol] = IndicatorEvent{}
	// fmt.Println(symbol)

	// indexFirst := indexLast
	// indexLast = strings.IndexByte(str, '(')
	// typeOfIndicator := str[indexFirst+1 : indexLast]
	// println(typeOfIndicator)

	// indexFirst = indexLast
	// indexLast = strings.IndexByte(str, ')')
	// args := str[indexFirst+1 : indexLast]
	// // args = strings.Trim(args, " ")
	// args = strings.ReplaceAll(args, " ", "")
	// println(args)

	// list := strings.Split(args, ",")
	// fmt.Println(list)

	// if 1 >= 0 {

	// }
}
