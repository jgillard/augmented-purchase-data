package httptransport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func marshallResponse(data interface{}) []byte {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

func unmarshallRequest(body []byte, got interface{}) {
	err := json.Unmarshal(body, got)
	// json.unmarshall will not error if fields don't match
	// the error below will catch invalid json
	if err != nil {
		log.Fatal(err)
		return
	}
}

func craftErrorPayload(errorString string) []byte {
	errorResponse := jsonErrors{}
	errorResponse.Errors = append(errorResponse.Errors, jsonError{errorString})
	payload := marshallResponse(errorResponse)
	return payload
}

func ensureJSONFieldsPresent(res http.ResponseWriter, got, desired interface{}) bool {
	// if after unmarshall got is empty...
	if got == desired {
		fmt.Println("json field(s) missing from request")
		res.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func ensureStringFieldNonEmpty(res http.ResponseWriter, key, title string) bool {
	if title == "" {
		fmt.Println(fmt.Sprintf(`"%s" missing from request`, key))
		res.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func ensureStringFieldTitle(res http.ResponseWriter, key, title string, possibleOptionTypes []string) bool {
	isValid := false

	for _, possible := range possibleOptionTypes {
		if title == possible {
			isValid = true
			break
		}
	}

	if !isValid {
		fmt.Printf(`"%s" must be one of %v`, key, possibleOptionTypes)
		res.WriteHeader(http.StatusBadRequest)
	}

	return isValid
}

func ensureNoDuplicates(res http.ResponseWriter, key string, strings []string) bool {
	noDuplicates := true

	counts := make(map[string]int)
	for _, str := range strings {
		counts[str]++
	}

	for _, count := range counts {
		if count > 1 {
			fmt.Printf(`"%s" contains duplicate strings`, key)
			res.WriteHeader(http.StatusBadRequest)
			noDuplicates = false
		}
	}

	return noDuplicates
}
