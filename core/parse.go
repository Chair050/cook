package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/giteshnxtlvl/cook/parse"
	"github.com/giteshnxtlvl/cook2/core"
)

var (
	Help       = false
	Verbose    = false
	ConfigPath = ""
	UpperCase  = false
	LowerCase  = false
)

var columnCases = make(map[int]map[string]bool)

func UpdateCases(caseValue string, noOfColumns int) map[int]map[string]bool {
	caseValue = strings.ToUpper(caseValue)

	for i := 0; i < noOfColumns; i++ {
		columnCases[i] = make(map[string]bool)
	}

	//Column Wise Cases
	if strings.Contains(caseValue, ":") {
		for _, val := range strings.Split(caseValue, ",") {
			v := strings.SplitN(val, ":", 2)
			i, err := strconv.Atoi(v[0])
			if err != nil {
				log.Fatalf("Err: Invalid column index %s", v[0])
			}
			for _, j := range strings.Split(v[1], "") {
				columnCases[i][j] = true
			}
		}
	} else {
		//Global Cases
		all := false
		if caseValue == "A" {
			all = true
			caseValue = ""
		}

		if all || strings.Contains(caseValue, "C") {
			columnCases[0]["L"] = true
			for i := 1; i < noOfColumns; i++ {
				columnCases[i]["T"] = true
			}
			caseValue = strings.ReplaceAll(caseValue, "C", "")
		}

		if all || strings.Contains(caseValue, "U") {
			UpperCase = true
			caseValue = strings.ReplaceAll(caseValue, "U", "")
		}

		if all || strings.Contains(caseValue, "L") {
			LowerCase = true
			caseValue = strings.ReplaceAll(caseValue, "L", "")
		}

		if all || strings.Contains(caseValue, "T") {
			for i := 0; i < noOfColumns; i++ {
				columnCases[i]["T"] = true
			}
		}

	}

	return columnCases
}

func PrintPattern(k string, v []string, search string) {
	fmt.Println(strings.ReplaceAll(k, search, "\u001b[48;5;239m"+search+core.Reset))
	fmt.Printf("    %s%s{\n", k, strings.ReplaceAll(v[0], search, core.Blue+search+core.Reset))
	for _, file := range v[1:] {
		fmt.Printf("\t%s\n", strings.ReplaceAll(file, search, core.Blue+search+core.Reset))
	}
	fmt.Println("    }")
}

//Checking for patterns/functions
func ParseFunc(value string, array *[]string) bool {

	if !(strings.Contains(value, "(") && strings.Contains(value, ")")) {
		return false
	}

	funcName, funcArgs := parse.ReadCrBrSepBy(value, ",")
	// fmt.Println(funcName)
	// fmt.Println(funcValues)

	fmt.Print("")

	if funcPatterns, exists := M["patterns"][funcName]; exists {

		funcDef := strings.Split(funcPatterns[0][1:len(funcPatterns[0])-1], ",")

		// fmt.Printf("Func Arg: %v", funcArgs)
		// fmt.Printf("\tFunc Def: %v", funcDef)

		if len(funcDef) != len(funcArgs) {
			log.Fatalln("\nErr: No of Arguments are different for")
			PrintPattern(funcName, funcPatterns, funcName)
		}

		for _, p := range funcPatterns[1:] {
			for index, arg := range funcDef {
				p = strings.ReplaceAll(p, arg, funcArgs[index])
			}
			*array = append(*array, p)
		}

		return true
	}
	return false
}

var InputFile = make(map[string]bool)

func ParseFile(param string, value string, array *[]string) bool {

	// Checking for file
	if InputFile[param] && !strings.Contains(value, ":") {
		// AddFilesToArray(value, array)
		FileValues(value, array)
		return true
	}

	if checkFileInYaml(value, array) {
		return true
	}

	// Checking for File and Regex
	if strings.Contains(value, ":") {
		// File may starts from E: C: D: for windows + Regex is supplied
		if strings.Count(value, ":") == 2 {
			tmp := strings.SplitN(value, ":", 3)

			one, two, three := tmp[0], tmp[1], tmp[2]
			test1, test2 := one+":"+two, two+":"+three

			if _, err := os.Stat(test1); err == nil {
				FileRegex(test1, three, array)
				return true
			} else if _, err := os.Stat(test2); err == nil {
				FileRegex(one, test2, array)
				return true
			}
		}

		// if strings.Count(value, ":") == 1 {
		// 	if _, err := os.Stat(value); err == nil {
		// 		AddFilesToArray(value, array)
		// 		return true
		// 	}
		// 	t := strings.SplitN(value, ":", 2)
		// 	file, reg := t[0], t[1]

		// 	if strings.HasSuffix(file, ".txt") {
		// 		FileRegex([]string{file}, reg, array)
		// 		return true
		// 	} else if files, exists := M["files"][file]; exists {
		// 		FileRegex(files, reg, array)
		// 		return true
		// 	}
		// }
	}
	return false
}

var pipe []string

func PipeInput(value string, array *[]string) bool {
	if value == "-" {
		sc := bufio.NewScanner(os.Stdin)
		if len(pipe) > 0 {
			*array = append(*array, pipe...)
		}
		for sc.Scan() {
			*array = append(*array, sc.Text())
			pipe = append(pipe, sc.Text())
		}
		return true
	}
	return false
}

func RawInput(value string, array *[]string) bool {
	if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
		lv := len(value)
		*array = append(*array, []string{value[1 : lv-1]}...)
		return true
	}
	return false
}

func ParseRanges(p string, array *[]string) bool {

	success := false
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if strings.Count(p, "-") == 1 {

		numRange := strings.SplitN(p, "-", 2)
		from := numRange[0]
		to := numRange[1]

		start, err1 := strconv.Atoi(from)
		stop, err2 := strconv.Atoi(to)

		if err1 == nil && err2 == nil {
			for start <= stop {
				*array = append(*array, strconv.Itoa(start))
				start++
			}
			success = true
		}

		if !success && len(from) == 1 && len(to) == 1 && strings.Contains(chars, from) && strings.Contains(chars, to) {
			start = strings.Index(chars, from)
			stop = strings.Index(chars, to)

			if start < stop {
				charsList := strings.Split(chars, "")
				for start <= stop {
					*array = append(*array, charsList[start])
					start++
				}
				success = true
			}
		}
	}
	return success
}

func ParsePorts(ports []string, array *[]string) {

	for _, p := range ports {
		if ParseRanges(p, array) {
			continue
		}
		port, err := strconv.Atoi(p)
		if err != nil {
			log.Printf("Err: Is this port number -_-?? '%s'", p)
		}
		*array = append(*array, strconv.Itoa(port))
	}
}
