package bevoegdheden

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/privacybydesign/kvk-extract/models"
)

func isMinderjarig(gb string) bool {
	geboorteDatum, err := time.Parse(time.RFC3339, gb[6:]+"-"+gb[3:5]+"-"+gb[0:2]+"T00:00:00+00:00")
	if err != nil {
		fmt.Println(err)
		return true
	}
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return geboorteDatum.After(eighteenYearsAgo)
}

func validateInput(kvkNummer string, identityNP models.IdentityNP) error {
	kvkRegexp := regexp.MustCompile(`^\d{8}.*$`)
	if kvkNummer == "" || !kvkRegexp.MatchString(kvkNummer) {
		fmt.Println("kvk nummer klopt niet")
		return ErrInvalidInput
	}
	if identityNP.Voornamen == "" || identityNP.Geslachtsnaam == "" || identityNP.Geboortedatum == "" || len(identityNP.Geboortedatum) != 10 {
		fmt.Println("persoonsgegeven klopt niet")
		return ErrInvalidInput
	}
	gb := identityNP.Geboortedatum
	_, err := time.Parse(time.RFC3339, gb[6:]+"-"+gb[3:5]+"-"+gb[0:2]+"T00:00:00+00:00")
	if err != nil {
		fmt.Println(err)
		return ErrInvalidInput
	}
	return nil
}

func GetBevoegdheid(kvkNummer string, identityNP models.IdentityNP, cert string, key string, useCache bool, env string) (*models.BevoegdheidResponse, error) {
	err := validateInput(kvkNummer, identityNP)
	if err != nil {
		return nil, err
	}

	ophalenInschrijvingResponse, err := GetExtract(kvkNummer, cert, key, useCache, env)
	if err != nil {
		return nil, err
	}

	bevoegdheidUittreksel := &models.BevoegdheidUittreksel{}
	paths := &models.Paths{}
	interpretatie := &models.Interpretatie{}
	bevoegdheidResponse := &models.BevoegdheidResponse{
		BevoegdheidUittreksel: bevoegdheidUittreksel,
		Paths:                 paths,
		Interpretatie:         interpretatie,
	}

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit

	bevoegdheidResponse.ExtractOriginalXML = ophalenInschrijvingResponse.ExtractOriginalXML
	bevoegdheidResponse.ExtractOriginal = ma

	getBevoegdheidUittreksel(bevoegdheidUittreksel, paths, ophalenInschrijvingResponse, identityNP)

	namePerson := fmt.Sprintf("%s %s %s", identityNP.Voornamen, identityNP.VoorvoegselGeslachtsnaam, identityNP.Geslachtsnaam)

	if bevoegdheidUittreksel.Functionaris == nil {
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s komt niet voor op de inschrijving van %s", namePerson, bevoegdheidUittreksel.KvkNummer)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.DatumUitschrijving != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s staat niet meer ingeschreven: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.DatumUitschrijving)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.RegistratieEinde != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s is niet (meer) actief: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.RegistratieEinde)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.BijzondereRechtstoestand != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een bijzondere rechtstoestand: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BijzondereRechtstoestand)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.BeperkingInRechtshandeling != "" && strings.Split(bevoegdheidUittreksel.Functionaris.BeperkingInRechtshandeling, ":")[1] != "WHOA" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een beperking in rechtshandeling: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BeperkingInRechtshandeling)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.BuitenlandseRechtstoestand != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een buitenlandse rechtstoestand: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BuitenlandseRechtstoestand)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.Overlijdensdatum != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s staat geregistreerd als overleden op %s", namePerson, bevoegdheidUittreksel.Functionaris.Overlijdensdatum)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.BijzondereRechtstoestand != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s heeft een bijzondere rechtstoestand: %s", namePerson, bevoegdheidUittreksel.Functionaris.BijzondereRechtstoestand)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.BeperkingInRechtshandeling != "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s heeft een beperking in rechtshandeling: %s", namePerson, bevoegdheidUittreksel.Functionaris.BeperkingInRechtshandeling)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.SchorsingAanvang != "" && bevoegdheidUittreksel.Functionaris.SchorsingEinde == "" {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s is geschorst sinds: %s", namePerson, bevoegdheidUittreksel.Functionaris.SchorsingAanvang)
		return bevoegdheidResponse, nil
	}

	if isMinderjarig(bevoegdheidUittreksel.Functionaris.Geboortedatum) {
		interpretatie.HeeftBeperking = true
		interpretatie.IsBevoegd = "Nee"

		if bevoegdheidUittreksel.PersoonRechtsvorm == "Eenmanszaak" && bevoegdheidUittreksel.Functionaris.TypeFunctionaris == "Eigenaar" {
			if bevoegdheidUittreksel.Functionaris.Handlichting == "" {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is een minderjarige eigenaar eenmanszaak zonder handlichting bij inschrijving %s alleen bevoegd met schriftelijke toestemming van een wettelijke vertegenwoordiger.", namePerson, bevoegdheidUittreksel.KvkNummer)
			} else {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is een minderjarige eigenaar eenmanszaak met handlichting bij inschrijving %s. Raadpleeg het Handelsregister voor meer informatie.", namePerson, bevoegdheidUittreksel.KvkNummer)
			}
		} else {
			if bevoegdheidUittreksel.Functionaris.Handlichting == "" {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is minderjarig zonder handlichting bij inschrijving %s alleen bevoegd met schriftelijke toestemming van een wettelijke vertegenwoordiger.", namePerson, bevoegdheidUittreksel.KvkNummer)
			} else {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is minderjarig met handlichting bij inschrijving %s. Raadpleeg het Handelsregister voor meer informatie.", namePerson, bevoegdheidUittreksel.KvkNummer)
			}
		}
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.PersoonRechtsvorm == "Eenmanszaak" && bevoegdheidUittreksel.Functionaris.TypeFunctionaris == "Eigenaar" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("De persoon %s is eigenaar van een eenmanszaak", namePerson)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.SoortBevoegdheid == "Alleen/zelfstandig bevoegd" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.SoortBevoegdheid)
		return bevoegdheidResponse, nil
	}
	if bevoegdheidUittreksel.Functionaris.SoortBevoegdheid == "Gezamenlijk bevoegd" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.SoortBevoegdheid)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.SoortBevoegdheid == "Onbeperkt bevoegd" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.SoortBevoegdheid)
		return bevoegdheidResponse, nil
	}
	if bevoegdheidUittreksel.Functionaris.SoortBevoegdheid == "Beperkt bevoegd" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.SoortBevoegdheid)
		return bevoegdheidResponse, nil
	}

	if bevoegdheidUittreksel.Functionaris.TypeVolmacht == "Volledige volmacht" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) heeft %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.TypeVolmacht)
		return bevoegdheidResponse, nil
	}
	if bevoegdheidUittreksel.Functionaris.TypeVolmacht == "Beperkte volmacht" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) heeft %s", namePerson, bevoegdheidUittreksel.Functionaris.Functie, bevoegdheidUittreksel.Functionaris.TypeVolmacht)
		return bevoegdheidResponse, nil
	}

	interpretatie.IsBevoegd = "Niet vastgesteld"
	interpretatie.Reden = "Geen reden gevonden"
	return bevoegdheidResponse, nil
}
