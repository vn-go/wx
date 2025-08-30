package wx

func Routes(baseUri string, ins ...any) error {
	return utils.Routes.Add(baseUri, ins...)
}
