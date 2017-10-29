package main

type User struct {
	Resource
	Name string
	Admin bool
}

func (u *User) Load() error {
	return u.Select([]string{"name", "admin"}, &u.Name, &u.Admin)
}

func (u *User) Save() error {
	return u.Sync([]string{"name", "admin"}, u.Name, u.Admin)
}

type Note struct {
	Resource
	Title string
	Content string
	Time string
	UserID int64
}

func (n *Note) Load() error {
	return n.Select([]string{"title", "content", "time", "user_id"}, &n.Title, &n.Content, &n.Time, &n.UserID)
}

func (n *Note) Save() error {
	return n.Sync([]string{"title", "content", "time", "user_id"}, n.Title, n.Content, n.Time, n.UserID)
}

func (n *Note) User() (u User, err error) {
	// Create an unloaded model for the user.
	u = User {
		Resource: Resource {
			ID: n.UserID,
			DB: n.DB,
			Table: "users",
		},
	}

	// Define err as the result of loading the user from their ID.
	err = u.Load()

	return
}

type Tag struct {
	Resource
	Title string
}

func (t *Tag) Load() error {
	return t.Select([]string{"title"}, &t.Title)
}

func (t *Tag) Save() error {
	return t.Sync([]string{"title"}, t.Title)
}
