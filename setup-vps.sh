#!/bin/bash

apt -y update
apt install -y openssl

# Helper functions
generate_random_str() {
    LENGTH=12

    random=$(openssl rand -base64 48 | tr -dc 'a-zA-Z0-9' | head -c$LENGTH)

    echo "$random"
}

TZ=Europe/Amsterdam

ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install docker
echo "Installing Docker"
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do apt-get remove $pkg; done

apt-get install -y ca-certificates curl gnupg apt-transport-https make
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null


apt-get -y update

apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-compose
groupadd docker
usermod -aG docker $(whoami)

# Install Nginx
echo "Installing Nginx"
apt install -y nginx
ufw allow 'Nginx Full'

# Install Certbot
echo "Installing certbot"
apt install snapd
snap install --classic certbot
ln -s /snap/bin/certbot /usr/bin/certbot

# Install PHP 8.2
PHP_VERSION=8.2
echo "Installing PHP $PHP_VERSION"
apt install -y lsb-release gnupg2 ca-certificates apt-transport-https software-properties-common
add-apt-repository -y ppa:ondrej/php
apt -y update
apt install -y php$PHP_VERSION php$PHP_VERSION-fpm php$PHP_VERSION-bcmath php$PHP_VERSION-curl php$PHP_VERSION-mbstring php$PHP_VERSION-mysql php$PHP_VERSION-tokenizer php$PHP_VERSION-xml php$PHP_VERSION-zip

# Install Composer 2.6.4
echo "Installing Composer 2.6.4"
php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
php -r "if (hash_file('sha384', 'composer-setup.php') === 'e21205b207c3ff031906575712edab6f13eb0b361f2085f1f1237b7126d785e826a450292b6cfd1d64d92e6563bbde02') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;"
php composer-setup.php
php -r "unlink('composer-setup.php');"
mv composer.phar /usr/local/bin/composer

# Install Go 1.21
echo "Installing Go 1.21"
apt-get -y update
apt-get install -y gnupg curl ca-certificates libcap2-bin libpng-dev build-essential
curl -sL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz | tar -C /usr/local -xz
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc
source ~/.bashrc

# Install Redis
echo "Installing Redis"
apt install -y redis-server
ufw allow 6379

# Setup MySQL
echo "Installing MySQL"
DB_DATABASE="cloud-services-site"
DB_USER=$DB_DATABASE
DB_PASSWORD=$(generate_random_str)

apt install -y mysql-server
systemctl start mysql.service
mysql -h localhost -u $DB_USER -p"$DB_PASSWORD" $DB_DATABASE -e "CREATE USER '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';CREATE DATABASE IF NOT EXISTS $DB_DATABASE;GRANT ALL PRIVILEGES ON $DB_DATABASE.* TO '$DB_USER'@'localhost' WITH GRANT OPTION;FLUSH PRIVILEGES;"

# Setup SFTP
echo "Installing SFTP server"

# Setting up microservices
echo "Installing Microservices"
git clone git@git.noahdev.nl:cloudservices/microservices.git ~/cloud-services/microservices
docker-compose up -d --build || echo "Setup Failed!!!"

# Setting up site
echo "Installing Site"
git clone git@git.noahdev.nl:cloudservices/site.git ~/cloud-services/site
chmod -R 655 storage/
chmod -R 655 bootstrap/
cp .env.production .env
composer install --no-plugins --optimize-autoloader --no-interaction --no-scripts
php artisan key:generate
php artisan test || echo "Setup Failed!!!"
php artisan install:cloud-services --db-database="$DB_DATABASE" --db-password="$DB_PASSWORD" || echo "Setup Failed!!!"

echo "Please enter your domain:"
read DOMAIN

FILES_PATH=/var/www/server/$DOMAIN
mkdir -p $FILES_PATH

nginx_config_template="server {
    listen 80;
    listen [::]:80;
    server_name \$DOMAIN_SERVER_NAME;
    root \$FILES_PATH/public;

    add_header X-Frame-Options \"SAMEORIGIN\";
    add_header X-XSS-Protection \"1; mode=block\";
    add_header X-Content-Type-Options \"nosniff\";

    index index.html index.htm index.php;

    charset utf-8;

    location / {
        try_files \$uri \$uri/ /index.php?\$query_string;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    error_page 404 /index.php;

    location ~ \.php\$ {
        fastcgi_pass unix:/var/run/php/php\$PHP_VERSION-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME \$realpath_root\$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.(?!well-known).* {
        deny all;
    }
}"

nginx_config=$(echo "$nginx_config_template" | \
               sed "s|\$DOMAIN_SERVER_NAME|${DOMAIN}|g; \
                    s|\$FILES_PATH|${FILES_PATH}|g; \
                    s|\$PHP_VERSION|${PHP_VERSION}|g")
echo "$nginx_config" > /etc/nginx/sites-available/server/$DOMAIN
ln -s /etc/nginx/sites-available/server/$DOMAIN /etc/nginx/sites-enabled/$DOMAIN
systemctl restart nginx
nginx -t

# Creating admin site user

echo "Finishing up..."
apt update && apt -y upgrade
apt autoremove