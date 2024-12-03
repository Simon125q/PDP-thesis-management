package validators

import (
	"log/slog"
	"thesis-management-app/types"
	"unicode"
)

func ValidateRealizedThesis(t types.RealizedThesisEntry) (types.RealizedThesisEntryErrors, bool) {
	sErr, sOk := ValidateStudent(t.Student)
	if !sOk {
		return types.RealizedThesisEntryErrors{
			Student: sErr,
		}, false
	}
	return types.RealizedThesisEntryErrors{}, true
}

func ValidateStudent(s types.Student) (types.StudentErrors, bool) {
	ok := true
	fErr, fOk := ValidateName(s.FirstName)
	lErr, lOk := ValidateName(s.LastName)
	nErr, nOk := ValidateIndex(s.StudentNumber)
	if !fOk || !lOk || !nOk {
		ok = false
	}
	return types.StudentErrors{
		StudentNumber: nErr,
		FirstName:     fErr,
		LastName:      lErr,
	}, ok
}

func ValidateIndex(index string) (string, bool) {
	if len(index) != 6 {
		slog.Info("valid index", "len", len(index))
		return "Indeks musi mieć długość 6", false
	}
	for _, char := range index {
		if !unicode.IsDigit(char) {
			return "Indeks może zawierać tylko liczby", false
		}
	}
	return "", true
}

func ValidateName(name string) (string, bool) {
	for _, char := range name {
		if unicode.IsDigit(char) {
			return "Imię nie może zawierać liczb", false
		}
	}
	return "", true
}
