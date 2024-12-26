package validators

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/types"
	"unicode"
)

func ValidateRealizedThesis(t types.RealizedThesisEntry) (types.RealizedThesisEntryErrors, bool) {
	sErr, sOk := ValidateStudent(t.Student)
	thesisNumberErr, tnOk := ValidateThesisNumber(t.ThesisNumber)
	supErr, supOk := ValidateEmployee(t.Supervisor)
	asErr, asOk := ValidateEmployee(t.AssistantSupervisor)
	rErr, rOk := ValidateEmployee(t.Reviewer)
	chErr, chOk := ValidateEmployee(t.Chair)
	hErr, hOk := ValidateHourlySettlement(t.HourlySettlement, t.Student.Degree)
	if !sOk || !tnOk || !supOk || !asOk || !rOk || !chOk || !hOk {
		return types.RealizedThesisEntryErrors{
			ThesisNumber:        thesisNumberErr,
			Student:             sErr,
			Supervisor:          supErr,
			AssistantSupervisor: asErr,
			Reviewer:            rErr,
			Chair:               chErr,
			HourlySettlement:    hErr,
		}, false
	}
	return types.RealizedThesisEntryErrors{}, true
}

func ValidateEmployee(e types.UniversityEmployeeEntry) (types.UniversityEmployeeEntryErrors, bool) {
	ok := true
	fErr, fOk := ValidateName(e.FirstName)
	lErr, lOk := ValidateName(e.LastName)
	if !fOk || !lOk {
		ok = false
	}
	return types.UniversityEmployeeEntryErrors{
		FirstName: fErr,
		LastName:  lErr,
	}, ok
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
	if len(index) != 6 && len(index) != 7 {
		slog.Info("valid index", "len", len(index))
		return "Indeks musi mieć długość 6 lub 7", false
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

func ValidateHourlySettlement(hours types.HourlySettlement, studyLevel string) (types.HourlySettlementErrors, bool) {
	hSum := hours.SupervisorHours + hours.AssistantSupervisorHours
	slog.Info("ValidateHourlySettlement", "hours sum", hSum)
	if studyLevel == "I stopień" {
		if hSum != 10 {
			return types.HourlySettlementErrors{Total: "Godziny promotorów w pracy inżynierskiej muszą sumowac sie do 10"}, false
		}
	} else if studyLevel == "II stopień" {
		if hSum != 15 {
			return types.HourlySettlementErrors{Total: "Godziny promotorów pracy magisterskiej muszą sumowac sie do 15"}, false
		}
	}
	return types.HourlySettlementErrors{}, true
}

func ValidateThesisNumber(num string) (string, bool) {
	pattern := `^[A-Za-z]\d+\/[A-Za-ż]{3}\/\d+\/\d{4}$`
	re := regexp.MustCompile(pattern)
	if re.MatchString(num) {
		return "", true
	}
	return "Numer pracy musi miec odpowiedni format: Number jednostki/stopien/numer pracy/rok", false
}

func CheckThesisNumber(thesNumber, degree string) string {
	slog.Info("CheckThesisNumber", "thesNumber", thesNumber)
	pattern := `^[A-Za-z]\d+\/stopien\/num\/\d{4}$`
	re := regexp.MustCompile(pattern)
	if re.MatchString(thesNumber) {
		slog.Info("CheckThesisNumber", "MatchedthesNumber", thesNumber)
		if degree == "I stopień" {
			thesNumber = strings.Replace(thesNumber, "stopien", "inż", -1)
		} else if degree == "II stopień" {
			thesNumber = strings.Replace(thesNumber, "stopien", "mgr", -1)
		}
		parts := strings.Split(thesNumber, "/")
		nextNum, err := server.MyS.DB.GetNextThesisNumber(parts[0], parts[1], parts[3])
		if err != nil {
			slog.Error("CheckThesisNumber", "err", err)
			return thesNumber
		}
		thesNumber = strings.Replace(thesNumber, "num", strconv.Itoa(nextNum), -1)
		return thesNumber
	}
	return thesNumber
}
