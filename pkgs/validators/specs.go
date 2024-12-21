package validators

import "thesis-management-app/types"

func ValidateSpecialization(spec types.Specialization) (types.SpecializationErrors, bool) {
	ok := true
	fErr, fOk := ValidateSpecName(spec.Name)

	if !fOk {
		ok = false
	}
	return types.SpecializationErrors{
		Name: fErr,
	}, ok
}

func ValidateSpecName(specName string) (string, bool) {
	return "", true
	//me love some good checking
}
