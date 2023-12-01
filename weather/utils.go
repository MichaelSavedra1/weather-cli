package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func input(prompt string) string { // Take user input form the ocndole
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func handleError(err error) { // Avoids repeat code by handling errors in a uniform way
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Color-code a value returned from the Met API based on it's value and type
func getColorEncoded(value string, threshold int, symbol string) string {
	// symbol arg will either hold "C", "mph" or "%", different thresholds assigned to each
	// value arg will hold an integer in string format
	val, err := strconv.Atoi(value) // convert string to int so threshold can be evealuated
	if err != nil {
		return value // return original string for non-numeric values
	}

	// define colors to return - should look to not have to assign these each time the func is called
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var colorString string
	// assign color and return val based on type and value
	if symbol == "C" {
		switch {
		case val > threshold+10:
			colorString = fmt.Sprint(red(value), red(symbol)) // Red
		case val < threshold/2:
			colorString = fmt.Sprint(red(value), red(symbol)) // Red
		case val > threshold/2 && val < threshold:
			colorString = fmt.Sprint(yellow(value), yellow(symbol)) // Yellow
		default:
			colorString = fmt.Sprint(green(value), green(symbol)) // Green
		}

	} else if symbol == "mph" {
		switch {
		case val > threshold+10:
			colorString = fmt.Sprint(red(value), red(symbol)) // Red
		case val > threshold:
			colorString = fmt.Sprint(yellow(value), yellow(symbol)) // Yellow
		default:
			colorString = fmt.Sprint(green(value), green(symbol)) // Green
		}

	} else if symbol == "%" {
		returnString := fmt.Sprint(value + "%%")
		switch {
		case val > threshold:
			colorString = fmt.Sprint(red(returnString)) // Red
		case val > threshold/2 && val < threshold:
			colorString = fmt.Sprint(yellow(returnString)) // Yellow
		default:
			colorString = fmt.Sprint(green(returnString)) // Green
		}
	} else {
		cyan := color.New(color.FgCyan).SprintFunc()
		colorString = fmt.Sprint(cyan(value))
	}

	return colorString
}

func formatColor(value string, choice string) string {
	// A more simple version of the above, not having to consider any value
	// thresholds, but instead simply choosing a color to convert the string into
	var returnString string
	switch choice {
	case "red":
		returnString = fmt.Sprint(color.New(color.FgRed).SprintFunc()(value))
	case "yellow":
		returnString = fmt.Sprint(color.New(color.FgYellow).SprintFunc()(value))
	case "green":
		returnString = fmt.Sprint(color.New(color.FgGreen).SprintFunc()(value))
	case "cyan":
		returnString = fmt.Sprint(color.New(color.FgCyan).SprintFunc()(value))
	default:
		returnString = fmt.Sprint(color.New(color.FgHiBlue).SprintFunc()(value))
	}

	return returnString

}
