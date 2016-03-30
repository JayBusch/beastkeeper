#!/bin/sh 
rec=false
stop=false
case "$1" in
	'record') 
		rec=true;;
	'stop')
		stop=true;;
esac

if "$rec" eq true
then
	if [ -f "img/animatedBuild.gif"] 
	then
		rm img/animatedBuild.gif		
	fi
	ttyrec -e "./buildAndTest.sh stop" buildAnimation.rec
	ttygif buildAnimation.rec
	convert -delay 10 -loop 1 *.gif ./img/animatedBuild.gif
	rm *.gif
	rm buildAnimation.rec


else
	cd bin
	go build ../src/bk
	go build ../src/test
	cd ..
	go test -coverprofile=coverage.out -v ./src/bk 
	go test -coverprofile=states_coverage.out -v ./src/bk/states
	go tool cover -html=coverage.out
	go tool cover -html=states_coverage.out
	./bin/test
fi

if "$stop" eq true
then
	exit
fi

