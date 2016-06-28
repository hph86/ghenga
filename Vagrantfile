# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure(2) do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = ENV["GHENGA_VAGRANT_BOX"] || "ubuntu/wily64"

  # fix slow network
  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--nictype1", "virtio"]
    v.memory = 1024
  end


  forwarded_port = ENV["GHENGA_PORT"] || 8080
  config.vm.network "forwarded_port", guest: 8080, host: forwarded_port.to_i
  config.vm.synced_folder ".", "/home/vagrant/ghenga"

  config.vm.provision "shell", inline: <<-SHELL
       curl -s https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add -
       echo 'deb https://deb.nodesource.com/node_4.x wily main' > /etc/apt/sources.list.d/nodesource.list
       echo 'deb-src https://deb.nodesource.com/node_4.x wily main' >> /etc/apt/sources.list.d/nodesource.list

       export DEBIAN_FRONTEND=noninteractive
       apt-get update
       apt-get -y dist-upgrade

       apt-get install -y \
          -o Dpkg::Options::="--force-confdef" \
          -o Dpkg::Options::="--force-confnew" \
          curl wget git vim tmux screen zsh moreutils silversearcher-ag nodejs postgresql

       locale-gen -a de_DE.UTF-8 en_US.UTF-8 en_GB.UTF-8
       grep -q LC_ALL /etc/environment || echo -en 'LC_ALL=en_US.UTF-8\nLANG=en_US.UTF-8\n' >> /etc/environment

       # create database 'vagrant'
       echo "create user vagrant with encrypted password 'vagrant';" | sudo -u postgres psql
       echo "create database vagrant with owner vagrant" | sudo -u postgres psql
       echo "create database test with owner vagrant" | sudo -u postgres psql
  SHELL

  config.vm.provision :reload

  config.vm.provision "shell", :privileged => false, inline: <<-SHELL
       wget -q -O /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
       mkdir -p .local
       cd .local/
       tar xzf /tmp/go.tar.gz

       cat /etc/skel/.profile > ~/.profile
       echo 'export GOROOT=$HOME/.local/go' >> ~/.profile
       echo 'export GOPATH=$HOME/go' >> ~/.profile
       echo 'export GOBIN=$HOME/bin' >> ~/.profile
       echo 'export PATH=$PATH:$GOROOT/bin:$GOBIN' >> ~/.profile

       source ~/.profile

       go get github.com/rubenv/modl-migrate/...
       go get github.com/constabulary/gb/...
       go get github.com/derekparker/delve/cmd/dlv
  SHELL

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Define a Vagrant Push strategy for pushing to Atlas. Other push strategies
  # such as FTP and Heroku are also available. See the documentation at
  # https://docs.vagrantup.com/v2/push/atlas.html for more information.
  # config.push.define "atlas" do |push|
  #   push.app = "YOUR_ATLAS_USERNAME/YOUR_APPLICATION_NAME"
  # end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  # config.vm.provision "shell", inline: <<-SHELL
  #   sudo apt-get update
  #   sudo apt-get install -y apache2
  # SHELL
end
