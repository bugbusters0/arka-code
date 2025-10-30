package utils

import (
	"strconv"
	"time"
)

// FormatDate formatea una fecha a string DD/MM/YYYY
func FormatDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// ParseDate parsea un string YYYY-MM-DD a time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatMoney formatea un float a moneda con 2 decimales
func FormatMoney(amount float64) string {
	return "S/ " + strconv.FormatFloat(amount, 'f', 2, 64)
}

// StringToInt convierte string a int con valor por defecto
func StringToInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// StringToFloat convierte string a float64 con valor por defecto
func StringToFloat(s string, defaultValue float64) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return defaultValue
}

// StringToBool convierte string a bool
func StringToBool(s string) bool {
	return s == "1" || s == "true" || s == "on"
}
