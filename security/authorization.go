package security

// UserDetails provides information regarding the user
// providers can create their own interfaces for what they need.
type UserDetails interface{
  Verify(interface{})      bool
  SatisfiesRole(Role)      bool     //Whether or not the role is satisfied
}

// Role describes a role
type Role interface{
  GetRole() string
}
