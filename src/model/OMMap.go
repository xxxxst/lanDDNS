
package model

// on-to-many map
type OMMap struct {
	mapKeyToVal map[interface{}] interface{};
	mapValToKey map[interface{}] *HashSet;
}

func CreateOMMap() *OMMap {
	md := &OMMap{};
	md.mapKeyToVal = make(map[interface{}] interface{});
	md.mapValToKey = make(map[interface{}] *HashSet);
	return md;
}

func (c *OMMap) Add(key interface{}, val interface{}) {
	oldVal, ok := c.mapKeyToVal[key];
	if(ok){
		hashVal, ok2 := c.mapValToKey[oldVal];
		if(ok2){
			hashVal.Remove(key);
		}
	}

	c.mapKeyToVal[key] = val;

	hashVal,ok3 := c.mapValToKey[val];
	if(!ok3){
		hashVal = CreateHashSet();
		c.mapValToKey[val] = hashVal;
	}

	hashVal.Add(key);
}

func (c *OMMap) RemoveKey(key interface{}) {
	val, ok := c.mapKeyToVal[key];
	if(!ok){
		return;
	}

	hashVal,ok2 := c.mapValToKey[val];
	
	delete(c.mapKeyToVal, key);
	if(!ok2){
		return;
	}

	hashVal.Remove(key);
	if(hashVal.Count() <= 0){
		delete(c.mapValToKey, val);
	}
}

func (c *OMMap) RemoveVal(val interface{}) {
	hashVal, ok := c.mapValToKey[val];
	if(!ok){
		return;
	}

	hashVal.Erg(func(key interface{}){
		delete(c.mapKeyToVal, key);
	});
}

func (c *OMMap) GetVal(key interface{}) interface{} {
	val, ok := c.mapKeyToVal[key];
	if(!ok){
		return nil;
	}
	return val;
}

func (c *OMMap) KeyExist(key interface{}) bool {
	_, ok := c.mapKeyToVal[key];
	return ok;
}

func (c *OMMap) GetKey(val interface{}) *HashSet {
	hashVal, ok := c.mapValToKey[val];
	if(!ok){
		return nil;
	}
	return hashVal;
}
