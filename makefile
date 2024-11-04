

doc:
	# echo "note: opens http://localhost:8080/ in default browser"
	pkgsite -open .
	


test:
	go test ./... -v


