package localstorage

import (
	"clave/objects"
	"errors"
	"log"
	"strings"
)

const ListOfSecrets = "ListOfSecrets"

func (kvs *PersistentStore) AddTotpSecretObject(issuer, secret string) error {
	object := objects.CreateNewTotpSecretObject(issuer, secret)

	listOfSecrets := []objects.TotpSecretObject{}

	err := kvs.Get(ListOfSecrets, &listOfSecrets)

	if err != nil {

		if errors.Is(err, ErrNotFound) {
			listOfSecrets = []objects.TotpSecretObject{}
		} else {
			log.Println("[AddTotpSecretObject][Error] Error getting list of secrets ", err)
			return err
		}

	}

	listOfSecrets = append(listOfSecrets, object)
	err = kvs.SetValue(ListOfSecrets, listOfSecrets)

	if err != nil {
		log.Println("Error saving list of secrets ", err)
		return err
	}

	return nil

}
func (kvs *PersistentStore) GetListOfTotpSecretObjects() objects.TotpSecretObjectList {

	listOfSecrets := []objects.TotpSecretObject{}

	err := kvs.Get(ListOfSecrets, &listOfSecrets)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return listOfSecrets
		} else {
			log.Println("Error getting list of secrets ", err)
			return listOfSecrets
		}
	}

	return listOfSecrets

}

func (kvs *PersistentStore) DeleteTotpSecretObject(secretId string) error {

	listOfSecrets := []objects.TotpSecretObject{}

	err := kvs.Get(ListOfSecrets, &listOfSecrets)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		} else {
			log.Println("Error getting list of secrets ", err)
			return err
		}
	}

	for i, _obj := range listOfSecrets {
		if strings.EqualFold(_obj.Id, secretId) {
			listOfSecrets = append(listOfSecrets[:i], listOfSecrets[i+1:]...)
			break
		}
	}

	err = kvs.SetValue(ListOfSecrets, listOfSecrets)

	if err != nil {
		log.Println("Error saving list of secrets ", err)
		return err
	}

	return nil

}

func (kvs *PersistentStore) CheckIfIssuerOrSecretExists(issuer, secret string) (bool, objects.TotpSecretObject) {

	listOfSecrets := []objects.TotpSecretObject{}

	err := kvs.Get(ListOfSecrets, &listOfSecrets)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, objects.TotpSecretObject{}
		} else {
			log.Println("Error getting list of secrets ", err)
			return false, objects.TotpSecretObject{}
		}
	}

	for _, _obj := range listOfSecrets {
		if strings.EqualFold(_obj.Issuer, issuer) || strings.EqualFold(_obj.Secret, secret) {
			log.Println("Issuer or secret already exists")
			return true, _obj
		}
	}

	return false, objects.TotpSecretObject{}
}
