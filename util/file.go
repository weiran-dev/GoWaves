package util

import "os"

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckOrCreateDir(path string) (bool,error) {
	isExist,err := PathExists(path)
	if err != nil {
		return false, err
	}
	if isExist {
		return true,nil
	}
	err = os.Mkdir(path,os.ModePerm)
	if err != nil {
		return false, err
	}
	return true,nil
}

func CheckOrCreateFile(path string) (bool,error) {
	isExist,err := PathExists(path)
	if err != nil {
		return false, err
	}
	if isExist {
		return true,nil
	}
	_,err = os.Create(path)
	if err != nil {
		return false, err
	}
	return true,nil
}