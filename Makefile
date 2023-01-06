install:
	curl https://dl.google.com/go/go1.19.4.linux-amd64.tar.gz --output go1.19.4.linux-amd64.tar.gz
	rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.4.linux-amd64.tar.gz
	export PATH=$PATH:/usr/local/go/bin
run:
	go run main.go