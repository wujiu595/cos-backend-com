package proto

type Data map[string]interface{}

func (d Data) Copy() Data {
	newData := make(Data)
	for key, value := range d {
		newData[key] = value
	}
	return newData
}
