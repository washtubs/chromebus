set -e
cd $GOPATH/src/github.com/washtubs/chromebus/cmd/
go install
cd -

echo sudo ln -s -T $GOPATH/bin/cmd /usr/bin/chromebus
echo sudo ln -s -T $PWD/systemd/chromebus.service /usr/lib/systemd/system/chromebus.service
[ -f /usr/bin/chromebus ] && sudo unlink /usr/bin/chromebus
[ -f /usr/bin/chromebus ] && sudo unlink /usr/lib/systemd/system/chromebus.service
sudo ln -s -T $GOPATH/bin/cmd /usr/bin/chromebus
sudo ln -s -T $PWD/systemd/chromebus.service /usr/lib/systemd/system/chromebus.service
