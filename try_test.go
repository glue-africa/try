package try_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/glue-africa/try"
)

func TestTryExample(t *testing.T) {
	try.MaxRetries = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		value, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = value
}

func TestTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}
	var value string
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5 // try 5 times
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		value, err = SomeFunction()
		return
	})
	_ = value
	if err != nil {
		t.Logf("error: %v", err)
	}
}

func TestTryDoSuccessful(t *testing.T) {
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
}

func TestTryDoFailed(t *testing.T) {
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	if err != theErr {
		t.Errorf("Expected error %v, got %v", theErr, err)
	}
	if callCount != 5 {
		t.Errorf("Expected callCount to be 5, got %d", callCount)
	}
}

func TestTryPanics(t *testing.T) {
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		callCount++
		if attempt > 2 {
			panic("I don't like three")
		}
		err = theErr
		return
	})
	expectedErrMsg := "panic: I don't like three"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
	}
	if callCount != 5 {
		t.Errorf("Expected callCount to be 5, got %d", callCount)
	}
}

func TestRetryLimit(t *testing.T) {
	err := try.Do(func(attempt int) (bool, error) {
		return true, errors.New("nope")
	})
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if !try.IsMaxRetries(err) {
		t.Error("Expected IsMaxRetries to return true")
	}
}
