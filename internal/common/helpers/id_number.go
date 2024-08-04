package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func ValidateIDNumber(nik string) error {
	provinceCode := nik[:2]    // 2-digit province code
	cityCode := nik[2:4]       // 2-digit city code
	birthDateCode := nik[6:12] // 6-digit birth date code (YYMMDD)
	serialNumber := nik[12:16] // 4-digit serial number

	if provinceCode < "01" || provinceCode > "99" {
		return errors.New("invalid province code")
	}

	if cityCode < "01" || cityCode > "99" {
		return errors.New("invalid city code")
	}

	// Validate birth date
	err := parseBirthDate(birthDateCode)
	if err != nil {
		return err
	}

	serial, err := strconv.Atoi(serialNumber)
	if err != nil || serial < 0 || serial > 9999 {
		return fmt.Errorf("invalid serial number")
	}
	return nil
}

func parseBirthDate(birthDateCode string) error {
	fmt.Printf("birthDateCode: %v\n", birthDateCode)
	if len(birthDateCode) != 6 {
		return errors.New("birthDateCode should be 8 digits long")
	}

	dayCode := birthDateCode[:2]
	monthCode := birthDateCode[2:4]
	yearCode := birthDateCode[4:]

	day, err := strconv.Atoi(dayCode)
	day = day - 40
	if err != nil || day < 1 || day > 31 {
		return errors.New("invalid day in birth date")
	}

	month, err := strconv.Atoi(monthCode)
	if err != nil || month < 1 || month > 12 {
		return errors.New("invalid month in birth date")
	}

	year, err := strconv.Atoi(yearCode)
	if err != nil {
		return errors.New("invalid year in birth date")
	}

	var fullYear int
	currentYear := time.Now().Year()

	if year >= 0 && year <= 22 {
		fullYear = 2000 + year
	} else {
		fullYear = 1900 + year
	}

	if fullYear > currentYear {
		return errors.New("year in birth date cannot be in the future")
	}

	if !isValidDate(fullYear, time.Month(month), day) {
		return errors.New("invalid date in birth date")
	}

	return nil
}

func isValidDate(year int, month time.Month, day int) bool {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return date.Year() == year && date.Month() == month && date.Day() == day
}
