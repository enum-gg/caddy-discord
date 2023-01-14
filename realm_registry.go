package caddydiscord

type RealmRegistry []*Realm

func (r RealmRegistry) ByName(name string) *Realm {
	for _, realm := range r {
		if realm.Ref == name {
			return realm
		}
	}

	return nil
}
