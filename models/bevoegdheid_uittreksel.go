package models

type IdentityNP struct {
	Geslachtsnaam            string `json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam string `json:"voorvoegselGeslachtsnaam"`
	Voornamen                string `json:"voornamen"`
	Geboortedatum            string `json:"geboortedatum"`
}

type Interpretatie struct {
	HeeftBeperking bool   `json:"heeftBeperking"`
	IsBevoegd      string `json:"isBevoegd"`
	Reden          string `json:"reden"`
}

type BevoegdheidUittreksel struct {
	KvkNummer          string `json:"kvkNummer"`
	Rsin               string `json:"rsin"`
	Naam               string `json:"naam"`
	TypeEigenaar       string `json:"typeEigenaar"`
	PersoonRechtsvorm  string `json:"persoonRechtsvorm"`
	Adres              string `json:"adres"`
	EmailAdres         string `json:"emailAdres"`
	Telefoon           string `json:"telefoon"`
	SbiActiviteit      string `json:"sbiActiviteit"`
	RegistratieAanvang string `json:"registratieAanvang"`

	DatumUitschrijving         string `json:"datumUitschrijving"`
	RegistratieEinde           string `json:"registratieEinde"`
	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand string `json:"buitenlandseRechtstoestand"`
	// Handlichting               string `json:"handlichting"`

	Functionaris        *Functionaris  `json:"functionaris,omitempty"`
	AlleFunctionarissen []Functionaris `json:"functionarissen"`
}

type BevoegdheidResponse struct {
	BevoegdheidUittreksel *BevoegdheidUittreksel `json:"bevoegdheidUittreksel"`
	Interpretatie         *Interpretatie         `json:"interpretatie"`

	ExtractOriginal    *MaatschappelijkeActiviteit `json:"extractOriginal"`
	ExtractOriginalXML string                      `json:"extractOriginalXML"`

	Paths *Paths `json:"paths"`
}

type Functionaris struct {
	Geslachtsnaam            string `json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam string `json:"voorvoegselGeslachtsnaam"`
	Voornamen                string `json:"voornamen"`
	Geboortedatum            string `json:"geboortedatum"`
	Overlijdensdatum         string `json:"overlijdensdatum"`
	VolledigeNaam            string `json:"volledigeNaam"`

	TypeFunctionaris string `json:"typeFunctionaris"`
	Functie          string `json:"functie"`
	Functietitel     string `json:"functietitel"`

	SoortBevoegdheid            string `json:"soortBevoegdheid"`
	BeperkingInEurosBevoegdheid string `json:"beperkingInEurosBevoegdheid"`
	OverigeBeperkingBevoegdheid string `json:"overigeBeperkingBevoegdheid"`
	IsBevoegdMetAnderePersonen  string `json:"isBevoegdMetAnderePersonen"`

	TypeVolmacht                 string `json:"typeVolmacht"`
	BeperkingInGeldVolmacht      string `json:"beperkingInGeldVolmacht"`
	BeperkingInHandelingVolmacht string `json:"beperkingInHandelingVolmacht"`
	HeeftOverigeVolmacht         string `json:"heeftOverigeVolmacht"`
	OmschrijvingOverigeVolmacht  string `json:"omschrijvingOverigeVolmacht"`
	MagOpgaveHandelsregisterDoen string `json:"magOpgaveHandelsregisterDoen"`

	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestandFunctionaris"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandelingFunctionaris"`
	SchorsingAanvang           string `json:"schorsingAanvang"`
	SchorsingEinde             string `json:"schorsingEinde"`
	Handlichting               string `json:"handlichting"`
}

type Paths struct {
	KvkNummer          string `json:"kvkNummer"`
	Rsin               string `json:"rsin"`
	Naam               string `json:"naam"`
	TypeEigenaar       string `json:"typeEigenaar"`
	PersoonRechtsvorm  string `json:"persoonRechtsvorm"`
	Adres              string `json:"adres"`
	EmailAdres         string `json:"emailAdres"`
	Telefoon           string `json:"telefoon"`
	SbiActiviteit      string `json:"sbiActiviteit"`
	RegistratieAanvang string `json:"registratieAanvang"`

	DatumUitschrijving         string `json:"datumUitschrijving"`
	RegistratieEinde           string `json:"registratieEinde"`
	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand string `json:"buitenlandseRechtstoestand"`
	Handlichting               string `json:"handlichting"`

	Functionaris Functionaris `json:"functionaris,omitempty"`
}
