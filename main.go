package main

import (
	"bufio"
	"fmt"
	"os"
	"products/products"
	"strings"
)

func main() {
	database := products.NewDatabase()
	cache := products.NewCache()
	notifier := products.NewNotifier()
	service := products.NewService(database, cache, notifier)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("1 for get product")
		fmt.Println("2 for exit")
		fmt.Println("Enter command:")
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if text == "1" {
			fmt.Println("Enter product id:")
			fmt.Print("-> ")
			id, _ := reader.ReadString('\n')
			id = strings.Replace(id, "\n", "", -1)
			product, err := service.GetProductById(id)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error occured: %s", err.Error()))
				continue
			}

			fmt.Println(fmt.Sprintf("ID:   %s", product.ID))
			fmt.Println(fmt.Sprintf("Name: %s", product.Name))
			fmt.Println("--------------------")
		}
		if text == "2" {
			break
		}
	}
}
