package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                            = 0.65 // средняя длина шага.
	mInKm                              = 1000 // количество метров в километре.
	minInH                             = 60   // количество минут в часе.
	runningCaloriesMeanSpeedMultiplier = 18.0
	runningCaloriesMeanSpeedShift      = 20.0
	walkingCaloriesWeightMultiplier    = 0.035
	walkingSpeedHeightMultiplier       = 0.029
)

// parseTraining разбирает строку данных, содержащую количество шагов, вид активности и продолжительность активности.
func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделить строку на слайс строк.
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int.
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}

	activity := strings.TrimSpace(parts[1])

	// Преобразовать третий элемент слайса в time.Duration.
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}

	// Проверка длительности на меньше или равно 0.
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("invalid duration: %w", err)
	}

	return steps, activity, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps)
	hours := duration.Hours()
	return dist / hours
}

// TrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	var calories float64
	var speed float64
	var distanceKm float64

	switch activity {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
		speed = meanSpeed(steps, duration)
		distanceKm = distance(steps)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
		speed = meanSpeed(steps, duration)
		distanceKm = distance(steps)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distanceKm, speed, calories)
}

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
}

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
