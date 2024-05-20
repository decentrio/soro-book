# Install Stellar Core on Ubuntu

From [official Github repo](https://github.com/stellar/packages)
```
wget -qO - https://apt.stellar.org/SDF.asc | sudo apt-key add -
echo "deb https://apt.stellar.org $(lsb_release -cs) stable" | sudo tee -a /etc/apt/sources.list.d/SDF.list
sudo apt-get update && sudo apt-get install stellar-core
```

# Install and run Sorobook

Clone [Sorobook Github repo](https://github.com/decentrio/soro-book) and build
```
git clone https://github.com/decentrio/soro-book
cd soro-book
make install
```

Config Stellar Core binary directory
```
nvim soro-book/config/config.go
```

Add Postgres Url config and run Sorobook directly
```
export POSTGRES_URL=postgresql://YOUR_USER:YOUR_PASSWORD@POSTGRES_IP:5432/DATABASE_NAME
sorobook start
```

Or run Sorobook as a service
```
sudo tee <<EOF >/dev/null /etc/systemd/system/sorobook.service
[Unit]
Description=sorobook

[Service]
WorkingDirectory=/root
Environment=POSTGRES_URL="postgresql://YOUR_USER:YOUR_PASSWORD@POSTGRES_IP:5432/DATABASE_NAME"
ExecStart=$(which sorobook) start
Type=simple
Restart=always
RestartSec=5
User=root

[Install]
WantedBy=multi-user.target
EOF
```
```
systemctl daemon-reload
systemctl start sorobook.service
```

