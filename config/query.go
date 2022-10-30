package config

import (
	"github.com/beevik/etree"
	pq "github.com/pasiol/gopq"
)

// Applicant struct
type Account struct {
	ID                 string
	NoLDAP             string
	StudentType        string
	UserGroup          string
	NewTypeUserAccount string
}

func MapData(data []string) Account {
	a := Account{}

	a.ID = data[0]
	a.NoLDAP = "Ei"
	a.StudentType = "N" // TODO: integer
	a.NewTypeUserAccount = "Ei"
	a.UserGroup = "nnnnnnnnn" // TODO: string
	return a
}

// Accounts query
func Accounts() pq.PrimusQuery {
	pq := pq.PrimusQuery{}
	pq.Charset = "UTF-8"
	pq.Database = "opphenk"
	pq.Sort = ""
	pq.Search = ""
	pq.Data = "#DATA{V1}"
	pq.Footer = ""

	return pq
}

// ArchieveXML generator
func UpdateAccountXML(a Account) (string, error) {
	updateDoc := etree.NewDocument()
	updateDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	primusquery := updateDoc.CreateElement("PRIMUSQUERY_IMPORT")
	primusquery.CreateAttr("ARCHIVEMODE", "0")
	primusquery.CreateAttr("CREATEIFNOTFOUND", "0")
	identity := updateDoc.CreateElement("IDENTITY")
	identity.CreateText("service-update-accounts")
	card := updateDoc.CreateElement("CARD")
	card.CreateAttr("FIND", a.ID)
	noLDAP := card.CreateElement("EILDAP")
	noLDAP.CreateText(a.NoLDAP)
	studentType := card.CreateElement("OPISKELIJALAJI")
	studentType.CreateText(a.StudentType)
	userGroup := card.CreateElement("KÄYTTÄJÄRYHMÄ")
	userGroup.CreateAttr("CMD", "MODIFY")
	userGroup.CreateAttr("LINE", "1")
	userGroup.CreateText(a.UserGroup)
	newUserAccount := card.CreateElement("UUSITUNNUS")
	newUserAccount.CreateText("0")
	newTypeUserAccount := card.CreateElement("UUSITUNNUSKAYTOSSA")
	newTypeUserAccount.CreateText(a.NewTypeUserAccount)
	loginDate := card.CreateElement("KIRJAUTUMISPÄIVÄ")
	loginDate.CreateText("")
	loginTime := card.CreateElement("KIRJAUDUTTU_VIIMEKSI")
	loginTime.CreateText("")

	updateDoc.Indent(2)
	xmlAsString, _ := updateDoc.WriteToString()
	filename, err := pq.CreateTMPFile(pq.StringWithCharset(128)+".xml", xmlAsString)
	if err != nil {
		return "", err
	}
	return filename, nil
}
