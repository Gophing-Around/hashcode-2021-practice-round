package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	files := []string{
		"a_example",
		"b_little_bit_of_everything.in",
		"c_many_ingredients.in",
		"d_many_pizzas.in",
		"e_many_teams.in",
	}

	desired := 0

	inputSet := readFile(fmt.Sprintf("./dataset/%s", files[desired]))

	parts := strings.Split(inputSet, "\n")
	config := getConfig(parts[0])
	pizzas := make([]Pizza, config.pizzaNumber)
	for i := 0; i < config.pizzaNumber; i++ {
		pizzas[i] = getPizzaIngredients(parts[i+1], i)
	}

	orders := firstPlaceOrder(config, pizzas)
	result := fmt.Sprintf("%d\n", len(orders))
	for _, order := range orders {
		result += fmt.Sprintf("%d %s\n", len(order.pizzas), strings.Join(order.pizzas, " "))
	}
	fmt.Printf("%s", result)

	ioutil.WriteFile(fmt.Sprintf("./result/%s", files[desired]), []byte(result), 0644)
}

func firstPlaceOrder(config Config, pizzas []Pizza) []OrderDelivery {
	orders := []OrderDelivery{}

	pizzaCounter := 0

	for i := 0; i < config.nTeamOfFour; i++ {
		if pizzaCounter+4 > len(pizzas) {
			return orders
		}

		order := OrderDelivery{
			pizzas: []string{
				pizzas[pizzaCounter].pizzaID,
				pizzas[pizzaCounter+1].pizzaID,
				pizzas[pizzaCounter+2].pizzaID,
				pizzas[pizzaCounter+3].pizzaID,
			},
		}
		pizzaCounter += 4
		orders = append(orders, order)
	}

	for i := 0; i < config.nTeamOfThree; i++ {
		if pizzaCounter+3 > len(pizzas) {
			return orders
		}

		order := OrderDelivery{
			pizzas: []string{
				pizzas[pizzaCounter].pizzaID,
				pizzas[pizzaCounter+1].pizzaID,
				pizzas[pizzaCounter+2].pizzaID,
			},
		}
		pizzaCounter += 3
		orders = append(orders, order)
	}

	for i := 0; i < config.nTeamOfTwo; i++ {
		if pizzaCounter+2 > len(pizzas) {
			return orders
		}

		order := OrderDelivery{
			pizzas: []string{
				pizzas[pizzaCounter].pizzaID,
				pizzas[pizzaCounter+1].pizzaID,
			},
		}
		pizzaCounter += 2
		orders = append(orders, order)
	}

	return orders
}

func readFile(source string) string {
	in, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	return string(in)
}

type OrderDelivery struct {
	pizzas []string
}

type Config struct {
	pizzaNumber  int
	nTeamOfTwo   int
	nTeamOfThree int
	nTeamOfFour  int
}

type Pizza struct {
	nIngr       int
	ingredients []string
	ingrMap     map[string]bool
	pizzaID     string
}

func getPizzaIngredients(pizzaLine string, pizzaID int) Pizza {
	parts := strings.Split(pizzaLine, " ")

	nIngr := toint(parts[0])

	ingredients := make([]string, nIngr)
	ingrMap := make(map[string]bool, nIngr)
	for i := 0; i < nIngr; i++ {
		ingr := parts[i+1]
		ingredients[i] = ingr
		ingrMap[ingr] = true
	}

	return Pizza{
		nIngr:       nIngr,
		ingredients: ingredients,
		ingrMap:     ingrMap,
		pizzaID:     fmt.Sprintf("%d", pizzaID),
	}
}

func getConfig(configLine string) Config {
	parts := strings.Split(configLine, " ")
	return Config{
		pizzaNumber:  toint(parts[0]),
		nTeamOfTwo:   toint(parts[1]),
		nTeamOfThree: toint(parts[2]),
		nTeamOfFour:  toint(parts[3]),
	}
}

func toint(num string) int {
	res, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	return res
}
