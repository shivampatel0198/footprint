# Install Go
VERSION='1.14.3'
OS='linux'
ARCH='amd64'

wget "https://dl.google.com/go/go${VERSION}.${OS}-${ARCH}.tar.gz"
tar -C /usr/local -xzf "go${VERSION}.${OS}-${ARCH}.tar.gz"

# Install Cassandra
sudo apt update
sudo apt install openjdk-8-jdk
sudo apt install apt-transport-https
wget -q -O - https://www.apache.org/dist/cassandra/KEYS | sudo apt-key add -
sudo sh -c 'echo "deb http://www.apache.org/dist/cassandra/debian 311x main" > /etc/apt/sources.list.d/cassandra.list'
sudo apt install cassandra