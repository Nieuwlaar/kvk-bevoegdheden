package bevoegdheden

import (
	"errors"
	"strconv"

	"github.com/privacybydesign/kvk-extract/models"
)

var ErrPersonNotOnExtract = errors.New("person not found on extract")
var ErrExtractNotFound = errors.New("extract not found")

func convertDate(str string) string {
	return str[6:] + "-" + str[4:6] + "-" + str[0:4]
}

func isSamePerson(p1 *models.Functionaris, p2 *models.Functionaris) bool {
	// rand.Seed(time.Now().UnixNano())
	// return rand.Intn(2) == 1
	// return true

	if p1.FirstNames == "" || p2.FirstNames == "" || p1.LastName == "" || p2.LastName == "" || p1.DateOfBirth == "" || p2.DateOfBirth == "" {
		return false
	}

	return (p1.LastName == p2.LastName && p1.FirstNames == p2.FirstNames && p1.DateOfBirth == p2.DateOfBirth && p1.Prefix == p2.Prefix)
}

func heeftGemachtigde(bevoegdheidExtract *models.BevoegdheidExtract, gemachtigde *models.Gemachtigde, basePath string) bool {
	np := gemachtigde.Door.NatuurlijkPersoon
	if np == nil {
		return false
	}
	functionaryData := models.Functionaris{
		LastName:    np.Geslachtsnaam,
		FirstNames:  np.Voornamen,
		Prefix:      np.VoorvoegselGeslachtsnaam,
		DateOfBirth: convertDate(np.Geboortedatum),
	}
	if isSamePerson(bevoegdheidExtract.Functionaris, &functionaryData) {
		bevoegdheidExtract.Functionaris.Type = "Gemachtigde"
		bevoegdheidExtract.Functionaris.Functie = gemachtigde.Functie.Omschrijving
		bevoegdheidExtract.Functionaris.Functietitel = gemachtigde.Functietitel.Titel
		bevoegdheidExtract.Functionaris.TypeVolmacht = gemachtigde.Volmacht.TypeVolmacht.Omschrijving
		bevoegdheidExtract.Functionaris.HeeftOverigeVolmacht = gemachtigde.Volmacht.BeperkteVolmacht.HeeftOverigeVolmacht.Omschrijving
		bevoegdheidExtract.Functionaris.OmschrijvingOverigeVolmacht = gemachtigde.Volmacht.BeperkteVolmacht.OmschrijvingOverigeVolmacht
		bevoegdheidExtract.Functionaris.BeperkingInGeld = gemachtigde.Volmacht.BeperkteVolmacht.BeperkingInGeld.Waarde + gemachtigde.Volmacht.BeperkteVolmacht.BeperkingInGeld.Valuta.Omschrijving
		bevoegdheidExtract.Functionaris.BijzondereRechtstoestand = gemachtigde.Door.NatuurlijkPersoon.BijzondereRechtstoestand.Soort.Omschrijving

		bevoegdheidExtract.Paths.Functionaris = models.Functionaris{
			LastName:    basePath + ".door.natuurlijkPersoon.geslachtsnaam",
			FirstNames:  basePath + ".door.natuurlijkPersoon.voornamen",
			Prefix:      basePath + ".door.natuurlijkPersoon.voorvoegselGeslachtsnaam",
			DateOfBirth: basePath + ".door.natuurlijkPersoon.geboortedatum",
		}
		bevoegdheidExtract.Paths.Functionaris.Type = basePath
		bevoegdheidExtract.Paths.Functionaris.Functie = basePath + ".functie.omschrijving"
		bevoegdheidExtract.Paths.Functionaris.Functietitel = basePath + ".functietitel.titel"
		bevoegdheidExtract.Paths.Functionaris.TypeVolmacht = basePath + ".volmacht.typeVolmacht.omschrijving"
		bevoegdheidExtract.Paths.Functionaris.HeeftOverigeVolmacht = basePath + ".volmacht.beperkteVolmacht.heeftOverigeVolmacht.omschrijving"
		bevoegdheidExtract.Paths.Functionaris.OmschrijvingOverigeVolmacht = basePath + ".volmacht.beperkteVolmacht.omschrijvingOverigeVolmacht"
		bevoegdheidExtract.Paths.Functionaris.BeperkingInGeld = basePath + ".volmacht.beperkteVolmacht.beperkingInGeld"
		bevoegdheidExtract.Paths.Functionaris.BijzondereRechtstoestand = basePath + ".door.natuurlijkPersoon.bijzondereRechtstoestand.soort.omschrijving"
		return true
	}
	return false
}

// GetBevoegdheidExtract TODO
func GetBevoegdheidExtract(bevoegdheidExtract *models.BevoegdheidExtract, cert string, key string, useCache bool) error {
	ophalenInschrijvingResponse, err := GetExtract(bevoegdheidExtract.KvkNummer, cert, key, useCache)
	if err != nil {
		return ErrExtractNotFound
	}

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit

	bevoegdheidExtract.ExtractOriginalXML = ophalenInschrijvingResponse.ExtractOriginalXML
	bevoegdheidExtract.ExtractOriginal = ma

	bevoegdheidExtract.Paths = &models.Paths{}
	bevoegdheidExtract.Paths.KvkNummer = "maatschappelijkeActiviteit.kvkNummer"

	bevoegdheidExtract.Naam = ma.Naam
	bevoegdheidExtract.Paths.Naam = "maatschappelijkeActiviteit.naam"
	bevoegdheidExtract.Adres = ma.BezoekLocatie.VolledigAdres
	bevoegdheidExtract.Paths.Adres = "maatschappelijkeActiviteit.bezoekLocatie.volledigAdres"
	bevoegdheidExtract.RegistratieAanvang = ma.Registratie.DatumAanvang
	bevoegdheidExtract.Paths.RegistratieAanvang = "maatschappelijkeActiviteit.registratie.datumAanvang"

	if len(ma.Communicatiegegevens.EmailAdres) != 0 {
		bevoegdheidExtract.EmailAdres = ma.Communicatiegegevens.EmailAdres[0]
		bevoegdheidExtract.Paths.EmailAdres = "maatschappelijkeActiviteit.communicatiegegevens.emailAdres.0"
	}

	for i, nr := range ma.Communicatiegegevens.Communicatienummer {
		if nr.Soort.Code == "T" {
			bevoegdheidExtract.Telefoon = nr.Toegangscode + " " + nr.Nummer[1:]
			bevoegdheidExtract.Paths.Telefoon = "maatschappelijkeActiviteit.communicatiegegevens.communicatienummer." + strconv.Itoa(i)
			break
		}
	}

	for i, sbi := range ma.SbiActiviteit {
		if sbi.IsHoofdactiviteit.Code == "J" {
			bevoegdheidExtract.SbiActiviteit = sbi.SbiCode.Code + ", " + sbi.SbiCode.Omschrijving
			bevoegdheidExtract.Paths.SbiActiviteit = "maatschappelijkeActiviteit.sbiActiviteit." + strconv.Itoa(i) + ".sbiCode"
			break
		}
	}

	if bevoegdheidExtract.SbiActiviteit == "" {
		for i, sbi := range ma.ManifesteertZichAls.Onderneming.SbiActiviteit {
			if sbi.IsHoofdactiviteit.Code == "J" {
				bevoegdheidExtract.SbiActiviteit = sbi.SbiCode.Code + ", " + sbi.SbiCode.Omschrijving
				bevoegdheidExtract.Paths.SbiActiviteit = "maatschappelijkeActiviteit.manifesteertZichAls.onderneming.sbiActiviteit." + strconv.Itoa(i) + ".sbiCode"
				break
			}
		}
	}

	if ma.HeeftAlsEigenaar.Eenmanszaak != nil {
		eenmanszaak := ma.HeeftAlsEigenaar.Eenmanszaak
		bevoegdheidExtract.TypeRechtsvorm = "Eenmanszaak"
		bevoegdheidExtract.Rechtsvorm = eenmanszaak.PersoonRechtsvorm
		basePath := "maatschappelijkeActiviteit.heeftAlsEigenaar.natuurlijkPersoon"
		bevoegdheidExtract.Paths.Rechtsvorm = basePath + ".persoonRechtsvorm"
		functionaryData := models.Functionaris{
			LastName:    eenmanszaak.Geslachtsnaam,
			FirstNames:  eenmanszaak.Voornamen,
			Prefix:      eenmanszaak.VoorvoegselGeslachtsnaam,
			DateOfBirth: convertDate(eenmanszaak.Geboortedatum),
		}

		if eenmanszaak.BijzondereRechtstoestand != nil {
			bevoegdheidExtract.BijzondereRechtstoestand = eenmanszaak.BijzondereRechtstoestand.Soort.Omschrijving
			bevoegdheidExtract.Paths.BijzondereRechtstoestand = basePath + ".bijzondereRechtstoestand.soort.omschrijving"
			bevoegdheidExtract.BijzondereRechtstoestandDatum = eenmanszaak.BijzondereRechtstoestand.Registratie.DatumAanvang
			bevoegdheidExtract.Paths.BijzondereRechtstoestandDatum = basePath + ".bijzondereRechtstoestand.registratie.datumAanvang"
		}
		if isSamePerson(bevoegdheidExtract.Functionaris, &functionaryData) {
			bevoegdheidExtract.Functionaris.Type = "Eigenaar"
			bevoegdheidExtract.Functionaris.SoortBevoegdheid = "Onbeperkt bevoegd"

			bevoegdheidExtract.Paths.Functionaris = models.Functionaris{
				LastName:    basePath + ".geslachtsnaam",
				FirstNames:  basePath + ".voornamen",
				Prefix:      basePath + ".voorvoegselGeslachtsnaam",
				DateOfBirth: basePath + ".geboortedatum",
			}
		} else {
			for i, heeft := range eenmanszaak.Heeft {
				if heeft.Gemachtigde != nil {
					found := heeftGemachtigde(bevoegdheidExtract, heeft.Gemachtigde, basePath+".heeft."+strconv.Itoa(i)+".gemachtigde")
					if found {
						break
					}
				}
				// } else if heeft.FunctionarisBijzondereRechtstoestand != nil {
				// TODO zouden we dit kunnen matchen?
			}
		}
	} else if ma.HeeftAlsEigenaar.Samenwerkingsverband != nil {
		samenwerkingsverband := ma.HeeftAlsEigenaar.Samenwerkingsverband
		bevoegdheidExtract.TypeRechtsvorm = "Samenwerkingsverband"
		bevoegdheidExtract.Rechtsvorm = samenwerkingsverband.PersoonRechtsvorm

		basePath := "maatschappelijkeActiviteit.heeftAlsEigenaar.samenwerkingsverband"
		bevoegdheidExtract.Paths.Rechtsvorm = basePath + ".persoonRechtsvorm"
		if samenwerkingsverband.BijzondereRechtstoestand != nil {
			bevoegdheidExtract.BijzondereRechtstoestand = samenwerkingsverband.BijzondereRechtstoestand.Soort.Omschrijving
			bevoegdheidExtract.BijzondereRechtstoestandDatum = samenwerkingsverband.BijzondereRechtstoestand.Registratie.DatumAanvang

			bevoegdheidExtract.Paths.BijzondereRechtstoestand = basePath + ".bijzondereRechtstoestand.soort.omschrijving"
			bevoegdheidExtract.Paths.BijzondereRechtstoestandDatum = basePath + ".bijzondereRechtstoestand.registratie.datumAanvang"
		}
		bevoegdheidExtract.OntbindingDatum = samenwerkingsverband.Ontbinding.Registratie.DatumAanvang
		bevoegdheidExtract.OntbindingAanleiding = samenwerkingsverband.Ontbinding.Aanleiding.Omschrijving

		bevoegdheidExtract.Paths.OntbindingDatum = basePath + ".ontbinding.registratie.datumAanvang"
		bevoegdheidExtract.Paths.OntbindingAanleiding = basePath + ".ontbinding.aanleiding.omschrijving"

		for i, heeft := range samenwerkingsverband.Heeft {
			if heeft.Aansprakelijke != nil {
				aansprakelijke := heeft.Aansprakelijke
				np := aansprakelijke.Door.NatuurlijkPersoon
				if np == nil {
					continue
				}
				functionaryData := models.Functionaris{
					LastName:    np.Geslachtsnaam,
					FirstNames:  np.Voornamen,
					Prefix:      np.VoorvoegselGeslachtsnaam,
					DateOfBirth: convertDate(np.Geboortedatum),
				}
				if isSamePerson(bevoegdheidExtract.Functionaris, &functionaryData) {
					bevoegdheidExtract.Functionaris.Type = "Aansprakelijke"
					bevoegdheidExtract.Functionaris.Functie = aansprakelijke.Functie.Omschrijving
					bevoegdheidExtract.Functionaris.SoortBevoegdheid = aansprakelijke.Bevoegdheid.Soort.Omschrijving
					bevoegdheidExtract.Functionaris.BeperkingInEuros = aansprakelijke.Bevoegdheid.BeperkingInEuros.Waarde + aansprakelijke.Bevoegdheid.BeperkingInEuros.Valuta.Omschrijving
					bevoegdheidExtract.Functionaris.OverigeBeperking = aansprakelijke.Bevoegdheid.OverigeBeperking.Omschrijving
					bevoegdheidExtract.Functionaris.IsBevoegdMetAnderePersonen = aansprakelijke.Bevoegdheid.IsBevoegdMetAnderePersonen.Omschrijving
					bevoegdheidExtract.Functionaris.BijzondereRechtstoestand = aansprakelijke.Door.NatuurlijkPersoon.BijzondereRechtstoestand.Soort.Omschrijving

					basePath = basePath + ".heeft." + strconv.Itoa(i) + ".aansprakelijke"
					bevoegdheidExtract.Paths.Functionaris = models.Functionaris{
						LastName:    basePath + ".door.natuurlijkPersoon.geslachtsnaam",
						FirstNames:  basePath + ".door.natuurlijkPersoon.voornamen",
						Prefix:      basePath + ".door.natuurlijkPersoon.voorvoegselGeslachtsnaam",
						DateOfBirth: basePath + ".door.natuurlijkPersoon.geboortedatum",
					}
					bevoegdheidExtract.Paths.Functionaris.Type = basePath
					bevoegdheidExtract.Paths.Functionaris.Functie = basePath + ".functie.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.SoortBevoegdheid = basePath + ".bevoegdheid.soort.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.BeperkingInEuros = basePath + ".bevoegdheid.beperkingInEuros"
					bevoegdheidExtract.Paths.Functionaris.OverigeBeperking = basePath + ".bevoegdheid.overigeBeperking.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.IsBevoegdMetAnderePersonen = basePath + ".bevoegdheid.isBevoegdMetAnderePersonen.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.BijzondereRechtstoestand = basePath + ".door.natuurlijkPersoon.bijzondereRechtstoestand.soort.omschrijving"

					break
				}
			} else if heeft.Gemachtigde != nil {
				found := heeftGemachtigde(bevoegdheidExtract, heeft.Gemachtigde, "maatschappelijkeActiviteit.heeftAlsEigenaar.samenwerkingsverband.heeft."+strconv.Itoa(i)+".gemachtigde")

				if found {
					break
				}
			}
			// } else if heeft.FunctionarisBijzondereRechtstoestand != nil {
			// TODO zouden we dit kunnen matchen?
		}
	} else if ma.HeeftAlsEigenaar.Rechtspersoon != nil {
		rechtspersoon := ma.HeeftAlsEigenaar.Rechtspersoon
		bevoegdheidExtract.TypeRechtsvorm = "Rechtspersoon"
		bevoegdheidExtract.Rechtsvorm = rechtspersoon.PersoonRechtsvorm
		basePath := "maatschappelijkeActiviteit.heeftAlsEigenaar.rechtspersoon"
		bevoegdheidExtract.Paths.Rechtsvorm = basePath + ".persoonRechtsvorm"
		if rechtspersoon.BijzondereRechtstoestand != nil {
			bevoegdheidExtract.BijzondereRechtstoestand = rechtspersoon.BijzondereRechtstoestand.Soort.Omschrijving
			bevoegdheidExtract.BijzondereRechtstoestandDatum = rechtspersoon.BijzondereRechtstoestand.Registratie.DatumAanvang

			bevoegdheidExtract.Paths.BijzondereRechtstoestand = basePath + ".bijzondereRechtstoestand.soort.omschrijving"
			bevoegdheidExtract.Paths.BijzondereRechtstoestandDatum = basePath + ".bijzondereRechtstoestand.registratie.datumAanvang"
		}
		bevoegdheidExtract.OntbindingDatum = rechtspersoon.Ontbinding.Registratie.DatumAanvang
		bevoegdheidExtract.OntbindingAanleiding = rechtspersoon.Ontbinding.Aanleiding.Omschrijving

		bevoegdheidExtract.Paths.OntbindingDatum = basePath + ".ontbinding.registratie.datumAanvang"
		bevoegdheidExtract.Paths.OntbindingAanleiding = basePath + ".ontbinding.aanleiding.omschrijving"

		for i, heeft := range rechtspersoon.Heeft {
			if heeft.Bestuursfunctie != nil {
				bestuursfunctie := heeft.Bestuursfunctie
				np := bestuursfunctie.Door.NatuurlijkPersoon
				if np == nil {
					continue
				}
				functionaryData := models.Functionaris{
					LastName:    np.Geslachtsnaam,
					FirstNames:  np.Voornamen,
					Prefix:      np.VoorvoegselGeslachtsnaam,
					DateOfBirth: convertDate(np.Geboortedatum),
				}
				if isSamePerson(bevoegdheidExtract.Functionaris, &functionaryData) {
					bevoegdheidExtract.Functionaris.Type = "Bestuursfunctie"
					bevoegdheidExtract.Functionaris.Functie = bestuursfunctie.Functie.Omschrijving
					bevoegdheidExtract.Functionaris.Functietitel = bestuursfunctie.Functietitel.Titel
					bevoegdheidExtract.Functionaris.SoortBevoegdheid = bestuursfunctie.Bevoegdheid.Soort.Omschrijving
					bevoegdheidExtract.Functionaris.BeperkingInEuros = bestuursfunctie.Bevoegdheid.BeperkingInEuros.Waarde + bestuursfunctie.Bevoegdheid.BeperkingInEuros.Valuta.Omschrijving
					bevoegdheidExtract.Functionaris.OverigeBeperking = bestuursfunctie.Bevoegdheid.OverigeBeperking.Omschrijving
					bevoegdheidExtract.Functionaris.IsBevoegdMetAnderePersonen = bestuursfunctie.Bevoegdheid.IsBevoegdMetAnderePersonen.Omschrijving
					bevoegdheidExtract.Functionaris.BijzondereRechtstoestand = bestuursfunctie.Door.NatuurlijkPersoon.BijzondereRechtstoestand.Soort.Omschrijving

					basePath = basePath + ".heeft." + strconv.Itoa(i) + ".bestuursfunctie"
					bevoegdheidExtract.Paths.Functionaris = models.Functionaris{
						LastName:    basePath + ".door.natuurlijkPersoon.geslachtsnaam",
						FirstNames:  basePath + ".door.natuurlijkPersoon.voornamen",
						Prefix:      basePath + ".door.natuurlijkPersoon.voorvoegselGeslachtsnaam",
						DateOfBirth: basePath + ".door.natuurlijkPersoon.geboortedatum",
					}
					bevoegdheidExtract.Paths.Functionaris.Type = basePath
					bevoegdheidExtract.Paths.Functionaris.Functie = basePath + ".functie.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.Functietitel = basePath + ".functietitel.titel"
					bevoegdheidExtract.Paths.Functionaris.SoortBevoegdheid = basePath + ".bevoegdheid.soort.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.BeperkingInEuros = basePath + ".bevoegdheid.beperkingInEuros"
					bevoegdheidExtract.Paths.Functionaris.OverigeBeperking = basePath + ".bevoegdheid.overigeBeperking.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.IsBevoegdMetAnderePersonen = basePath + ".bevoegdheid.isBevoegdMetAnderePersonen.omschrijving"
					bevoegdheidExtract.Paths.Functionaris.BijzondereRechtstoestand = basePath + ".door.natuurlijkPersoon.bijzondereRechtstoestand.soort.omschrijving"

					break
				}
			} else if heeft.Gemachtigde != nil {
				found := heeftGemachtigde(bevoegdheidExtract, heeft.Gemachtigde, "maatschappelijkeActiviteit.heeftAlsEigenaar.rechtspersoon.heeft."+strconv.Itoa(i)+".gemachtigde")
				if found {
					break
				}
			}
			// } else if heeft.FunctionarisBijzondereRechtstoestand != nil {
			// TODO zouden we dit kunnen matchen?
		}
	}
	// else if product.HeeftAlsEigenaar.BuitenlandseVennootschap != nil {
	// 	bevoegdheid.Rechtsvorm = product.HeeftAlsEigenaar.BuitenlandseVennootschap.PersoonRechtsvorm
	// }

	if bevoegdheidExtract.Functionaris.Type == "" {
		// Person not found on extract
		return ErrPersonNotOnExtract
	}

	return nil
}
