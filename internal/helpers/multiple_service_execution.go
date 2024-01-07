package helpers

import (
	"strings"
)

func MutipleExec(inputPrompt string, checkString string) []string {
	var serviceList []string

	if MoreThanOneExpressionCheck(inputPrompt, checkString) {
		serviceList = ParseMoreThanOneExpression(inputPrompt, checkString)
		return serviceList
	}

	return append(serviceList, inputPrompt)
}

func MoreThanOneExpressionCheck(inputPrompt string, detectString string) bool {
	return strings.Contains(inputPrompt, detectString)
}

func ParseMoreThanOneExpression(inputPrompt string, splitParameter string) []string {
	return strings.Split(inputPrompt, splitParameter)
}

func ParseToSingleLine(input []string, joinString string) string {
	if len(input) == 0 {
		return ""
	}

	return strings.Join(input, joinString)
}
