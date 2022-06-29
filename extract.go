package bevoegdheden

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/privacybydesign/kvk-extract/models"
	"github.com/privacybydesign/kvk-extract/soap"
)

func GetExtract(kvkNummer, cert, key string, useCache bool) (*models.OphalenInschrijvingResponse, error) {
	cachePath := "cache-extract"
	ophalenInschrijvingResponse := models.OphalenInschrijvingResponse{}

	if useCache {
		respBody, err := os.ReadFile(cachePath + "/" + kvkNummer + ".xml")
		if err == nil {
			fmt.Println("using cache")
			envelope := soap.NewEnvelope(&ophalenInschrijvingResponse)

			if err := xml.Unmarshal(respBody, envelope); err != nil {
				panic(err)
			}
			r := envelope.Body.Content.(*models.OphalenInschrijvingResponse)

			r.ExtractOriginalXML = string(respBody)
			return r, nil
		}
	}

	fmt.Println("not using cache")

	wsseInfo, authErr := soap.NewWSSEAuthInfo(cert, key)
	if authErr != nil {
		fmt.Printf("Auth error: %s\n", authErr.Error())
		return nil, authErr
	}

	ophalenInschrijvingRequest := models.OphalenInschrijvingRequest{
		KvkNummer: kvkNummer,
	}

	soapReq := soap.NewRequest("http://es.kvk.nl/ophalenInschrijving", "https://webservices.preprod.kvk.nl/postbus2", ophalenInschrijvingRequest, &ophalenInschrijvingResponse, nil)

	soapReq.AddHeader(soap.ActionHeader{
		ID:    "_2",
		Value: "http://es.kvk.nl/ophalenInschrijving",
	})
	soapReq.AddHeader(soap.MessageIDHeader{
		ID:    "_3",
		Value: "uuid:" + uuid.New().String(),
	})
	soapReq.AddHeader(soap.ToHeader{
		ID:      "_4",
		Address: "http://es.kvk.nl/KVK-DataservicePP/2015/02",
	})

	soapReq.SignWith(wsseInfo)

	certificate, _ := tls.LoadX509KeyPair(cert, key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
			},
		},
	}

	soapClient := soap.NewClient(client)

	soapResp, err := soapClient.Do(context.Background(), soapReq)
	if err != nil {
		fmt.Printf("Unable to validate: %s\n", err.Error())
		return nil, err
	} else if soapResp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to validate (status code invalid): %d\n", soapResp.StatusCode)
		return nil, err
	} else if ophalenInschrijvingResponse.Meldingen.Fout != nil {
		fmt.Printf("SOAP fault experienced during call: %s\n", ophalenInschrijvingResponse.Meldingen.Fout.Omschrijving)
		return nil, errors.New(ophalenInschrijvingResponse.Meldingen.Fout.Omschrijving)
	}

	if useCache {
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			os.MkdirAll(cachePath, 0700)
		}
		_ = os.WriteFile(cachePath+"/"+kvkNummer+".xml", soapResp.RespBody, 0644)
	}

	ophalenInschrijvingResponse.ExtractOriginalXML = string(soapResp.RespBody)

	// ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	// jsonMA, _ := json.MarshalIndent(ma, "", "  ")
	// _ = os.WriteFile("extract.json", jsonMA, 0644)

	// data, _ := os.ReadFile("all.json")
	// bvgn := []models.Bevoegdheid{}
	// _ = json.Unmarshal(data, &bvgn)

	// if ma.HeeftAlsEigenaar.Eenmanszaak != nil {
	// 	eenmanszaak := ma.HeeftAlsEigenaar.Eenmanszaak
	// 	bvgn = append(bvgn, models.Bevoegdheid{
	// 		KvkNummer:      ma.KvkNummer,
	// 		TypeRechtsvorm: "Eenmanszaak",
	// 		Rechtsvorm:     eenmanszaak.PersoonRechtsvorm,
	// 		Person: models.Person{
	// 			LastName:    eenmanszaak.Geslachtsnaam,
	// 			FirstNames:  eenmanszaak.Voornamen,
	// 			Prefix:      eenmanszaak.VoorvoegselGeslachtsnaam,
	// 			DateOfBirth: ConvertDate(eenmanszaak.Geboortedatum),
	// 		},
	// 		Functionaris: &models.Functionaris{
	// 			Type: "Bevoegde",
	// 		},
	// 	})
	// 	for _, heeft := range eenmanszaak.Heeft {
	// 		if heeft.Gemachtigde != nil {
	// 			gemachtigde := heeft.Gemachtigde
	// 			np := gemachtigde.Door.NatuurlijkPersoon
	// 			if np == nil {
	// 				break
	// 			}
	// 			bvgn = append(bvgn, models.Bevoegdheid{
	// 				KvkNummer:      ma.KvkNummer,
	// 				TypeRechtsvorm: "Eenmanszaak",
	// 				Rechtsvorm:     eenmanszaak.PersoonRechtsvorm,
	// 				Person: models.Person{
	// 					LastName:    np.Geslachtsnaam,
	// 					FirstNames:  np.Voornamen,
	// 					Prefix:      np.VoorvoegselGeslachtsnaam,
	// 					DateOfBirth: ConvertDate(np.Geboortedatum),
	// 				},
	// 				Functionaris: &models.Functionaris{
	// 					Type:         "Gemachtigde",
	// 					TypeVolmacht: gemachtigde.Volmacht.TypeVolmacht.Omschrijving,
	// 				},
	// 			})
	// 		}
	// 	}
	// } else if ma.HeeftAlsEigenaar.Samenwerkingsverband != nil {
	// 	samenwerkingsverband := ma.HeeftAlsEigenaar.Samenwerkingsverband
	// 	for _, heeft := range samenwerkingsverband.Heeft {
	// 		if heeft.Aansprakelijke != nil {
	// 			aansprakelijke := heeft.Aansprakelijke
	// 			np := aansprakelijke.Door.NatuurlijkPersoon
	// 			if np == nil {
	// 				break
	// 			}
	// 			bvgn = append(bvgn, models.Bevoegdheid{
	// 				KvkNummer:      ma.KvkNummer,
	// 				TypeRechtsvorm: "Samenwerkingsverband",
	// 				Rechtsvorm:     samenwerkingsverband.PersoonRechtsvorm,
	// 				Person: models.Person{
	// 					LastName:    np.Geslachtsnaam,
	// 					FirstNames:  np.Voornamen,
	// 					Prefix:      np.VoorvoegselGeslachtsnaam,
	// 					DateOfBirth: ConvertDate(np.Geboortedatum),
	// 				},
	// 				Functionaris: &models.Functionaris{
	// 					Type:             "Bevoegde",
	// 					SoortBevoegdheid: aansprakelijke.Bevoegdheid.Soort.Omschrijving,
	// 				},
	// 			})
	// 		} else if heeft.Gemachtigde != nil {
	// 			gemachtigde := heeft.Gemachtigde
	// 			np := gemachtigde.Door.NatuurlijkPersoon
	// 			if np == nil {
	// 				break
	// 			}
	// 			bvgn = append(bvgn, models.Bevoegdheid{
	// 				KvkNummer:      ma.KvkNummer,
	// 				TypeRechtsvorm: "Samenwerkingsverband",
	// 				Rechtsvorm:     samenwerkingsverband.PersoonRechtsvorm,
	// 				Person: models.Person{
	// 					LastName:    np.Geslachtsnaam,
	// 					FirstNames:  np.Voornamen,
	// 					Prefix:      np.VoorvoegselGeslachtsnaam,
	// 					DateOfBirth: ConvertDate(np.Geboortedatum),
	// 				},
	// 				Functionaris: &models.Functionaris{
	// 					Type:         "Gemachtigde",
	// 					TypeVolmacht: gemachtigde.Volmacht.TypeVolmacht.Omschrijving,
	// 				},
	// 			})
	// 		}
	// 	}
	// } else if ma.HeeftAlsEigenaar.Rechtspersoon != nil {
	// 	rechtspersoon := ma.HeeftAlsEigenaar.Rechtspersoon
	// 	for _, heeft := range rechtspersoon.Heeft {
	// 		if heeft.Bestuursfunctie != nil {
	// 			bestuursfunctie := heeft.Bestuursfunctie
	// 			np := bestuursfunctie.Door.NatuurlijkPersoon
	// 			if np == nil {
	// 				break
	// 			}
	// 			bvgn = append(bvgn, models.Bevoegdheid{
	// 				KvkNummer:      ma.KvkNummer,
	// 				TypeRechtsvorm: "Rechtspersoon",
	// 				Rechtsvorm:     rechtspersoon.PersoonRechtsvorm,
	// 				Person: models.Person{
	// 					LastName:    np.Geslachtsnaam,
	// 					FirstNames:  np.Voornamen,
	// 					Prefix:      np.VoorvoegselGeslachtsnaam,
	// 					DateOfBirth: ConvertDate(np.Geboortedatum),
	// 				},
	// 				Functionaris: &models.Functionaris{
	// 					Type:             "Bevoegde",
	// 					SoortBevoegdheid: bestuursfunctie.Bevoegdheid.Soort.Omschrijving,
	// 				},
	// 			})

	// 		} else if heeft.Gemachtigde != nil {
	// 			gemachtigde := heeft.Gemachtigde
	// 			np := gemachtigde.Door.NatuurlijkPersoon
	// 			if np == nil {
	// 				break
	// 			}
	// 			bvgn = append(bvgn, models.Bevoegdheid{
	// 				KvkNummer:      ma.KvkNummer,
	// 				TypeRechtsvorm: "Rechtspersoon",
	// 				Rechtsvorm:     rechtspersoon.PersoonRechtsvorm,
	// 				Person: models.Person{
	// 					LastName:    np.Geslachtsnaam,
	// 					FirstNames:  np.Voornamen,
	// 					Prefix:      np.VoorvoegselGeslachtsnaam,
	// 					DateOfBirth: ConvertDate(np.Geboortedatum),
	// 				},
	// 				Functionaris: &models.Functionaris{
	// 					Type:         "Gemachtigde",
	// 					TypeVolmacht: gemachtigde.Volmacht.TypeVolmacht.Omschrijving,
	// 				},
	// 			})
	// 		}
	// 	}
	// }

	// jsonBvgn, _ := json.MarshalIndent(bvgn, "", "  ")
	// _ = os.WriteFile("all.json", jsonBvgn, 0644)

	// fmt.Println("KvkNummer: ", maatschappelijkeActiviteit.KvkNummer)
	// fmt.Println("Naam: ", maatschappelijkeActiviteit.Naam)
	// fmt.Println("Adres: ", maatschappelijkeActiviteit.BezoekLocatie.VolledigAdres)
	// fmt.Println("DatumOprichting: ", maatschappelijkeActiviteit.Registratie.DatumAanvang)
	// for _, emailAdres := range maatschappelijkeActiviteit.Communicatiegegevens.EmailAdres {
	// 	fmt.Println("Emailadres: " + emailAdres)
	// }
	// for _, communicatienummer := range maatschappelijkeActiviteit.Communicatiegegevens.Communicatienummer {
	// 	fmt.Println(communicatienummer.Soort.Omschrijving + ": " + communicatienummer.Nummer)
	// }
	// for _, sbiActiviteit := range maatschappelijkeActiviteit.ManifesteertZichAls.Onderneming.SbiActiviteit {
	// 	fmt.Println("Activiteit: " + sbiActiviteit.SbiCode.Code + ", " + sbiActiviteit.SbiCode.Omschrijving + ", hoofd: " + sbiActiviteit.IsHoofdactiviteit.Omschrijving)
	// }
	// for _, sbiActiviteit := range maatschappelijkeActiviteit.SbiActiviteit {
	// 	fmt.Println("Activiteit: " + sbiActiviteit.SbiCode.Code + ", " + sbiActiviteit.SbiCode.Omschrijving + ", hoofd: " + sbiActiviteit.IsHoofdactiviteit.Omschrijving)
	// }

	// if maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak != nil {
	// 	fmt.Println("Eenmanszaak.PersoonRechtsvorm: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.PersoonRechtsvorm)
	// 	fmt.Println("Eenmanszaak.Geslachtsnaam: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.Geslachtsnaam)
	// 	fmt.Println("Eenmanszaak.Geboortedatum: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.Geboortedatum)
	// 	if maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.BijzondereRechtstoestand != nil {
	// 		fmt.Println("Eenmanszaak.BijzondereRechtstoestand.Soort: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.BijzondereRechtstoestand.Soort.Omschrijving)
	// 		fmt.Println("Eenmanszaak.BijzondereRechtstoestand.Registratie.DatumAanvang: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.BijzondereRechtstoestand.Registratie.DatumAanvang)
	// 	}
	// 	for _, heeft := range maatschappelijkeActiviteit.HeeftAlsEigenaar.Eenmanszaak.Heeft {
	// 		fmt.Println("\n ")
	// 		if heeft.Gemachtigde != nil {
	// 			fmt.Println("Gemachtigde.Functie: ", heeft.Gemachtigde.Functie.Omschrijving)
	// 			if heeft.Gemachtigde.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.NatuurlijkPersoon.Voornamen: ", heeft.Gemachtigde.Door.NatuurlijkPersoon.Voornamen)
	// 				fmt.Println("Gemachtigde.Door.NatuurlijkPersoon.Geboortedatum: ", heeft.Gemachtigde.Door.NatuurlijkPersoon.Geboortedatum)
	// 			} else if heeft.Gemachtigde.Door.Rechtspersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.VolledigeNaam: ", heeft.Gemachtigde.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 			// fmt.Println("Gemachtigde.Volmacht.TypeVolmacht: ", heeft.Gemachtigde.Volmacht.TypeVolmacht.Omschrijving)
	// 		}
	// 	}
	// }

	// if maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband != nil {
	// 	fmt.Println("Samenwerkingsverband.PersoonRechtsvorm: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.PersoonRechtsvorm)
	// 	fmt.Println("Samenwerkingsverband.Ontbinding.Registratie.DatumAanvang: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.Ontbinding.Registratie.DatumAanvang)
	// 	fmt.Println("Aantal functionarissen: ", len(maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.Heeft))
	// 	if maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.BijzondereRechtstoestand != nil {
	// 		fmt.Println("Samenwerkingsverband.BijzondereRechtstoestand: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.BijzondereRechtstoestand.Soort.Omschrijving)
	// 		fmt.Println("Samenwerkingsverband.BijzondereRechtstoestand.Registratie.DatumAanvang: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.BijzondereRechtstoestand.Registratie.DatumAanvang)
	// 	}

	// 	for _, heeft := range maatschappelijkeActiviteit.HeeftAlsEigenaar.Samenwerkingsverband.Heeft {
	// 		fmt.Println("\n ")
	// 		if heeft.Aansprakelijke != nil {
	// 			fmt.Println("Aansprakelijke.Functie: ", heeft.Aansprakelijke.Functie.Omschrijving)
	// 			if heeft.Aansprakelijke.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("Aansprakelijke.Door.NatuurlijkPersoon.Geslachtsnaam: ", heeft.Aansprakelijke.Door.NatuurlijkPersoon.Geslachtsnaam)
	// 				fmt.Println("Aansprakelijke.Door.NatuurlijkPersoon.Geboortedatum: ", heeft.Aansprakelijke.Door.NatuurlijkPersoon.Geboortedatum)
	// 			} else {
	// 				fmt.Println("Aansprakelijke.Door.Rechtspersoon.VolledigeNaam: ", heeft.Aansprakelijke.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("Aansprakelijke.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.Aansprakelijke.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("Aansprakelijke.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.Aansprakelijke.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 			fmt.Println("Aansprakelijke.Bevoegdheid.Soort: ", heeft.Aansprakelijke.Bevoegdheid.Soort.Omschrijving)
	// 			fmt.Println("Aansprakelijke.Bevoegdheid.BeperkingInEuros: ", heeft.Aansprakelijke.Bevoegdheid.BeperkingInEuros.Waarde)
	// 			fmt.Println("Aansprakelijke.Bevoegdheid.OverigeBeperking: ", heeft.Aansprakelijke.Bevoegdheid.OverigeBeperking.Omschrijving)
	// 		}
	// 		if heeft.Gemachtigde != nil {
	// 			fmt.Println("Gemachtigde.Functie: ", heeft.Gemachtigde.Functie.Omschrijving)
	// 			if heeft.Gemachtigde.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.NatuurlijkPersoon.Voornamen: ", heeft.Gemachtigde.Door.NatuurlijkPersoon.Voornamen)
	// 			} else if heeft.Gemachtigde.Door.Rechtspersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.VolledigeNaam: ", heeft.Gemachtigde.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 			fmt.Println("Gemachtigde.Volmacht.TypeVolmacht: ", heeft.Gemachtigde.Volmacht.TypeVolmacht.Omschrijving)
	// 		}
	// 		if heeft.FunctionarisBijzondereRechtstoestand != nil {
	// 			fmt.Println("FunctionarisBijzondereRechtstoestand.Door.NaamPersoon.Naam: ", heeft.FunctionarisBijzondereRechtstoestand.Door.NaamPersoon.Naam)
	// 			fmt.Println("FunctionarisBijzondereRechtstoestand.Functie: ", heeft.FunctionarisBijzondereRechtstoestand.Functie.Omschrijving)
	// 		}
	// 	}
	// }

	// if maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon != nil {
	// 	fmt.Println("Rechtspersoon.PersoonRechtsvorm: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.PersoonRechtsvorm)
	// 	// fmt.Println("Rechtspersoon.DatumAkteOprichting: ", product.HeeftAlsEigenaar.Rechtspersoon.DatumAkteOprichting)
	// 	fmt.Println("Rechtspersoon.Ontbinding.Aanleiding: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.Ontbinding.Aanleiding.Omschrijving)
	// 	fmt.Println("Rechtspersoon.Ontbinding.Registratie.DatumAanvang: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.Ontbinding.Registratie.DatumAanvang)
	// 	fmt.Println("Aantal functionarissen: ", len(maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.Heeft))
	// 	if maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.BijzondereRechtstoestand != nil {
	// 		fmt.Println("Rechtspersoon.BijzondereRechtstoestand: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.BijzondereRechtstoestand.Soort.Omschrijving)
	// 		// fmt.Println("Rechtspersoon.BijzondereRechtstoestand.Registratie.DatumAanvang: ", product.HeeftAlsEigenaar.Rechtspersoon.BijzondereRechtstoestand.Registratie.DatumAanvang)
	// 	}

	// 	for _, heeft := range maatschappelijkeActiviteit.HeeftAlsEigenaar.Rechtspersoon.Heeft {
	// 		fmt.Println("\n ")
	// 		if heeft.Bestuursfunctie != nil {
	// 			fmt.Println("Bestuursfunctie.Functie: ", heeft.Bestuursfunctie.Functie.Omschrijving)
	// 			if heeft.Bestuursfunctie.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("Bestuursfunctie.Door.NatuurlijkPersoon: ", heeft.Bestuursfunctie.Door.NatuurlijkPersoon.Geslachtsnaam)
	// 			} else if heeft.Bestuursfunctie.Door.Rechtspersoon != nil {
	// 				fmt.Println("Bestuursfunctie.Door.Rechtspersoon.VolledigeNaam: ", heeft.Bestuursfunctie.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("Bestuursfunctie.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.Bestuursfunctie.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("Bestuursfunctie.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.Bestuursfunctie.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 			fmt.Println("Bestuursfunctie.Functietitel: ", heeft.Bestuursfunctie.Functietitel.Titel)
	// 			fmt.Println("Bestuursfunctie.Bevoegdheid.Soort: ", heeft.Bestuursfunctie.Bevoegdheid.Soort.Omschrijving)
	// 			fmt.Println("Bestuursfunctie.Bevoegdheid.BeperkingInEuros: ", heeft.Bestuursfunctie.Bevoegdheid.BeperkingInEuros.Waarde)
	// 			fmt.Println("Bestuursfunctie.Bevoegdheid.IsBevoegdMetAnderePersonen: ", heeft.Bestuursfunctie.Bevoegdheid.IsBevoegdMetAnderePersonen.Omschrijving)
	// 		}
	// 		if heeft.Gemachtigde != nil {
	// 			fmt.Println("Gemachtigde.Functie: ", heeft.Gemachtigde.Functie.Omschrijving)
	// 			if heeft.Gemachtigde.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.NatuurlijkPersoon: ", heeft.Gemachtigde.Door.NatuurlijkPersoon.Geslachtsnaam)
	// 			} else if heeft.Gemachtigde.Door.Rechtspersoon != nil {
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.VolledigeNaam: ", heeft.Gemachtigde.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.Gemachtigde.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.Gemachtigde.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 			fmt.Println("Gemachtigde.Volmacht.TypeVolmacht: ", heeft.Gemachtigde.Volmacht.TypeVolmacht.Omschrijving)
	// 			fmt.Println("Gemachtigde.Volmacht.BeperkteVolmacht.HeeftOverigeVolmacht: ", heeft.Gemachtigde.Volmacht.BeperkteVolmacht.HeeftOverigeVolmacht.Omschrijving)
	// 			fmt.Println("Gemachtigde.Volmacht.BeperkteVolmacht: ", heeft.Gemachtigde.Volmacht.BeperkteVolmacht.OmschrijvingOverigeVolmacht)
	// 		}
	// 		if heeft.OverigeFunctionaris != nil {
	// 			fmt.Println("OverigeFunctionaris.Functie: ", heeft.OverigeFunctionaris.Functie.Omschrijving)
	// 			if heeft.OverigeFunctionaris.Door.NatuurlijkPersoon != nil {
	// 				fmt.Println("OverigeFunctionaris.Door.NatuurlijkPersoon: ", heeft.OverigeFunctionaris.Door.NatuurlijkPersoon.Geslachtsnaam)
	// 			} else if heeft.OverigeFunctionaris.Door.Rechtspersoon != nil {
	// 				fmt.Println("OverigeFunctionaris.Door.Rechtspersoon.VolledigeNaam: ", heeft.OverigeFunctionaris.Door.Rechtspersoon.VolledigeNaam)
	// 				fmt.Println("OverigeFunctionaris.Door.Rechtspersoon.PersoonRechtsvorm: ", heeft.OverigeFunctionaris.Door.Rechtspersoon.PersoonRechtsvorm)
	// 				fmt.Println("OverigeFunctionaris.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer: ", heeft.OverigeFunctionaris.Door.Rechtspersoon.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer)
	// 			}
	// 		}
	// 		if heeft.FunctionarisBijzondereRechtstoestand != nil {
	// 			fmt.Println("FunctionarisBijzondereRechtstoestand.Door.NaamPersoon.Naam: ", heeft.FunctionarisBijzondereRechtstoestand.Door.NaamPersoon.Naam)
	// 			fmt.Println("FunctionarisBijzondereRechtstoestand.Functie: ", heeft.FunctionarisBijzondereRechtstoestand.Functie.Omschrijving)
	// 		}
	// 	}
	// }

	// if maatschappelijkeActiviteit.HeeftAlsEigenaar.BuitenlandseVennootschap != nil {
	// 	fmt.Println("BuitenlandseVennootschap.PersoonRechtsvorm: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.BuitenlandseVennootschap.PersoonRechtsvorm)
	// 	fmt.Println("BuitenlandseVennootschap.LandVanVestiging: ", maatschappelijkeActiviteit.HeeftAlsEigenaar.BuitenlandseVennootschap.LandVanVestiging.Omschrijving)
	// }

	return &ophalenInschrijvingResponse, nil
}
