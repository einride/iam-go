package iamgooglemember

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestUserInfo_UnmarshalBase64(t *testing.T) {
	var ui UserInfo
	assert.NilError(t, ui.UnmarshalBase64(`eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI0Nzk0NzQwMjEzNi1samxiNDAyaXA4MnQ5MnNvYzc0Y29mMnQ5ZTdrdm4ycC5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6IjQ3OTQ3NDAyMTM2LWxqbGI0MDJpcDgydDkyc29jNzRjb2YydDllN2t2bjJwLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTA0MDIyOTgxMjU1Nzk5OTk3NzgxIiwiaGQiOiJlaW5yaWRlLnRlY2giLCJlbWFpbCI6ImpvaGFubmVzLndpdHRlbnN0YW1AZWlucmlkZS50ZWNoIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJaS01zbExUTnEzRFJfamlxM0w5QnVnIiwibmFtZSI6IkpvaGFubmVzIFdpdHRlbnN0YW0iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EtL0FPaDE0R2dfblJHNFcwS1E3LVRTdlpuUWdzLTlfRmVfMmNLT1NkaVppbVVaPXM5Ni1jIiwiZ2l2ZW5fbmFtZSI6IkpvaGFubmVzIiwiZmFtaWx5X25hbWUiOiJXaXR0ZW5zdGFtIiwibG9jYWxlIjoiZW4iLCJpYXQiOjE2MjIxNzgzMjIsImV4cCI6MTYyMjE4MTkyMn0=`))
}
