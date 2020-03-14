package a

func WithDefer() {
	defer func() { _ = 1 + 1 }()
}

func WithoutDefer() {
	func() { _ = 1 + 1 }()
}
