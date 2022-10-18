package models

import "encoding/xml"

type OphalenInschrijvingRequest struct {
	XMLName   xml.Name `xml:"ophalenInschrijvingRequest" json:"ophalenInschrijvingRequest"`
	KvkNummer string   `xml:"kvkNummer" json:"kvkNummer"`
}

type OphalenInschrijvingResponse struct {
	Meldingen *Meldingen `xml:"meldingen" json:"meldingen"`
	Product   struct {
		MaatschappelijkeActiviteit *MaatschappelijkeActiviteit `xml:"maatschappelijkeActiviteit" json:"maatschappelijkeActiviteit"`
	} `xml:"product"`
	ExtractOriginalXML string `json:"extractOriginal"`
}

type Meldingen struct {
	Informatie *Enumeratie `xml:"informatie" json:"informatie"`
	Fout       *Enumeratie `xml:"fout" json:"fout"`
}

type MaatschappelijkeActiviteit struct {
	KvkNummer     string      `xml:"kvkNummer" json:"kvkNummer"`
	Naam          string      `xml:"naam" json:"naam"`
	Registratie   Registratie `xml:"registratie" json:"registratie"`
	BezoekLocatie struct {
		VolledigAdres string `xml:"volledigAdres" json:"volledigAdres"`
	} `xml:"bezoekLocatie" json:"bezoekLocatie"`
	Communicatiegegevens struct {
		EmailAdres         []string `xml:"emailAdres" json:"emailAdres"`
		Communicatienummer []struct {
			Toegangscode string     `xml:"toegangscode" json:"toegangscode"`
			Nummer       string     `xml:"nummer" json:"nummer"`
			Soort        Enumeratie `xml:"soort" json:"soort"`
		} `xml:"communicatienummer" json:"communicatienummer"`
	} `xml:"communicatiegegevens" json:"communicatiegegevens"`
	SbiActiviteit []struct {
		SbiCode           Enumeratie `xml:"sbiCode" json:"sbiCode"`
		IsHoofdactiviteit Enumeratie `xml:"isHoofdactiviteit" json:"isHoofdactiviteit"`
	} `xml:"sbiActiviteit" json:"sbiActiviteit"`
	ManifesteertZichAls struct {
		Onderneming struct {
			SbiActiviteit []struct {
				SbiCode           Enumeratie `xml:"sbiCode" json:"sbiCode"`
				IsHoofdactiviteit Enumeratie `xml:"isHoofdactiviteit" json:"isHoofdactiviteit"`
			} `xml:"sbiActiviteit" json:"sbiActiviteit"`
		} `xml:"onderneming" json:"onderneming"`
	} `xml:"manifesteertZichAls" json:"manifesteertZichAls"`
	HeeftAlsEigenaar struct {
		NaamPersoon                     *NietNatuurlijkPersoon `xml:"naamPersoon,omitempty" json:"naamPersoon,omitempty"`
		Eenmanszaak                     *Eenmanszaak           `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
		BuitenlandseVennootschap        *NietNatuurlijkPersoon `xml:"buitenlandseVennootschap,omitempty" json:"buitenlandseVennootschap,omitempty"`
		EenmanszaakMetMeerdereEigenaren *NietNatuurlijkPersoon `xml:"eenmanszaakMetMeerdereEigenaren,omitempty" json:"eenmanszaakMetMeerdereEigenaren,omitempty"`
		Rechtspersoon                   *NietNatuurlijkPersoon `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
		RechtspersoonInOprichting       *NietNatuurlijkPersoon `xml:"rechtspersoonInOprichting,omitempty" json:"rechtspersoonInOprichting,omitempty"`
		Samenwerkingsverband            *NietNatuurlijkPersoon `xml:"samenwerkingsverband,omitempty" json:"samenwerkingsverband,omitempty"`
		AfgeslotenMoeder                *NietNatuurlijkPersoon `xml:"afgeslotenMoeder,omitempty" json:"afgeslotenMoeder,omitempty"`
	} `xml:"heeftAlsEigenaar" json:"heeftAlsEigenaar"`
}

type Eenmanszaak struct { // in KVK productstore this is 'natuurlijkPersoon' but that conflicts with NatuurlijkPersoon
	Registratie                Registratie                `xml:"registratie" json:"registratie"`
	PersoonRechtsvorm          string                     `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
	Geslachtsnaam              string                     `xml:"geslachtsnaam" json:"geslachtsnaam"`
	Voornamen                  string                     `xml:"voornamen" json:"voornamen"`
	VoorvoegselGeslachtsnaam   string                     `xml:"voorvoegselGeslachtsnaam" json:"voorvoegselGeslachtsnaam"`
	Geboortedatum              string                     `xml:"geboortedatum" json:"geboortedatum"`
	Overlijdensdatum           string                     `xml:"overlijdensdatum" json:"overlijdensdatum"`
	VolledigeNaam              string                     `xml:"volledigeNaam" json:"volledigeNaam"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
	Handlichting               Handlichting               `xml:"handlichting" json:"handlichting"`
	Heeft                      []Functievervulling        `xml:"heeft" json:"heeft"`
}

type NietNatuurlijkPersoon struct {
	Registratie                Registratie                `xml:"registratie" json:"registratie"`
	DatumUitschrijving         string                     `xml:"datumUitschrijving" json:"datumUitschrijving"`
	PersoonRechtsvorm          string                     `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
	Ontbinding                 Ontbinding                 `xml:"ontbinding" json:"ontbinding"`
	Heeft                      []Functievervulling        `xml:"heeft" json:"heeft"`
	Rsin                       string                     `xml:"rsin" json:"rsin"`
	// LandVanVestiging           Enumeratie                  `xml:"landVanVestiging" json:"landVanVestiging"`
}

// type NaamPersoon struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

// type BuitenlandseVennootschap struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// 	// LandVanVestiging           Enumeratie                  `xml:"landVanVestiging" json:"landVanVestiging"`
// }

// type EenmanszaakMetMeerdereEigenaren struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

// type Rechtspersoon struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

// type RechtspersoonInOprichting struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

// type Samenwerkingsverband struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

// type AfgeslotenMoeder struct {
// 	Registratie                Registratie                 `xml:"registratie" json:"registratie"`
// 	DatumUitschrijving         string                      `xml:"datumUitschrijving" json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `xml:"buitenlandseRechtstoestand" json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `xml:"ontbinding" json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `xml:"heeft" json:"heeft"`
// }

type Functievervulling struct {
	Aansprakelijke                       *FunctionarisOfGemachtigde `xml:"aansprakelijke,omitempty" json:"aansprakelijke,omitempty"`
	Bestuursfunctie                      *FunctionarisOfGemachtigde `xml:"bestuursfunctie,omitempty" json:"bestuursfunctie,omitempty"`
	FunctionarisBijzondereRechtstoestand *FunctionarisOfGemachtigde `xml:"functionarisBijzondereRechtstoestand,omitempty" json:"functionarisBijzondereRechtstoestand,omitempty"`
	Gemachtigde                          *FunctionarisOfGemachtigde `xml:"gemachtigde,omitempty" json:"gemachtigde,omitempty"`
	OverigeFunctionaris                  *FunctionarisOfGemachtigde `xml:"overigeFunctionaris,omitempty" json:"overigeFunctionaris,omitempty"`
	PubliekrechtelijkeFunctionaris       *FunctionarisOfGemachtigde `xml:"publiekrechtelijkeFunctionaris,omitempty" json:"publiekrechtelijkeFunctionaris,omitempty"`
}

type Door struct {
	NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon,omitempty" json:"natuurlijkPersoon,omitempty"`
	Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon,omitempty" json:"rechtspersoon,omitempty"`
}

type FunctionarisOfGemachtigde struct {
	Functie      Enumeratie `xml:"functie" json:"functie"`
	Functietitel struct {
		Titel string `xml:"titel" json:"titel"`
	} `xml:"functietitel" json:"functietitel"`
	Bevoegdheid  Bevoegdheid  `xml:"bevoegdheid" json:"bevoegdheid"`
	Volmacht     Volmacht     `xml:"volmacht" json:"volmacht"`
	Handlichting Handlichting `xml:"handlichting" json:"handlichting"`
	Schorsing    Schorsing    `xml:"schorsing" json:"schorsing"`
	Door         Door         `xml:"door" json:"door"`
}

// type Gemachtigde struct {
// 	Functie      Enumeratie `xml:"functie" json:"functie"`
// 	Functietitel struct {
// 		Titel string `xml:"titel" json:"titel"`
// 	} `xml:"functietitel" json:"functietitel"`
// 	Volmacht  Volmacht  `xml:"volmacht" json:"volmacht"`
// 	Schorsing Schorsing `xml:"schorsing" json:"schorsing"`
// 	Door      Door      `xml:"door" json:"door"`
// }

// type Aansprakelijke struct {
// 	Functie      Enumeratie    `xml:"functie" json:"functie"`
// 	Bevoegdheid  Bevoegdheid   `xml:"bevoegdheid" json:"bevoegdheid"`
// 	Handlichting Handlichting `xml:"handlichting" json:"handlichting"`
// 	Schorsing    Schorsing    `xml:"schorsing" json:"schorsing"`
// 	Door         struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon" json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon" json:"rechtspersoon"`
// 	} `xml:"door" json:"door"`
// }

// type Bestuursfunctie struct {
// 	Functie      Enumeratie `xml:"functie" json:"functie"`
// 	Functietitel struct {
// 		Titel string `xml:"titel" json:"titel"`
// 	} `xml:"functietitel" json:"functietitel"`
// 	Bevoegdheid Bevoegdheid `xml:"bevoegdheid" json:"bevoegdheid"`
// 	Schorsing   Schorsing  `xml:"schorsing" json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon" json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon" json:"rechtspersoon"`
// 	} `xml:"door" json:"door"`
// }

// type FunctionarisBijzondereRechtstoestand struct {
// 	Functie   Enumeratie `xml:"functie" json:"functie"`
// 	Schorsing Schorsing `xml:"schorsing" json:"schorsing"`
// 	Door      struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon" json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon" json:"rechtspersoon"`
// 	} `xml:"door" json:"door"`
// }

// type OverigeFunctionaris struct {
// 	Functie     Enumeratie  `xml:"functie" json:"functie"`
// 	Bevoegdheid Bevoegdheid `xml:"bevoegdheid" json:"bevoegdheid"`
// 	Schorsing   Schorsing  `xml:"schorsing" json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon" json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon" json:"rechtspersoon"`
// 	} `xml:"door" json:"door"`
// }

// type PubliekrechtelijkeFunctionaris struct {
// 	Functie     Enumeratie  `xml:"functie" json:"functie"`
// 	Bevoegdheid Bevoegdheid `xml:"bevoegdheid" json:"bevoegdheid"`
// 	Schorsing   Schorsing  `xml:"schorsing" json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `xml:"natuurlijkPersoon" json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `xml:"rechtspersoon" json:"rechtspersoon"`
// 	} `xml:"door" json:"door"`
// }

type Bevoegdheid struct {
	Soort            Enumeratie `xml:"soort" json:"soort"`
	BeperkingInEuros struct {
		Waarde string     `xml:"waarde" json:"waarde"`
		Valuta Enumeratie `xml:"valuta" json:"valuta"`
	} `xml:"beperkingInEuros" json:"beperkingInEuros"`
	OverigeBeperking           Enumeratie `xml:"overigeBeperking" json:"overigeBeperking"`
	IsBevoegdMetAnderePersonen Enumeratie `xml:"isBevoegdMetAnderePersonen" json:"isBevoegdMetAnderePersonen"`
}

type Volmacht struct {
	TypeVolmacht     Enumeratie `xml:"typeVolmacht" json:"typeVolmacht"`
	BeperkteVolmacht struct {
		BeperkingInHandeling struct {
			SoortHandeling Enumeratie `xml:"soortHandeling" json:"soortHandeling"`
		} `xml:"beperkingInHandeling" json:"beperkingInHandeling"`
		BeperkingInGeld struct {
			Waarde string     `xml:"waarde" json:"waarde"`
			Valuta Enumeratie `xml:"valuta" json:"valuta"`
		} `xml:"beperkingInGeld" json:"beperkingInGeld"`
		MagOpgaveHandelsregisterDoen Enumeratie `xml:"magOpgaveHandelsregisterDoen" json:"magOpgaveHandelsregisterDoen"`
		HeeftOverigeVolmacht         Enumeratie `xml:"heeftOverigeVolmacht" json:"heeftOverigeVolmacht"`
		OmschrijvingOverigeVolmacht  string     `xml:"omschrijvingOverigeVolmacht" json:"omschrijvingOverigeVolmacht"`
	} `xml:"beperkteVolmacht" json:"beperkteVolmacht"`
}

type NatuurlijkPersoon struct {
	Geslachtsnaam              string                     `xml:"geslachtsnaam" json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam   string                     `xml:"voorvoegselGeslachtsnaam" json:"voorvoegselGeslachtsnaam"`
	Voornamen                  string                     `xml:"voornamen" json:"voornamen"`
	Geboortedatum              string                     `xml:"geboortedatum" json:"geboortedatum"`
	Overlijdensdatum           string                     `xml:"overlijdensdatum" json:"overlijdensdatum"`
	VolledigeNaam              string                     `xml:"volledigeNaam" json:"volledigeNaam"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `xml:"bijzondereRechtstoestand" json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `xml:"beperkingInRechtshandeling" json:"beperkingInRechtshandeling"`
}

type RechtspersoonAlsFunctionaris struct {
	PersoonRechtsvorm string `xml:"persoonRechtsvorm" json:"persoonRechtsvorm"`
	VolledigeNaam     string `xml:"volledigeNaam" json:"volledigeNaam"`
	IsEigenaarVan     struct {
		MaatschappelijkeActiviteit struct {
			KvkNummer string `xml:"kvkNummer" json:"kvkNummer"`
		} `xml:"maatschappelijkeActiviteit" json:"maatschappelijkeActiviteit"`
	} `xml:"isEigenaarVan" json:"isEigenaarVan"`
}

type BijzondereRechtstoestand struct {
	Registratie Registratie `xml:"registratie" json:"registratie"`
	Soort       Enumeratie  `xml:"soort" json:"soort"`
}

type BeperkingInRechtshandeling struct {
	Registratie Registratie `xml:"registratie" json:"registratie"`
	Soort       Enumeratie  `xml:"soort" json:"soort"`
}

type BuitenlandseRechtstoestand struct {
	Registratie  Registratie `xml:"registratie" json:"registratie"`
	Beschrijving string      `xml:"beschrijving" json:"beschrijving"`
}

type Handlichting struct {
	Registratie Registratie `xml:"registratie" json:"registratie"`
	IsVerleend  Enumeratie  `xml:"isVerleend" json:"isVerleend"`
}

type Ontbinding struct {
	Registratie Registratie `xml:"registratie" json:"registratie"`
	Aanleiding  Enumeratie  `xml:"aanleiding" json:"aanleiding"`
	Liquidatie  struct {
		Registratie Registratie `xml:"registratie" json:"registratie"`
	} `xml:"liquidatie" json:"liquidatie"`
}

type Schorsing struct {
	Registratie Registratie `xml:"registratie" json:"registratie"`
}

type Registratie struct {
	RegistratieTijdstip string `xml:"registratieTijdstip" json:"registratieTijdstip"`
	DatumAanvang        string `xml:"datumAanvang" json:"datumAanvang"`
	DatumEinde          string `xml:"datumEinde" json:"datumEinde"`
}

type Enumeratie struct {
	Code           string `xml:"code" json:"code"`
	Omschrijving   string `xml:"omschrijving" json:"omschrijving"`
	ReferentieType string `xml:"referentieType" json:"referentieType"`
}
