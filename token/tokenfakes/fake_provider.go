// This file was generated by counterfeiter
package tokenfakes

import (
	"sync"

	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/token"
	"golang.org/x/oauth2"
)

type FakeProvider struct {
	RequestTokenStub        func(code string) (*oauth2.Token, error)
	requestTokenMutex       sync.RWMutex
	requestTokenArgsForCall []struct {
		code string
	}
	requestTokenReturns struct {
		result1 *oauth2.Token
		result2 error
	}
	LoginURLStub        func(state string) string
	loginURLMutex       sync.RWMutex
	loginURLArgsForCall []struct {
		state string
	}
	loginURLReturns struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProvider) RequestToken(code string) (*oauth2.Token, error) {
	fake.requestTokenMutex.Lock()
	fake.requestTokenArgsForCall = append(fake.requestTokenArgsForCall, struct {
		code string
	}{code})
	fake.recordInvocation("RequestToken", []interface{}{code})
	fake.requestTokenMutex.Unlock()
	if fake.RequestTokenStub != nil {
		return fake.RequestTokenStub(code)
	} else {
		return fake.requestTokenReturns.result1, fake.requestTokenReturns.result2
	}
}

func (fake *FakeProvider) RequestTokenCallCount() int {
	fake.requestTokenMutex.RLock()
	defer fake.requestTokenMutex.RUnlock()
	return len(fake.requestTokenArgsForCall)
}

func (fake *FakeProvider) RequestTokenArgsForCall(i int) string {
	fake.requestTokenMutex.RLock()
	defer fake.requestTokenMutex.RUnlock()
	return fake.requestTokenArgsForCall[i].code
}

func (fake *FakeProvider) RequestTokenReturns(result1 *oauth2.Token, result2 error) {
	fake.RequestTokenStub = nil
	fake.requestTokenReturns = struct {
		result1 *oauth2.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeProvider) LoginURL(state string) string {
	fake.loginURLMutex.Lock()
	fake.loginURLArgsForCall = append(fake.loginURLArgsForCall, struct {
		state string
	}{state})
	fake.recordInvocation("LoginURL", []interface{}{state})
	fake.loginURLMutex.Unlock()
	if fake.LoginURLStub != nil {
		return fake.LoginURLStub(state)
	} else {
		return fake.loginURLReturns.result1
	}
}

func (fake *FakeProvider) LoginURLCallCount() int {
	fake.loginURLMutex.RLock()
	defer fake.loginURLMutex.RUnlock()
	return len(fake.loginURLArgsForCall)
}

func (fake *FakeProvider) LoginURLArgsForCall(i int) string {
	fake.loginURLMutex.RLock()
	defer fake.loginURLMutex.RUnlock()
	return fake.loginURLArgsForCall[i].state
}

func (fake *FakeProvider) LoginURLReturns(result1 string) {
	fake.LoginURLStub = nil
	fake.loginURLReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestTokenMutex.RLock()
	defer fake.requestTokenMutex.RUnlock()
	fake.loginURLMutex.RLock()
	defer fake.loginURLMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ token.Provider = new(FakeProvider)
