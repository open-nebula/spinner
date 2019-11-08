IMAGE := nebula-spinner

build:
	cd .. && docker build -t $(IMAGE) -f build/Dockerfile .

run: build
	docker run -it -p 5912:5912 $(IMAGE)

test-send:
	
