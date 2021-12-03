package data

import "testing"

func TestCheckUserValidation(t *testing.T) {

	u := &User{
		FirstName: "Dush",
		LastName: "Desh",
		Email : "dushdesh@email.com",
	}

	err := u.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
