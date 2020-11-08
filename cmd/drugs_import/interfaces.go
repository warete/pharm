package drugs_import

type DrugsImporter interface {
	Run() error
	getDataFromService(...interface{}) (interface{}, error)
	processData(string) interface{}
}
