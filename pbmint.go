package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]int
	width, height int
	magicNumber   string
}

func main() {
	pbm, err := ReadPBM("P1.txt")
	if err != nil {
		fmt.Println("erreur lors de la lecture du fichier", err)
		return
	}

	fmt.Println("Magic number:", pbm.magicNumber)
	fmt.Printf("Width: %d Height: %d\n\n", pbm.width, pbm.height)

	fmt.Println("Data:")
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			fmt.Printf("%d ", pbm.data[i][j])
		}
		fmt.Println()
	}
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string

	// Fonction pour lire la ligne suivante non commentée
	readNextLine := func() (string, error) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Ignorer les lignes vides ou les lignes commençant par "#"
			if line != "" && !strings.HasPrefix(line, "#") {
				return line, nil
			}
		}
		return "", scanner.Err()
	}

	// Lire la première ligne non commentée pour obtenir le nombre magique
	if magicNumber, err = readNextLine(); err != nil {
		return nil, err
	}

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("unsupported file type")
	}

	// Lire les dimensions
	dimensions, err := readNextLine()
	if err != nil {
		return nil, err
	}

	dimTokens := strings.Fields(dimensions)
	if len(dimTokens) != 2 {
		return nil, errors.New("invalid image dimensions")
	}

	width, _ := strconv.Atoi(dimTokens[0])
	height, _ := strconv.Atoi(dimTokens[1])

	var data [][]int
	// Si l'image est vide, initialiser les données avec une tranche vide
	if width > 0 && height > 0 {
		data = make([][]int, height)
		for i := range data {
			data[i] = make([]int, width)
		}

		if magicNumber == "P1" {
			for i := 0; i < height; i++ {
				line, err := readNextLine()
				if err != nil {
					return nil, err
				}

				tokens := strings.Fields(line)
				for j, token := range tokens {
					pixel, err := strconv.Atoi(token)
					if err != nil {
						return nil, err
					}
					data[i][j] = pixel
				}
			}
		} else if magicNumber == "P4" {
			// Calculer le nombre de bits de remplissage
			paddingBits := (8 - width%8) % 8

			// Calculer le nombre d'octets par ligne, en tenant compte du remplissage
			bytesPerRow := (width + paddingBits + 7) / 8

			// Créer un tampon pour lire les données binaires
			buffer := make([]byte, bytesPerRow)
			for i := 0; i < height; i++ {
				_, err := file.Read(buffer)
				if err != nil {
					return nil, err
				}

				// Traiter les bits du tampon
				for j := 0; j < width; j++ {
					// Obtenir l'octet contenant le bit
					byteIndex := j / 8
					bitIndex := 7 - (j % 8)
					bit := int((buffer[byteIndex] >> bitIndex) & 1)
					data[i][j] = bit
				}
			}
		}
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}
