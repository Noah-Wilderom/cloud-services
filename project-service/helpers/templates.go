package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadTemplateFiles(folder string) ([]os.DirEntry, error) {
	return os.ReadDir(fmt.Sprintf("/templates/%s", folder))
}

func ReplaceStubVariables(file string, outputPath string, vars map[string]string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Create an output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(outputFile)

	// Read each line from the input file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into words
		words := strings.Fields(line)

		// Replace specific words in the line
		for i, word := range words {
			if strings.HasPrefix(word, "{{") {
				fmt.Println("Woord:", word)
				for placeholder, replacement := range vars {
					if strings.Contains(word, placeholder) {
						if strings.Contains(word, "{{") || strings.Contains(word, "}}") {
							words[i] = strings.Replace(word, "{{", "", -1)
							words[i] = strings.Replace(word, "}}", "", -1)
						}

						words[i] = strings.Replace(word, placeholder, replacement, -1)
						break
					}
					//if strings.Contains(strings.ToUpper(word), strings.ToUpper(placeholder)) {
					//	words[i] = strings.Replace(word, replacement, "", -1)
					//} else if strings.Contains(strings.ToUpper(word), strings.ToUpper(fmt.Sprintf("{{%s}}", placeholder))) {
					//	words[i] = strings.Replace(word, fmt.Sprintf("{{%s}}", replacement), "", -1)
					//}
				}
			}
		}

		// Join the modified words back into a line
		modifiedLine := strings.Join(words, " ")

		// Write the modified line to the output file
		_, err := fmt.Fprintln(writer, modifiedLine)
		if err != nil {
			log.Println("ERROR WRITER", err)
			return err
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	// Flush the writer to ensure all data is written to the output file
	err = writer.Flush()
	if err != nil {
		log.Println("ERROR WRITER FLUSH", err)
		return err
	}

	return nil
}
