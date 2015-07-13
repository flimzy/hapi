package hapi

// A common function for various test cases
func (h *HypermediaAPI) TestRegister(ctype,id string) {
    h.Register("GET", "/", ctype, func(c *Context) { c.Stash["id"] = id })
}
