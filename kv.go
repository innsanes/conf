package conf

type kv[V interface{}] struct {
	keyValue map[string]V
}

func newKV[V interface{}]() *kv[V] {
	return &kv[V]{
		keyValue: make(map[string]V),
	}
}

func (k *kv[V]) Get(str string) (V, bool) {
	v, has := k.keyValue[str]
	return v, has
}

func (k *kv[V]) Set(str string, v V) {
	k.keyValue[str] = v
}

func (k *kv[V]) Range(f func(key string, value V) bool) {
	for key, value := range k.keyValue {
		if !f(key, value) {
			break
		}
	}
}
