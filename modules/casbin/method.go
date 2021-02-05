package casbin

import "strings"

//====================================
// func
//====================================

// DomainMatch domain
func DomainMatch(key1 string, key2 string) bool {
	if key2[:1] == "." {
		return strings.HasSuffix(key1, key2)
	}
	i := strings.Index(key2, "*") + 1
	if i == 0 {
		return key1 == key2
	}
	l := len(key2)
	if i == l {
		return true
	}
	if li := len(key1) - (l - i); li > 0 {
		// 截取key1可用部分
		return key1[li:] == key2[i:]
	}
	return key1 == key2[i:]
}

// DomainMatchFunc domain
func DomainMatchFunc(args ...interface{}) (interface{}, error) {
	domain1 := args[0].(string)
	domain2 := args[1].(string)
	if domain2 == "" || domain2 == "*" {
		return true, nil
	}
	return DomainMatch(domain1, domain2), nil
}

// AudienceMatchFunc domain
func AudienceMatchFunc(args ...interface{}) (interface{}, error) {
	domain1 := args[0].(string)
	domain2 := args[1].(string)
	audience := args[2].(string)
	if domain2 == "" || domain2 == "*" {
		return true, nil
	}
	if domain2 == "jwt" {
		return DomainMatch(domain1, audience), nil
	}
	return DomainMatch(domain1, domain2), nil
}

// MethodMatchFunc action
func MethodMatchFunc(args ...interface{}) (interface{}, error) {
	method := args[0].(string)
	methods := args[1].(string)
	if methods == "" || methods == "*" {
		return true, nil
	}
	return strings.Contains(methods, method), nil
}
