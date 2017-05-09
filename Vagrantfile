# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "minimal/xenial64"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
    vb.cpus = 2
  end

  config.ssh.insert_key = false

  config.vm.provision "shell", inline: <<-SHELL
    sed -i.bak -e 's!//archive.ubuntu.com/!//ftp.jaist.ac.jp/pub/Linux/!g' /etc/apt/sources.list
    apt update
    apt install -y debconf-utils
    debconf-set-selections <<< 'mysql-server mysql-server/root_password password mysql123'
    debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password mysql123'
    apt install -y mysql-server
    apt install -y postgresql postgresql-contrib
    touch .hushlogin
  SHELL
end
