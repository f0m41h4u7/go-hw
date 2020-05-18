//+build generation

package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

type Validated interface {
	Validate() ([]ValidationError, error)
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func TestUserValidation(t *testing.T) {
	var v interface{} = User{}
	_, ok := v.(Validated)
	require.True(t, ok)

	someUser := User{
		ID:        "441e0ae8644611eab8a0632c74ca9988",
		Name:      "Test",
		Age:       27,
		Email:     "owl@otus.ru",
		Role:      "admin",
		Addresses: []string{RandomString(250)},
	}

	t.Run("ID length", func(t *testing.T) {
		user := someUser
		user.ID = "123"
		errs, _ := user.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "ID")
		require.NotNil(t, errs[0].Err)
	})

	t.Run("email regexp", func(t *testing.T) {
		user := someUser
		user.Email = "isnotvalid@@email"
		errs, _ := user.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "Email")
		require.NotNil(t, errs[0].Err)
	})

	t.Run("age borders", func(t *testing.T) {
		user := someUser
		user.Age = 17
		errs, _ := user.Validate()
		require.NotEqual(t, len(errs), 0)

		for _, a := range []int{18, 34, 50} {
			user.Age = a
			errs, _ := user.Validate()
			require.Len(t, errs, 0)
		}

		user.Age = 51
		errs, _ = user.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "Age")
		require.NotNil(t, errs[0].Err)
	})

	t.Run("role", func(t *testing.T) {
		user := someUser
		user.Role = "manager"
		errs, _ := user.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "Role")
		require.NotNil(t, errs[0].Err)
	})

	t.Run("addresses slice", func(t *testing.T) {
		user := someUser
		user.Addresses = []string{RandomString(250), RandomString(249)}
		errs, _ := user.Validate()
		require.Equal(t, len(errs), 1)
	})
}

func TestAppValidation(t *testing.T) {
	var v interface{} = App{}
	_, ok := v.(Validated)
	require.True(t, ok)

	t.Run("version length", func(t *testing.T) {
		errs, _ := App{"0.1"}.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "Version")
		require.NotNil(t, errs[0].Err)
	})
}

func TestTokenValidation(t *testing.T) {
	var v interface{} = Token{}
	_, ok := v.(Validated)
	require.False(t, ok)
}

func TestResponseValidation(t *testing.T) {
	var v interface{} = Response{}
	_, ok := v.(Validated)
	require.True(t, ok)

	t.Run("code set", func(t *testing.T) {
		for _, c := range []int{200, 404, 500} {
			errs, _ := Response{Code: c}.Validate()
			require.Len(t, errs, 0)
		}

		errs, _ := Response{Code: 133}.Validate()
		require.Equal(t, len(errs), 1)
		require.Equal(t, errs[0].Field, "Code")
		require.NotNil(t, errs[0].Err)
	})
}
