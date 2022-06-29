package bevoegdheden

import (
	"fmt"

	"github.com/privacybydesign/kvk-extract/models"
)

func GetBevoegdheid(bevoegdheidExtract *models.BevoegdheidExtract, cert string, key string, useCache bool) error {
	err := GetBevoegdheidExtract(bevoegdheidExtract, cert, key, useCache)
	if err != nil {
		return err
	}
	if bevoegdheidExtract.Functionaris == nil {
		bevoegdheidExtract.IsBevoegd = false
		bevoegdheidExtract.Reason = "Persoon komt niet voor op het uittreksel"
		return nil
	}
	if bevoegdheidExtract.OntbindingDatum != "" {
		bevoegdheidExtract.IsBevoegd = false
		bevoegdheidExtract.Reason = "Organisatie is ontbonden"
		return nil
	}
	if bevoegdheidExtract.BijzondereRechtstoestand != "" {
		bevoegdheidExtract.IsBevoegd = false
		bevoegdheidExtract.Reason = "Organisatie verkeert in bijzondere rechtstoestand: " + bevoegdheidExtract.BijzondereRechtstoestand
		return nil
	}
	if bevoegdheidExtract.Functionaris.BijzondereRechtstoestand != "" {
		bevoegdheidExtract.IsBevoegd = false
		bevoegdheidExtract.Reason = "Functionaris verkeert in bijzondere rechtstoestand: " + bevoegdheidExtract.Functionaris.BijzondereRechtstoestand
		return nil
	}
	if bevoegdheidExtract.Functionaris.Type == "Bevoegde" {
		if bevoegdheidExtract.Rechtsvorm == "Eenmanszaak" {
			bevoegdheidExtract.IsBevoegd = true
			bevoegdheidExtract.Reason = "Persoon is eigenaar van een eenmanszaak"
			return nil
		}
		if bevoegdheidExtract.Functionaris.SoortBevoegdheid == "Alleen/zelfstandig bevoegd" {
			bevoegdheidExtract.IsBevoegd = true
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) is %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.SoortBevoegdheid)
			return nil
		}
		if bevoegdheidExtract.Functionaris.SoortBevoegdheid == "Onbeperkt bevoegd" {
			bevoegdheidExtract.IsBevoegd = true
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) is %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.SoortBevoegdheid)
			return nil
		}
		if bevoegdheidExtract.Functionaris.SoortBevoegdheid == "Gezamenlijk bevoegd" {
			bevoegdheidExtract.IsBevoegd = false
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) is %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.SoortBevoegdheid)
			return nil
		}
		if bevoegdheidExtract.Functionaris.SoortBevoegdheid == "Beperkt bevoegd" {
			bevoegdheidExtract.IsBevoegd = false
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) is %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.SoortBevoegdheid)
			return nil
		}
	}
	if bevoegdheidExtract.Functionaris.Type == "Gemachtigde" {
		if bevoegdheidExtract.Functionaris.TypeVolmacht == "Volledige volmacht" {
			bevoegdheidExtract.IsBevoegd = true
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) heeft %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.TypeVolmacht)
			return nil
		}
		if bevoegdheidExtract.Functionaris.TypeVolmacht == "Beperkte volmacht" {
			bevoegdheidExtract.IsBevoegd = false
			bevoegdheidExtract.Reason = fmt.Sprintf("%s (%s) heeft %s", bevoegdheidExtract.Functionaris.Functie, bevoegdheidExtract.Functionaris.Functietitel, bevoegdheidExtract.Functionaris.TypeVolmacht)
			return nil
		}
	}

	bevoegdheidExtract.IsBevoegd = false
	bevoegdheidExtract.Reason = "Geen reden gevonden"
	return nil
}
