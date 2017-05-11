# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "minimal/xenial64"

  # for PostgreSQL
  config.vm.network :forwarded_port, guest:5432, host:5432
  # for MySQL
  config.vm.network :forwarded_port, guest:3306, host:3306

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
    vb.cpus = 2
  end

  config.ssh.insert_key = false

  config.vm.provision "shell", inline: <<-SHELL
    update-locale LANG=C.UTF-8 LANGUAGE=
    sed -i.bak -e 's!//archive.ubuntu.com/!//ftp.jaist.ac.jp/pub/Linux/!g' /etc/apt/sources.list
    apt update
    apt install -y debconf-utils

    # PostgreSQL
    apt install -y postgresql postgresql-contrib
    sudo -u postgres createuser -DRs vagrant
    sudo -u postgres psql -c "ALTER USER vagrant WITH PASSWORD 'db1234'"
    sudo -u postgres createdb -O vagrant vagrant
    echo "listen_addresses = '*'" >> /etc/postgresql/9.5/main/postgresql.conf
    echo "host vagrant vagrant 0.0.0.0/0 md5" >> /etc/postgresql/9.5/main/pg_hba.conf
    systemctl restart postgresql

    # MySQL
    debconf-set-selections <<< 'mysql-server mysql-server/root_password password mysql123'
    debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password mysql123'
    apt install -y mysql-server
    echo "default-character-set=utf8" >> /etc/mysql/conf.d/mysql.cnf
    echo "default-character-set=utf8" >> /etc/mysql/conf.d/mysqldump.cnf
    cat >> /etc/mysql/mysql.conf.d/mysqld.cnf <<__EOS__
bind-address=0.0.0.0
character-set-server=utf8
skip-character-set-client-handshake
default-storage-engine=INNODB
__EOS__
    systemctl restart mysql
    mysql -u root --password=mysql123 <<__EOS__
CREATE DATABASE vagrant;
GRANT ALL ON vagrant.* TO vagrant@"%" IDENTIFIED BY 'db1234';
FLUSH PRIVILEGES;
__EOS__

    touch .hushlogin
  SHELL

end
