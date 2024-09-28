package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var apiURL = "https://api-inference.huggingface.co/models/"
var apiKey = os.Getenv("huggingface_api_key")

type ImageRequest struct {
	Inputs string `json:"inputs"`
}

func generateImage(prompt string) ([]byte, error) {
	// Prepare the request payload
	requestBody, err := json.Marshal(ImageRequest{
		Inputs: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %v", err)
	}

	// Create the POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers, including the API token for authentication
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors in response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response code: %v", resp.Status)
	}

	// Read the response body (the image data)
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return responseData, nil
}

// readFileAndGenerateImages reads the text file, splits paragraphs, and generates images for each paragraph
func readFileAndGenerateImages(fileName string) error {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file contents
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Buffer to accumulate lines of a paragraph
	var paragraphBuilder strings.Builder
	var paragraphIndex int

	// Loop through the lines in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is empty, indicating the end of a paragraph
		if line == "" {
			paragraph := strings.TrimSpace(paragraphBuilder.String())
			if paragraph != "" {
				err := processParagraph(paragraph, paragraphIndex)
				if err != nil {
					return err
				}
				paragraphIndex++
			}
			paragraphBuilder.Reset()
		} else {
			paragraphBuilder.WriteString(line + " ")
		}
	}

	// Handle the last paragraph if the file doesn't end with a blank line
	paragraph := strings.TrimSpace(paragraphBuilder.String())
	if paragraph != "" {
		err := processParagraph(paragraph, paragraphIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

// processParagraph sends the paragraph as a prompt to generate an image and saves the result
func processParagraph(paragraph string, index int) error {
	fmt.Printf("Generating image for paragraph %d: %s\n", index+1, paragraph)

	// Generate the image based on the paragraph
	imageData, err := generateImage("Generate a cartoon illustration for the following text: " + paragraph)
	if err != nil {
		return fmt.Errorf("error generating image for paragraph %d: %v", index+1, err)
	}

	// Save the image to a file
	imageFileName := fmt.Sprintf("generated_image_%d.png", index+1)
	err = ioutil.WriteFile(imageFileName, imageData, 0644)
	if err != nil {
		return fmt.Errorf("error saving image for paragraph %d: %v", index+1, err)
	}

	fmt.Printf("Image for paragraph %d successfully saved as %s\n", index+1, imageFileName)
	return nil
}

func Execute() {
	fileName := flag.String("f", "file", "Specify the text file to read")
	module := flag.String("m", "black-forest-labs/FLUX.1-schnell", "Specify the model to use")
	flag.Parse()
	apiURL += *module

	if *fileName == "" || *module == "" {
		usage()
	}

	// Read the file and generate images for each paragraph
	err := readFileAndGenerateImages(*fileName)
	if err != nil {
		log.Fatalf("Error processing file: %v", err)
	}

	fmt.Println("All images have been generated and saved.")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("Example: huggingface -f input.txt -m black-forest-labs/FLUX.1-schnell")
	fmt.Println("RECOMMENDED MODELS")
	fmt.Println("black-forest-labs/FLUX.1-dev")
	fmt.Println("black-forest-labs/FLUX.1-schnell")
	fmt.Println("stabilityai/stable-diffusion-xl-base-1.0")
}
