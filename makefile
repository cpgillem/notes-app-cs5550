rundev:
	app/app 8080

keygen:
	rm -f app/keys/*
	ssh-keygen -t rsa -b 4096 -f app/keys/app.rsa
	openssl rsa -in app/keys/app.rsa -pubout -outform PEM -out app/keys/app.rsa.pub
	# https://gist.github.com/ygotthilf/baa58da5c3dd1f69fae9
