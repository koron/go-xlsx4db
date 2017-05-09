# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "minimal/xenial64"

  config.vm.network :forwarded_port, guest:5432, host:5432

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
    vb.cpus = 2
  end

  config.ssh.insert_key = false

  config.vm.provision "shell", inline: <<-SHELL
    sed -i.bak -e 's!//archive.ubuntu.com/!//ftp.jaist.ac.jp/pub/Linux/!g' /etc/apt/sources.list
    apt update
    apt install -y debconf-utils

    # PostgreSQL
    apt install -y postgresql postgresql-contrib
    sudo -u postgres createuser -DRs vagrant
    sudo -u postgres psql -c "ALTER USER vagrant WITH PASSWORD 'db1234'"
    sudo -u postgres createdb vagrant
    echo "listen_addresses = '*'" >> /etc/postgresql/9.5/main/postgresql.conf
    echo "host vagrant vagrant 0.0.0.0/0 md5" >> /etc/postgresql/9.5/main/pg_hba.conf
    systemctl restart postgresql

    # MySQL
    debconf-set-selections <<< 'mysql-server mysql-server/root_password password mysql123'
    debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password mysql123'
    apt install -y mysql-server

    touch .hushlogin
  SHELL

end
