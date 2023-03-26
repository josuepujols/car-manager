package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Car struct {
	id    int64
	brand string
	model string
	year  int64
}

func main() {
	db, err := ConnectDataBase()
	if err != nil {
		return
	}
	for i := 1; i == 1; {
		userPromt, _ := strconv.Atoi(strings.TrimSpace(promtText("\n 1 - Listar carros, \n 2 - Agregar un carro, \n 3 - Actualizar un carro, \n 4 - Borrar un carro, \n 5 - Salir \n :")))
		switch userPromt {
		case 1:
			listCars((db))
		case 2:
			fmt.Println(addCar(db))
		case 3:
			fmt.Println(updateCar(db))
		case 4:
			fmt.Println(deleteCar(db))
		default:
			i = 0
		}
	}
	defer db.Close()
}

// logical functions
func listCars(db *sql.DB) {
	result, err := db.Query("SELECT * FROM CARS")
	if err != nil {
		fmt.Println(err.Error())
	}
	for result.Next() {
		var car Car
		errorSelect := result.Scan(&car.id, &car.brand, &car.model, &car.year)
		if errorSelect != nil {
			fmt.Println(errorSelect.Error())
		}
		fmt.Println(car.id)
		fmt.Println(car.brand)
		fmt.Println(car.model)
		fmt.Println(car.year)
		fmt.Println("-------------------------------")
	}
}

func addCar(db *sql.DB) string {
	car := Car{
		brand: promtText("Enter the car brand: "),
		model: promtText("Enter the car model: "),
		year:  2000,
	}
	year, _ := strconv.Atoi(strings.TrimSpace(promtText("Enter the car year: ")))
	car.year = int64(year)
	query := "INSERT INTO cars (brand, model, year) VALUES (?, ?, ?)"
	errorResult := handleErrorResult(db.ExecContext(context.Background(), query, car.brand, car.model, car.year))
	if errorResult != nil {
		return "Error inserting this record " + errorResult.Error()
	}
	return "Record inserted"
}

func updateCar(db *sql.DB) string {
	id := strings.TrimSpace(promtText("Enter the id of the car that you wanna update: "))
	//local function to updatea entity in the database
	update := func(query, value string) string {
		errorResult := handleErrorResult(db.ExecContext(context.Background(), query, value))
		if errorResult != nil {
			return errorResult.Error()
		}
		return "Record updated"
	}
	updateYear := func(query string, value int) string {
		errorResult := handleErrorResult(db.ExecContext(context.Background(), query, value))
		if errorResult != nil {
			return errorResult.Error()
		}
		return "Record updated"
	}
	//Get the car to update
	var car Car
	carToUpdate := db.QueryRow("SELECT * FROM cars WHERE id = " + id)
	errorUpdate := carToUpdate.Scan(&car.id, &car.brand, &car.model, &car.year)
	if errorUpdate != nil {
		return "Error: " + errorUpdate.Error()
	}
	option, _ := strconv.Atoi(strings.TrimSpace(promtText("1 - Change the brand, 2 - Change the model, 3 - Change the year: ")))
	switch option {
	case 1:
		car.brand = promtText("Inngrese el nuevo nombre de la marca: ")
		query := "UPDATE CARS SET brand = ? WHERE id = " + id
		return update(query, car.brand)
	case 2:
		car.model = promtText("Inngrese el nuevo nombre del modelo: ")
		query := "UPDATE CARS SET model = ? WHERE id = " + id
		return update(query, car.model)
	case 3:
		year, _ := strconv.Atoi(strings.TrimSpace(promtText("Inngrese el nuevo year: ")))
		car.year = int64(year)
		query := "UPDATE CARS SET year = ? WHERE id = " + id
		return updateYear(query, int(car.year))
	default:
		return "Option is not allowed"
	}
}

func deleteCar(db *sql.DB) string {
	id := strings.TrimSpace(promtText("Enter the id of the car that you wanna delete: "))
	query := "DELETE FROM CARS WHERE ID = ?"
	errorResult := handleErrorResult(db.ExecContext(context.Background(), query, id))
	if errorResult != nil {
		return errorResult.Error()
	}
	return "Record deleted"
}

// extension functions
func handleErrorResult(result sql.Result, err error) error {
	return err
}

func promtText(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "An error occured while reading input. Please try again" + err.Error()
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}
