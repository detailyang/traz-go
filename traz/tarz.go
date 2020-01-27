// Package traz transforms differernt data source to the go struct which is inspired from
package traz

type Transformer interface {
	Apply(s interface{}) error
}

type TransformerFunc func(s interface{}) error

func (tf TransformerFunc) Apply(s interface{}) error {
	return tf(s)
}

type MultiTraz struct {
	transformers []Transformer
}

func NewMultiTraz(transformers ...Transformer) *MultiTraz {
	return &MultiTraz{
		transformers: transformers,
	}
}

func (mt *MultiTraz) Append(trazs ...Transformer) {
	mt.transformers = append(mt.transformers, trazs...)
}

func (mt *MultiTraz) Apply(s interface{}) error {
	for i := range mt.transformers {
		if err := mt.transformers[i].Apply(s); err != nil {
			return err
		}
	}
	return nil
}
