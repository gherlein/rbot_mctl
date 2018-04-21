NAME   = rbot_mctl
FILE   = ./${NAME}

build: dependencies 
	go build 

${FILE}: build

dependencies:
	go get ./...

install: ${FILE}
	sudo systemctl stop ${NAME}
	sudo cp ${FILE} /usr/local/bin/
	sudo cp ${FILE}.toml /etc/
	sudo cp ${FILE}.service /etc/systemd/system
	sudo systemctl daemon-reload
	sudo systemctl enable gonetmon
	sudo systemctl start gonetmon

stop:
	sudo service ${NAME} stop

start:
	sudo service ${NAME} start

restart:
	sudo service ${NAME} restart

clean:
	-rm -f ${FILE}
	-rm -f *~

#install_xboxdrv:
#	sudo service xboxdrv stop
#	sudo cp xboxdrv.service /lib/systemd/system/
#	sudo cp xboxdrv /etc/defaults
#	sudo cp uxvars.bak /usr/share/ubuntu-xboxdrv/uxvars
#	sudo cp xboxdrv.ini /etc
#	systemctl daemon-reload
#	sudo service xboxdrv start
#	sudo service xboxdrv status

