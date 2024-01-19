package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM représente une image PGM.
type PGM struct {
	data         [][]uint8
	width, height int
	magicNumber  string
	max          int
}

// PBM représente une image PBM.
type PBM struct {
	data         [][]int
	width, height int
	magicNumber  string
}

// ReadPGM lit une image PGM depuis un fichier et retourne une structure représentant l'image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string

	// Fonction pour lire la prochaine ligne non commentée
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

	// Lire la première ligne non commentée pour obtenir le numéro magique
	if magicNumber, err = readNextLine(); err != nil {
		return nil, err
	}

	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("type de fichier non pris en charge")
	}

	// Lire les dimensions et la valeur maximale
	dimensions, err := readNextLine()
	if err != nil {
		return nil, err
	}

	dimTokens := strings.Fields(dimensions)
	if len(dimTokens) != 2 {
		return nil, errors.New("dimensions d'image invalides")
	}

	width, _ := strconv.Atoi(dimTokens[0])
	height, _ := strconv.Atoi(dimTokens[1])

	maxValueStr, err := readNextLine()
	if err != nil {
		return nil, err
	}

	maxValue, err := strconv.Atoi(maxValueStr)
	if err != nil {
		return nil, errors.New("valeur maximale invalide")
	}

	var data [][]uint8

	// Si l'image n'est pas vide, initialisez les données avec une slice vide
	if width > 0 && height > 0 {
		data = make([][]uint8, height)
		for i := range data {
			data[i] = make([]uint8, width)
		}

		if magicNumber == "P2" {
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
					data[i][j] = uint8(pixel)
				}
			}
		} else if magicNumber == "P5" {
			buffer := make([]byte, width*height)
			_, err := file.Read(buffer)
			if err != nil {
				return nil, err
			}

			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {
					data[i][j] = uint8(buffer[i*width+j])
				}
			}
		}
	}

	return &PGM{
		data:         data,
		width:        width,
		height:       height,
		magicNumber:  magicNumber,
		max:          maxValue,
	}, nil
}

// Size retourne la largeur et la hauteur de l'image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At retourne la valeur du pixel à la position (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set définit la valeur du pixel à la position (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save enregistre l'image PGM dans un fichier et retourne une erreur s'il y a un problème.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Écrire le numéro magique, les dimensions et la valeur maximale
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Écrire les données
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintf(writer, "%d ", pgm.data[i][j])
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j], pgm.data[pgm.height-i-1][j] = pgm.data[pgm.height-i-1][j], pgm.data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate
