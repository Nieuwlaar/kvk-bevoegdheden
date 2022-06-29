package models

import "encoding/xml"

type OphalenInschrijvingRequest struct {
	XMLName   xml.Name `xml:"ophalenInschrijvingRequest" json:"ophalenInschrijvingRequest,omitempty"`
	KvkNummer string   `xml:"kvkNummer" json:"kvkNummer,omitempty"`
}

type OphalenInschrijvingResponse struct {
	Meldingen *Meldingen `xml:"meldingen" json:"meldingen,omitempty"`
	Product   struct {
		MaatschappelijkeActiviteit *MaatschappelijkeActiviteit `xml:"maatschappelijkeActiviteit" json:"maatschappelijkeActiviteit,omitempty"`
	} `xml:"product"`
	ExtractOriginalXML string `json:"extractOriginal,omitempty"`
}

type Meldingen struct {
	Informatie *Code `xml:"informatie" json:"informatie,omitempty"`
	Fout       *Code `xml:"fout" json:"fout,omitempty"`
}

type MaatschappelijkeActiviteit struct {
	KvkNummer   string `xml:"kvkNummer" json:"kvkNummer,omitempty"`
	Naam        string `xml:"naam,omitempty" json:"naam,omitempty"`
	Registratie struct {
		DatumAanvang string `xml:"datumAanvang,omitempty" json:"datumAanvang,omitempty"`
	} `xml:"registratie,omitempty" json:"registratie,omitempty"`
	BezoekLocatie struct {
		VolledigAdres string `xml:"volledigAdres,omitempty" json:"volledigAdres,omitempty"`
	} `xml:"bezoekLocatie,omitempty" json:"bezoekLocatie,omitempty"`
	Communicatiegegevens struct {
		EmailAdres         []string `xml:"emailAdres,omitempty" json:"emailAdres,omitempty"`
		Communicatienummer []struct {
			Toegangscode string `xml:"toegangscode,omitempty" json:"toegangscode,omitempty"`
			Nummer       string `xml:"nummer,omitempty" json:"nummer,omitempty"`
			Soort        Code   `xml:"soort,omitempty" json:"soort,omitempty"`
		} `xml:"communicatienummer,omitempty" json:"communicatienummer,omitempty"`
	} `xml:"communicatiegegevens,omitempty" json:"communicatiegegevens,omitempty"`
	SbiActiviteit []struct {
		SbiCode           Code `xml:"sbiCode,omitempty" json:"sbiCode,omitempty"`
		IsHoofdactiviteit Code `xml:"isHoofdactiviteit,omitempty" json:"isHoofdactiviteit,omitempty"`
	} `xml:"sbiActiviteit,omitempty" json:"sbiActiviteit,omitempty"`
	ManifesteertZichAls struct {
		Onderneming struct {
			SbiActiviteit []struct {
				SbiCode           Code `xml:"sbiCode,omitempty" json:"sbiCode,omitempty"`
				IsHoofdactiviteit Code `xml:"isHoofdactiviteit,omitempty" json:"isHoofdactiviteit,omitempty"`
			} `xml:"sbiActiviteit,omitempty" json:"sbiActiviteit,omitempty"`
		} `xml:"onderneming,omitempty" json:"onderneming,omitempty"`
	} `xml:"manifesteertZichAls,omitempty" json:"manifesteertZichAls,omitempty"`
	HeeftAlsEigenaar struct {
		Samenwerkingsverband     *Samenwerkingsverband     `xml:"samenwerkingsverband,omitempty" json:"samenwerkingsverband,omitempty"`
		Eenmanszaak              *Eenmanszaak              `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		Rechtspersoon            *Rechtspersoon            `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
		BuitenlandseVennootschap *BuitenlandseVennootschap `xml:"buitenlandseVennootschap,omitempty" json:"buitenlandseVennootschap,omitempty"`
	} `xml:"heeftAlsEigenaar,omitempty" json:"heeftAlsEigenaar,omitempty"`
}

type Eenmanszaak struct {
	PersoonRechtsvorm        string                    `xml:"persoonRechtsvorm,omitempty" json:"persoonRechtsvorm,omitempty"`
	Geslachtsnaam            string                    `xml:"geslachtsnaam,omitempty" json:"geslachtsnaam,omitempty"`
	VoorvoegselGeslachtsnaam string                    `xml:"voorvoegselGeslachtsnaam,omitempty" json:"voorvoegselGeslachtsnaam,omitempty"`
	Voornamen                string                    `xml:"voornamen,omitempty" json:"voornamen,omitempty"`
	Geboortedatum            string                    `xml:"geboortedatum,omitempty" json:"geboortedatum,omitempty"`
	BijzondereRechtstoestand *BijzondereRechtstoestand `xml:"bijzondereRechtstoestand,omitempty" json:"bijzondereRechtstoestand,omitempty"`
	Heeft                    []struct {
		Gemachtigde                          *Gemachtigde                          `xml:"gemachtigde,omitempty" json:"gemachtigde,omitempty"`
		FunctionarisBijzondereRechtstoestand *FunctionarisBijzondereRechtstoestand `xml:"functionarisBijzondereRechtstoestand,omitempty" json:"functionarisBijzondereRechtstoestand,omitempty"`
	} `xml:"heeft,omitempty" json:"heeft,omitempty"`
}

type Samenwerkingsverband struct {
	PersoonRechtsvorm        string                    `xml:"persoonRechtsvorm,omitempty" json:"persoonRechtsvorm,omitempty"`
	BijzondereRechtstoestand *BijzondereRechtstoestand `xml:"bijzondereRechtstoestand,omitempty" json:"bijzondereRechtstoestand,omitempty"`
	Heeft                    []struct {
		Aansprakelijke                       *Aansprakelijke                       `xml:"aansprakelijke,omitempty" json:"aansprakelijke,omitempty"`
		Gemachtigde                          *Gemachtigde                          `xml:"gemachtigde,omitempty" json:"gemachtigde,omitempty"`
		FunctionarisBijzondereRechtstoestand *FunctionarisBijzondereRechtstoestand `xml:"functionarisBijzondereRechtstoestand,omitempty" json:"functionarisBijzondereRechtstoestand,omitempty"`
	} `xml:"heeft,omitempty" json:"heeft,omitempty"`
	Ontbinding Ontbinding `xml:"ontbinding,omitempty" json:"ontbinding,omitempty"`
}

type Rechtspersoon struct {
	PersoonRechtsvorm        string                    `xml:"persoonRechtsvorm,omitempty" json:"persoonRechtsvorm,omitempty"`
	BijzondereRechtstoestand *BijzondereRechtstoestand `xml:"bijzondereRechtstoestand,omitempty" json:"bijzondereRechtstoestand,omitempty"`
	// DatumAkteOprichting      string                    `xml:"datumAkteOprichting,omitempty"`
	Heeft []struct {
		Bestuursfunctie                      *Bestuursfunctie                      `xml:"bestuursfunctie,omitempty" json:"bestuursfunctie,omitempty"`
		Gemachtigde                          *Gemachtigde                          `xml:"gemachtigde,omitempty" json:"gemachtigde,omitempty"`
		OverigeFunctionaris                  *OverigeFunctionaris                  `xml:"overigeFunctionaris,omitempty" json:"overigeFunctionaris,omitempty"`
		FunctionarisBijzondereRechtstoestand *FunctionarisBijzondereRechtstoestand `xml:"functionarisBijzondereRechtstoestand,omitempty" json:"functionarisBijzondereRechtstoestand,omitempty"`
	} `xml:"heeft,omitempty" json:"heeft,omitempty"`
	Ontbinding Ontbinding `xml:"ontbinding,omitempty" json:"ontbinding,omitempty"`
}

type BuitenlandseVennootschap struct {
	PersoonRechtsvorm string `xml:"persoonRechtsvorm,omitempty" json:"persoonRechtsvorm,omitempty"`
	LandVanVestiging  Code   `xml:"landVanVestiging,omitempty" json:"landVanVestiging,omitempty"`
}

type Bestuursfunctie struct {
	Functie      Code `xml:"functie,omitempty" json:"functie,omitempty"`
	Functietitel struct {
		Titel string `xml:"titel,omitempty" json:"titel,omitempty"`
	} `xml:"functietitel,omitempty" json:"functietitel,omitempty"`
	Bevoegdheid Bevoegdheid `xml:"bevoegdheid,omitempty" json:"bevoegdheid,omitempty"`
	Door        struct {
		NatuurlijkPersoon *NatuurlijkPersoon      `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		Rechtspersoon     *RechtspersoonInFunctie `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
	} `xml:"door,omitempty" json:"door,omitempty"`
}

type Aansprakelijke struct {
	Functie     Code        `xml:"functie,omitempty" json:"functie,omitempty"`
	Bevoegdheid Bevoegdheid `xml:"bevoegdheid,omitempty" json:"bevoegdheid,omitempty"`
	Door        struct {
		NatuurlijkPersoon *NatuurlijkPersoon      `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		Rechtspersoon     *RechtspersoonInFunctie `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
	} `xml:"door,omitempty" json:"door,omitempty"`
}

type Gemachtigde struct {
	Functie      Code `xml:"functie,omitempty" json:"functie,omitempty"`
	Functietitel struct {
		Titel string `xml:"titel,omitempty" json:"titel,omitempty"`
	} `xml:"functietitel,omitempty" json:"functietitel,omitempty"`
	Volmacht Volmacht `xml:"volmacht,omitempty" json:"volmacht,omitempty"`
	Door     struct {
		NatuurlijkPersoon *NatuurlijkPersoon      `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		Rechtspersoon     *RechtspersoonInFunctie `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
	} `xml:"door,omitempty" json:"door,omitempty"`
}

type OverigeFunctionaris struct {
	Functie Code `xml:"functie,omitempty" json:"functie,omitempty"`
	Door    struct {
		NatuurlijkPersoon *NatuurlijkPersoon      `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		Rechtspersoon     *RechtspersoonInFunctie `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
	} `xml:"door,omitempty" json:"door,omitempty"`
}

type FunctionarisBijzondereRechtstoestand struct {
	Functie Code `xml:"functie,omitempty" json:"functie,omitempty"`
	Door    struct {
		NaamPersoon struct {
			Naam string `xml:"volledigeNaam,omitempty" json:"volledigeNaam,omitempty"`
		} `xml:"naamPersoon,omitempty" json:"naamPersoon,omitempty"`
	} `xml:"door,omitempty" json:"door,omitempty"`
}

type Bevoegdheid struct {
	Soort            Code `xml:"soort,omitempty" json:"soort,omitempty"`
	BeperkingInEuros struct {
		Waarde string `xml:"waarde,omitempty" json:"waarde,omitempty"`
		Valuta Code   `xml:"code,omitempty" json:"code,omitempty"`
	} `xml:"beperkingInEuros,omitempty" json:"beperkingInEuros,omitempty"`
	OverigeBeperking           Code `xml:"overigeBeperking,omitempty" json:"overigeBeperking,omitempty"`
	IsBevoegdMetAnderePersonen Code `xml:"isBevoegdMetAnderePersonen,omitempty" json:"isBevoegdMetAnderePersonen,omitempty"`
}

type Volmacht struct {
	TypeVolmacht     Code `xml:"typeVolmacht,omitempty" json:"typeVolmacht,omitempty"`
	BeperkteVolmacht struct {
		BeperkingInGeld struct {
			Waarde string `xml:"waarde,omitempty" json:"waarde,omitempty"`
			Valuta Code   `xml:"code,omitempty" json:"code,omitempty"`
		} `xml:"beperkingInGeld,omitempty" json:"beperkingInGeld,omitempty"`
		MagOpgaveHandelsregisterDoen Code   `xml:"magOpgaveHandelsregisterDoen,omitempty" json:"magOpgaveHandelsregisterDoen,omitempty"`
		HeeftOverigeVolmacht         Code   `xml:"heeftOverigeVolmacht,omitempty" json:"heeftOverigeVolmacht,omitempty"`
		OmschrijvingOverigeVolmacht  string `xml:"omschrijvingOverigeVolmacht,omitempty" json:"omschrijvingOverigeVolmacht,omitempty"`
	} `xml:"beperkteVolmacht,omitempty" json:"beperkteVolmacht,omitempty"`
}

type NatuurlijkPersoon struct {
	Geslachtsnaam            string `xml:"geslachtsnaam,omitempty" json:"geslachtsnaam,omitempty"`
	VoorvoegselGeslachtsnaam string `xml:"voorvoegselGeslachtsnaam,omitempty" json:"voorvoegselGeslachtsnaam,omitempty"`
	Voornamen                string `xml:"voornamen,omitempty" json:"voornamen,omitempty"`
	Geboortedatum            string `xml:"geboortedatum,omitempty" json:"geboortedatum,omitempty"`
	BijzondereRechtstoestand struct {
		Soort Code `xml:"soort,omitempty" json:"soort,omitempty"`
	} `xml:"bijzondereRechtstoestand,omitempty" json:"bijzondereRechtstoestand,omitempty"`
}

type RechtspersoonInFunctie struct {
	PersoonRechtsvorm string `xml:"persoonRechtsvorm,omitempty" json:"persoonRechtsvorm,omitempty"`
	VolledigeNaam     string `xml:"volledigeNaam,omitempty" json:"volledigeNaam,omitempty"`
	IsEigenaarVan     struct {
		MaatschappelijkeActiviteit struct {
			KvkNummer string `xml:"kvkNummer,omitempty" json:"kvkNummer,omitempty"`
		} `xml:"maatschappelijkeActiviteit,omitempty" json:"maatschappelijkeActiviteit,omitempty"`
	} `xml:"isEigenaarVan,omitempty" json:"isEigenaarVan,omitempty"`
}

type BijzondereRechtstoestand struct {
	Registratie struct {
		DatumAanvang string `xml:"datumAanvang,omitempty" json:"datumAanvang,omitempty"`
	} `xml:"registratie,omitempty" json:"registratie,omitempty"`
	Soort Code `xml:"soort,omitempty" json:"soort,omitempty"`
}

type Ontbinding struct {
	Registratie struct {
		DatumAanvang string `xml:"datumAanvang,omitempty" json:"datumAanvang,omitempty"`
	} `xml:"registratie,omitempty" json:"registratie,omitempty"`
	Aanleiding Code `xml:"aanleiding,omitempty" json:"aanleiding,omitempty"`
}

type Code struct {
	Code           string `xml:"code,omitempty" json:"code,omitempty"`
	Omschrijving   string `xml:"omschrijving,omitempty" json:"omschrijving,omitempty"`
	ReferentieType string `xml:"referentieType,omitempty" json:"referentieType,omitempty"`
}
