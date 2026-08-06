package service

func KategoriJumlahMinimal(kategori string) int {
	switch kategori {
	case "admin":
		return 20
	case "pejabat", "auditor":
		return 40
	default:
		return 0
	}
}
