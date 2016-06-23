package helper

import (
	"math"
)

func PlantNormalization(PlantName string) string {
	retVal := PlantName
	switch PlantName {
	case "POWER PLANT #9":
		retVal = "PP9"
		break
	case "RABIGH POWER PLANT":
		retVal = "Rabigh"
		break
	case "Rabigh 2":
		retVal = "Rabigh"
		break
	case "Rabigh PP":
		retVal = "Rabigh"
		break
	case "Shuaiba Power Plant":
		retVal = "Shoaiba"
		break
	case "Sha'iba (CC)":
		retVal = "Shoaiba"
		break
	case "Sha'iba (SEC)":
		retVal = "Shoaiba"
		break
	case "GHAZLAN POWER PLANT":
		retVal = "Ghazlan"
		break
	case "GHZLAN":
		retVal = "Ghazlan"
		break
	case "Qurayyah Power Plant":
		retVal = "Qurayyah"
		break
	case "Qurayyah -Steam":
		retVal = "Qurayyah"
		break
	case "Qurayyah Combined Cycle Power Plant":
		retVal = "Qurayyah CC"
		break
	case "Qurayyah- Combined Cycle":
		retVal = "Qurayyah CC"
		break
	case "QurayyahCC":
		retVal = "Qurayyah CC"
		break
	case "QURAYYAH CC":
		retVal = "Qurayyah CC"
		break
	case "QurayyahPP":
		retVal = "Qurayyah"
		break
	case "Qurayyah Steam":
		retVal = "Qurayyah"
		break
	case "QURAYYAH":
		retVal = "Qurayyah"
		break
	default:
		break
	}
	return retVal
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
