package main

import (
	"fmt"
	"io/ioutil"
	"sort"
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

	for _, fileName := range files {
		fmt.Printf("--------------------------------------------------------")
		fmt.Printf("****************** INPUT: %s\n", fileName)

		inputSet := readFile(fmt.Sprintf("./dataset/%s", fileName))

		parts := strings.Split(inputSet, "\n")
		config := getConfig(parts[0])
		pizzas := make([]Pizza, config.pizzaNumber)
		for i := 0; i < config.pizzaNumber; i++ {
			pizzas[i] = getPizzaIngredients(parts[i+1], i)
		}

		sortPizzas(pizzas)
		orders := firstPlaceOrder(config, pizzas)
		result := fmt.Sprintf("%d\n", len(orders))
		for _, order := range orders {
			result += fmt.Sprintf("%d %s\n", len(order.pizzas), strings.Join(order.pizzas, " "))
		}

		ioutil.WriteFile(fmt.Sprintf("./result/%s", fileName), []byte(result), 0644)
	}
}

func unservedPizzas(pizzas []Pizza) (count int) {
	for _, pizza := range pizzas {
		if !pizza.taken {
			count++
		}
	}
	return
}

func sortPizzas(pizzas []Pizza) {
	sort.Slice(pizzas, func(i, j int) bool {
		return pizzas[i].nIngr > pizzas[j].nIngr
	})
}

func firstPlaceOrder(config Config, pizzas []Pizza) []OrderDelivery {
	orders := []OrderDelivery{}
	pizzaCounter := 0
	teamOfFour := config.nTeamOfFour
	teamOfThree := config.nTeamOfThree
	teamOfTwo := config.nTeamOfTwo

	fmt.Printf("****************** Team requests:\n")
	fmt.Printf("****************** 4 people team: %d\n", config.nTeamOfFour)
	fmt.Printf("****************** 4 people team: %d\n", config.nTeamOfThree)
	fmt.Printf("****************** 4 people team: %d\n", config.nTeamOfTwo)

	for i := 0; i < config.nTeamOfFour; i++ {
		pizzasOrder := make([]Pizza, 4)

		for pizzaCounter < len(pizzas) {
			if !pizzas[pizzaCounter].taken {
				break
			}
			pizzaCounter++
		}
		if pizzaCounter == len(pizzas) {
			break
		}

		pizzasOrder[0] = pizzas[pizzaCounter]
		pizzas[pizzaCounter].taken = true
		pizzaCounter++

		initC := 1
		for j := pizzaCounter; j < len(pizzas); j++ {
			if pizzas[j].taken {
				continue
			}

			var hasMatches bool
			for _, piOrder := range pizzasOrder {
				threshold := float64(piOrder.nIngr) / 30.0
				if float64(howManyIngredientEquals(piOrder.ingrMap, pizzas[j].ingrMap)) > threshold {
					hasMatches = true
				}
			}
			if !hasMatches {
				pizzas[j].taken = true
				pizzasOrder[initC] = pizzas[j]
				initC++
			}
			if initC == 4 {
				break
			}
		}

		order := OrderDelivery{
			pizzas: getPizzasIDs(pizzasOrder),
		}

		if initC == 4 {
			orders = append(orders, order)
			teamOfFour--
		} else {
			for idx, pizza := range pizzas {
				for _, pizzaOrder := range pizzasOrder {
					if pizzaOrder.pizzaID == pizza.pizzaID {
						pizzas[idx].taken = false
					}
				}
			}
		}
	}

	pizzaCounter = 0
	for i := 0; i < config.nTeamOfThree; i++ {
		pizzasOrder := make([]Pizza, 3)

		for pizzaCounter < len(pizzas) {
			if !pizzas[pizzaCounter].taken {
				break
			}
			pizzaCounter++
		}
		if pizzaCounter == len(pizzas) {
			break
		}

		pizzasOrder[0] = pizzas[pizzaCounter]
		pizzas[pizzaCounter].taken = true
		pizzaCounter++

		initC := 1
		for j := pizzaCounter; j < len(pizzas); j++ {
			if pizzas[j].taken {
				continue
			}

			var hasMatches bool
			for _, piOrder := range pizzasOrder {
				threshold := float64(piOrder.nIngr) / 50.0
				if float64(howManyIngredientEquals(piOrder.ingrMap, pizzas[j].ingrMap)) > threshold {
					hasMatches = true
				}
			}
			if !hasMatches {
				pizzas[j].taken = true
				pizzasOrder[initC] = pizzas[j]
				initC++
			}
			if initC == 3 {
				break
			}
		}

		order := OrderDelivery{
			pizzas: getPizzasIDs(pizzasOrder),
		}

		if initC == 3 {
			orders = append(orders, order)
			teamOfThree--
		} else {
			for idx, pizza := range pizzas {
				for _, pizzaOrder := range pizzasOrder {
					if pizzaOrder.pizzaID == pizza.pizzaID {
						pizzas[idx].taken = false
					}
				}
			}
		}
	}

	pizzaCounter = 0
	for i := 0; i < config.nTeamOfTwo; i++ {
		pizzasOrder := make([]Pizza, 2)

		for pizzaCounter < len(pizzas) {
			if !pizzas[pizzaCounter].taken {
				break
			}
			pizzaCounter++
		}
		if pizzaCounter == len(pizzas) {
			break
		}

		pizzasOrder[0] = pizzas[pizzaCounter]
		pizzas[pizzaCounter].taken = true
		pizzaCounter++

		initC := 1
		for j := pizzaCounter; j < len(pizzas); j++ {
			if pizzas[j].taken {
				continue
			}

			var hasMatches bool
			for _, piOrder := range pizzasOrder {
				if howManyIngredientEquals(piOrder.ingrMap, pizzas[j].ingrMap) != 0 {
					hasMatches = true
				}
			}
			if !hasMatches {
				pizzas[j].taken = true
				pizzasOrder[initC] = pizzas[j]
				initC++
			}
			if initC == 2 {
				break
			}
		}

		order := OrderDelivery{
			pizzas: getPizzasIDs(pizzasOrder),
		}

		if initC == 2 {
			orders = append(orders, order)
			teamOfTwo--
		} else {
			for idx, pizza := range pizzas {
				for _, pizzaOrder := range pizzasOrder {
					if pizzaOrder.pizzaID == pizza.pizzaID {
						pizzas[idx].taken = false
					}
				}
			}
		}
	}

	fmt.Println("*********************************************************")
	fmt.Println("FILL ORDERS start")
	fmt.Println("*********************************************************")

	pizzaCounter = 0
	filteredPizzas := []Pizza{}
	for _, pizza := range pizzas {
		if !pizza.taken {
			filteredPizzas = append(filteredPizzas, pizza)
		}
	}

	for i := 0; i < teamOfFour; i++ {
		if pizzaCounter+4 > len(filteredPizzas) {
			break
		}

		filteredPizzas[pizzaCounter].taken = true
		filteredPizzas[pizzaCounter+1].taken = true
		filteredPizzas[pizzaCounter+2].taken = true
		filteredPizzas[pizzaCounter+3].taken = true

		order := OrderDelivery{
			pizzas: []string{
				filteredPizzas[pizzaCounter].pizzaID,
				filteredPizzas[pizzaCounter+1].pizzaID,
				filteredPizzas[pizzaCounter+2].pizzaID,
				filteredPizzas[pizzaCounter+3].pizzaID,
			},
		}
		pizzaCounter += 4
		orders = append(orders, order)
	}

	for i := 0; i < teamOfThree; i++ {
		if pizzaCounter+3 > len(filteredPizzas) {
			break
		}

		filteredPizzas[pizzaCounter].taken = true
		filteredPizzas[pizzaCounter+1].taken = true
		filteredPizzas[pizzaCounter+2].taken = true

		order := OrderDelivery{
			pizzas: []string{
				filteredPizzas[pizzaCounter].pizzaID,
				filteredPizzas[pizzaCounter+1].pizzaID,
				filteredPizzas[pizzaCounter+2].pizzaID,
			},
		}
		pizzaCounter += 3
		orders = append(orders, order)
	}

	for i := 0; i < teamOfTwo; i++ {
		if pizzaCounter+2 > len(filteredPizzas) {
			break
		}

		filteredPizzas[pizzaCounter].taken = true
		filteredPizzas[pizzaCounter+1].taken = true

		order := OrderDelivery{
			pizzas: []string{
				filteredPizzas[pizzaCounter].pizzaID,
				filteredPizzas[pizzaCounter+1].pizzaID,
			},
		}
		pizzaCounter += 2
		orders = append(orders, order)
	}

	fmt.Printf("****************** Unserved pizzas: %d :(\n", unservedPizzas(filteredPizzas))

	return orders
}

func getPizzasIDs(pizzas []Pizza) []string {
	result := make([]string, len(pizzas))
	for i, pizza := range pizzas {
		result[i] = pizza.pizzaID
	}
	return result
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
	taken       bool
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

func howManyIngredientEquals(pizzaIngredientsA map[string]bool, pizzaIngredientsB map[string]bool) (matches int) {
	for ingrA := range pizzaIngredientsA {
		if pizzaIngredientsB[ingrA] {
			matches++
		}
	}
	return
}
