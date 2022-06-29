package models

type BevoegdheidExtract struct {
	KvkNummer          string `json:"kvkNummer,omitempty"`
	Naam               string `json:"naam,omitempty"`
	TypeRechtsvorm     string `json:"typeRechtsvorm,omitempty"`
	Rechtsvorm         string `json:"rechtsvorm,omitempty"`
	Adres              string `json:"adres,omitempty"`
	EmailAdres         string `json:"emailAdres,omitempty"`
	Telefoon           string `json:"telefoon,omitempty"`
	SbiActiviteit      string `json:"sbiActiviteit,omitempty"`
	RegistratieAanvang string `json:"registratieAanvang,omitempty"`

	BijzondereRechtstoestand      string `json:"bijzondereRechtstoestand,omitempty"`
	BijzondereRechtstoestandDatum string `json:"bijzondereRechtstoestandDatum,omitempty"`
	OntbindingDatum               string `json:"ontbindingDatum,omitempty"`
	OntbindingAanleiding          string `json:"ontbindingAanleiding,omitempty"`

	Functionaris *Functionaris `json:"functionaris,omitempty"`

	IsBevoegd bool   `json:"isBevoegd,omitempty"`
	Reason    string `json:"reason,omitempty"`

	ExtractOriginal    *MaatschappelijkeActiviteit `json:"extractOriginal,omitempty"`
	ExtractOriginalXML string                      `json:"extractOriginalXML,omitempty"`

	Paths *Paths `json:"paths,omitempty"`
}

type Functionaris struct {
	DateOfBirth string `json:"dateOfBirth,omitempty"`
	FirstNames  string `json:"firstNames,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	FullName    string `json:"fullName,omitempty"`

	Type         string `json:"type,omitempty"`
	Functie      string `json:"functie,omitempty"`
	Functietitel string `json:"functietitel,omitempty"`

	SoortBevoegdheid           string `json:"soortBevoegdheid,omitempty"`
	BeperkingInEuros           string `json:"beperkingInEuros,omitempty"`
	OverigeBeperking           string `json:"overigeBeperking,omitempty"`
	IsBevoegdMetAnderePersonen string `json:"isBevoegdMetAnderePersonen,omitempty"`

	TypeVolmacht                string `json:"typeVolmacht,omitempty"`
	HeeftOverigeVolmacht        string `json:"heeftOverigeVolmacht,omitempty"`
	OmschrijvingOverigeVolmacht string `json:"omschrijvingOverigeVolmacht,omitempty"`
	BeperkingInGeld             string `json:"beperkingInGeld,omitempty"`

	BijzondereRechtstoestand string `json:"bijzondereRechtstoestand,omitempty"`
}

type Paths struct {
	KvkNummer          string `json:"kvkNummer,omitempty"`
	Naam               string `json:"naam,omitempty"`
	TypeRechtsvorm     string `json:"typeRechtsvorm,omitempty"`
	Rechtsvorm         string `json:"rechtsvorm,omitempty"`
	Adres              string `json:"adres,omitempty"`
	EmailAdres         string `json:"emailAdres,omitempty"`
	Telefoon           string `json:"telefoon,omitempty"`
	SbiActiviteit      string `json:"sbiActiviteit,omitempty"`
	RegistratieAanvang string `json:"registratieAanvang,omitempty"`

	BijzondereRechtstoestand      string       `json:"bijzondereRechtstoestand,omitempty"`
	BijzondereRechtstoestandDatum string       `json:"bijzondereRechtstoestandDatum,omitempty"`
	OntbindingDatum               string       `json:"ontbindingDatum,omitempty"`
	OntbindingAanleiding          string       `json:"ontbindingAanleiding,omitempty"`
	Functionaris                  Functionaris `json:"functionaris"`
}
