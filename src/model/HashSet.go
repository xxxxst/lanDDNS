
package model

type HashSet struct {
	data map[interface{}] bool;
}

func CreateHashSet() *HashSet {
	md := &HashSet{};
	md.data = make(map[interface{}] bool);
	return md;
}

// type HashSetString struct {
// 	data map[string] bool;
// }

func (c *HashSet) Add(key interface{}) {
	c.data[key] = true;
}

func (c *HashSet) Exist(key interface{}) bool {
	_,ok := c.data[key];
	return ok;
}

func (c *HashSet) Remove(key interface{}) {
	delete(c.data, key);
}

func (c *HashSet) Erg(fun func(interface{})) {
	for key, _ := range c.data {
		fun(key);
	}
}

func (c *HashSet) Count() int {
	return len(c.data);
}