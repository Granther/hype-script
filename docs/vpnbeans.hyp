
# Default packages !required!
PACKAGES="openvpn stunnel"
# Package manager command (ie: dnf, apt, yum)
PACKAGETYPE=
# Command used to remove a package for a specific package manager
REMWORD=
# Wether to enable to disable required services (ie: auto start or not). 1=False. 0=True
ENABLE_SRV=1
# For when status is called, used to check for DNS issues
#STATUS_IP="8.8.8.8"
#STATUS_DOMAIN="google.com"
MAX_CURL_TIME=4

#PUB_DNS_SERVER="1.1.1.1"

# For text bolding 
bold=$(tput bold)
normal=$(tput sgr0)

help() {
    echo "Options, must be run as sudo user"
    echo "${bold}sudo ./vpnbeans.sh up${normal}         | bring the vpn UP"
    echo "${bold}sudo ./vpnbeans.sh down${normal}       | bring the vpn DOWN"
    echo "${bold}sudo ./vpnbeans.sh restart${normal}    | bring the vpn UP and then DOWN"
    echo "${bold}sudo ./vpnbeans.sh status${normal}     | see the status of your vpn connection"
    echo "${bold}sudo ./vpnbeans.sh install${normal}    | install vpn configs"
    echo "${bold}sudo ./vpnbeans.sh uninstall${normal}  | uninstall vpn configs"
    echo "${bold}sudo ./vpnbeans.sh help${normal}       | see this nice help message :)"
}

src_env() {
    if test -f "./vpnbeans.env"; then
		source "./vpnbeans.env"
    else
		echo "vpnbeans.env file must be present in the same dir as this script"
		exit
    fi 
}

read_answer() {
    if [ "$1" = "y" ]; then
	return 0
    fi
    if [ "$1" = "n" ]; then
	return 1
    else
	echo "Answer $1 is not valid, either y or n"
	exit
    fi
}

set_package_type() {
	if command -v apt; then 
		echo "Using apt package manager"
		PACKAGETYPE="apt"
		return
	fi 

	if command -v dnf; then
		echo "Using dnf package manager"
		PACKAGETYPE="dnf"
		return
	fi

	if command -v yum; then
		echo "Using yum package manager"
		PACKAGETYPE="yum"
		return
	else
		if [ "$1" = "uninstall" ]; then
			echo -n "This script only manages packages with dnf or apt. Please uninstall theses packages using your desired method: {$bold}$PACKAGES{$normal}"
			return
		fi
		echo -n "This script only installs using dnf or apt package manager. Please install {$bold}$PACKAGES{$normal}. Type 'n' to exit, 'y' if you already have these packages installed. [y/n]: "
		read answer
		if ! read_answer $answer; then
			exit
		fi
	fi
}

#$1 is var name, $2 is old var val, $3 is new val
set_env_var() {
    echo "Changing $1 from $2 to $3"
    if ! $(sed -i "s/$1=$2/$1=$3/g" "./vpnbeans.env"); then
		echo "Unable to change environment variable value in vpnbeans.env, ensure file is present"
		exit
    fi
	src_env # Re-src after setting var to refresh
}

# Check that VPN_SERVER_IP != nil
check_vpn_srv_ip() {
    if [ -z "$VPN_SERVER_IP" ]; then
	echo "Hmm, the VPN_SERVER_IP environment variable is null. This usually means you changed it in the file, unless you saved the VPN_SERVER_IP elsewhere I would reccomend redownloading your VPN configs"
	exit
    fi
    return 0 # return true
}

# Check if GATEWAY_IP is empty, if so, fill it
check_gateway_ip() {
    if [ -z "$GATEWAY_IP" ]; then # If GTWAY == None
		tmp_ip=$(ip route | grep default | awk '{print $3}')
		num_per=$(echo $tmp_ip | tr -dc "." | wc -c) # Number of periods in string
		echo "$num_per"
		if [ $num_per -eq "3" ]; then
			GATEWAY_IP=$tmp_ip
			set_env_var "GATEWAY_IP" "" $tmp_ip # Set new GATEWAY_IP in env
		else # Not enough periods in output
			echo "Unable to auto-retrieve gateway ip, please set this manually in vpnbeans.env. This is usually the local ip address of your router"
			exit
		fi
	fi

	default_gtway=$(ip route | grep default | awk '{print $3}')
	if [ "$default_gtway" != "$GATEWAY_IP" ]; then
		echo -n "WARNING: Your default gateway is $default_gtway, but your config expects $GATEWAY_IP, change to $default_gtway? [y/n]: "
		read answer

		if read_answer $answer; then
			set_env_var "GATEWAY_IP" $GATEWAY_IP $default_gtway # Set new GATEWAY_IP in env	
		fi
	fi

	return 0
}

give_overview() {
    echo "--- Install overview ---"
    echo "2 packages: $PACKAGES"
    echo "2 systemd services added, NOT enabled by default. openvpn-client@client.service & stunnel.service. Pass --enable-srv to enable these (i.e. they start when your computer starts)"
    #echo "2 systemd services added, NOT enabled by default. openvpn-client@client.service & stunnel.service"
    echo -n "Continue? [y/n]: "
    read answer

    if ! read_answer $answer; then 
	exit
    fi
}

set_dns_resolv() {
	if test -f "/etc/resolv.conf"; then
		echo nameserver $1 > "/etc/resolv.conf"
		echo "INFO: Set nameserver to $1 in /etc/resolv.conf"
		#Please run 'sudo ./vpnbeans.sh status' again to see any changes
	else
		echo "/etc/resolv.conf does not exist"
	fi
}

confs_exist() {
	if test -f "ovpn.conf" && test -f "stunnel.pem" && test -f "stunnel.conf"; then
		return 0
	else
		return 1
	fi
}

ensure_stunnel_run() {
	if ! test -d "/var/run/stunnel4"; then
		echo "Creating stunnel dir for .pid file"
		if ! mkdir /var/run/stunnel4; then 
			echo "Unable to create required dir at /var/run/stunnel4"
			exit
		fi
	fi
}

install_confs() {
    echo "Moving stunnel.pem"
    if ! cp "./stunnel.pem" "./stunnel.conf" "/etc/stunnel/"; then
	echo "Unable to copy stunnel.pem & stunnel.conf to /etc/stunnel"
	exit
    fi 

    echo "Moving ovpn.conf"
    if ! cp "./ovpn.conf" "/etc/openvpn/client/client.conf"; then
	echo "Unable to copy ovpn.conf to /etc/openvpn/client/client.conf"
	exit
    fi
}

install_packages() {
    echo "Installing packages: $PACKAGES with $PACKAGETYPE"
    $PACKAGETYPE install -y $PACKAGES
}

finished_installing() {
    echo "--- Finished installing! ---"
    help
}

enable_srv() {
    if [ $ENABLE_SRV -eq 0 ]; then 
		echo "Enabling services..."
		systemctl enable openvpn-client@client
		systemctl enable stunnel
    fi
}

ask_proxy_dns() {
	echo -n "Proxy DNS though VPN server? This helps in preventing DNS leaks. [y/n]: "
	read answer
	if read_answer $answer; then
		set_env_var "PROXY_DNS" "" "YES"
	else
		set_env_var "PROXY_DNS" "" "NO"
	fi
}

# Ensure ovpn.conf & stunnel.pem are adjacent
install() {
    echo "Installing VPN..."

    set_package_type

    if [ "$2" = "--enable-srv" ]; then # Handle --enable-srv being passed
	ENABLE_SRV=0
    fi

    if ! confs_exist; then
	    echo "ERROR: Either ovpn.conf, stunnel.pem and/or stunnel.conf are not adjacent to this script, ensure they are in the same directory as this script"
	    exit
    fi

    give_overview
	ask_proxy_dns
    install_packages $PACKAGES
    enable_srv
    install_confs
    finished_installing
}

uninstall_confs() {
    echo "Removing stunnel.pem & stunnel.conf from /etc/stunnel"
    if ! rm /etc/stunnel/stunnel.{conf,pem}; then
	echo "Unable to remove stunnel.conf and/or stunnel.pem from /etc/stunnel"
	return 1
    fi 

    echo "Removing client.conf from /etc/openvpn/client"
    if ! rm /etc/openvpn/client/client.conf; then
	echo "Unable to remove client.conf from /etc/openvpn/client"
	return 1
    fi
}

uninstall_packages() {
    echo -n "Would you like to remove the packages: $PACKAGES? [y/n]: "
    read answer

    if check_answer $answer; then
	$PACKAGETYPE remove -y $PACKAGES
    fi

    if ! check_anwer $answer; then
	exit
    fi
}

uninstall() {
    echo "Uninstalling VPN..."

    set_package_type "uninstall"

    if ! uninstall_confs; then 
	echo -n "Having trouble removing some files. This might because the files are already removed. continue? [y/n]: "
	read answer
	if ! read_answer $answer; then
	    exit
	fi
    fi

    uninstall_packages
    echo "Successfully uninstalled VPN, you may now delete this script and any files downloaded with it :)"
}

add_routes() {
    ROUTE="$VPN_SERVER_IP/32 via $GATEWAY_IP"
    echo "Adding route $ROUTE"
    ip route add $ROUTE
}

del_routes() {
    ROUTE="$VPN_SERVER_IP/32 via $GATEWAY_IP"
    echo "Removing route $ROUTE"
    ip route del $ROUTE
}

proxy_dns() {
    if [ "$PROXY_DNS" = "YES" ]; then
		return 0
    else
		return 1
	fi
}

set_dns() {
	if proxy_dns; then
		set_dns_resolv $VPN_DNS_IP	
	fi
}

rm_dns() {
	if proxy_dns; then
		set_dns_resolv $NON_VPN_DNS_IP
	fi
}

up() {
    check_vpn_srv_ip
    check_gateway_ip

	set_dns

    start_stunnel
    sleep 1s
    add_routes
    start_openvpn
		
	echo -n "Giving things a second to settle"
	sleep 1s
	echo -n "."
	sleep 1s
	echo -n "."
	sleep 1s
	echo "."

	status
}

down() {
    stop_openvpn
    sleep 1s
    del_routes
    stop_stunnel

	rm_dns
}

status() {
	echo -n "STATUS: Attempting test connection to domain: $DNS_CHECK_DOMAIN..."
	curl -4 -s --connect-timeout $MAX_CURL_TIME --max-time $MAX_CURL_TIME $DNS_CHECK_DOMAIN > /dev/null
	if [ $? -eq 0 ]; then # Can curl google
		echo "Good :)"
	else
	#if ! [ $? -eq 0 ]; then # Cant curl google
		echo "Bad :(" # Newline after status hasn't
		echo "ERROR: DNS Check failed to connect to domain $DNS_CHECK_DOMAIN"
		
		echo -n "STATUS: Attempting test ping to ip: $PING_CHECK_IP..."
		ping -c 1 -w 2 $PING_CHECK_IP > /dev/null
		if [ $? -eq 0 ]; then # Can ping google ip
			echo "Good, but..." # Newline after status hasn't
			echo "ERROR: Was able to ping known IP ($PING_CHECK_IP) but not domain name ($DNS_CHECK_DOMAIN), you, my friend, have a DNS issue"
			if ! proxy_dns; then # If DNS is not proxied, ask if they would like to
				ask_proxy_dns
			fi
			exit 1
		else
			echo "Bad :(" # Newline after status hasn't
			echo "ERROR: DNS Check, failed to reach internet with target test IP $PING_CHECK_IP"
			exit 1
		fi
	fi

	sleep 1
    pub_ip=$(curl -4 -s ip.me)

	#if [ "$CURRENT_DNS_IP" = "$VPN_DNS_IP" ]; then
	#	proxy_dns="Yes"
	#else
	#	proxy_dns="No"
	#fi

    if [ "$pub_ip" = "$VPN_SERVER_IP" ]; then
		echo "--- CONNECTED OVER VPN ---"
		echo "Current Public IP: ${bold}"$pub_ip"${normal}"
		echo "Proxying DNS: $PROXY_DNS"
		#echo "Primary DNS Server: ${bold}"$CURRENT_DNS_IP"${normal}"
		#echo "Proxying DNS: ${bold}"$proxy_dns"${normal}"
		exit 0
    else
		echo "--- DISCONNECTED FROM VPN ---"
		echo "Current Public IP: ${bold}"$pub_ip"${normal}"
		echo "Proxying DNS: NO"
		#echo "Primary DNS Server: ${bold}"$CURRENT_DNS_IP"${normal}" 
		#echo "Proxying DNS: ${bold}"$proxy_dns"${normal}"
		exit 0
    fi
    # If tun exists, openvpn-client@client & stunnel is running without issue
    # Show ip addr, server public ip
}

restart() {
    down
    up
}

enable_stunnel() {
    echo "Enabling Stunnel"
    systemctl enable stunnel.service
}

start_stunnel() {
    echo "Starting Stunnel"
    ensure_stunnel_run
    systemctl start stunnel.service
}

stop_stunnel() {
    echo "Stopping Stunnel"
    systemctl stop stunnel.service
}

start_openvpn() {
    echo "Starting OpenVPN"
    systemctl start openvpn-client@client.service
}

stop_openvpn() {
    echo "Stopping OpenVPN"
    systemctl stop openvpn-client@client.service
}

# source vpnbeans.env file, containing VPN server ip and gatway ip
src_env

#if test -z "$1"; then
#	help
#fi

if [ $# -gt 0 ]; then
	case $1 in
	    "help")
	        help
		;;
	    "up")
	        up
		;;
	    "down")
		down
		;;
	    "status")
		status
		;;
	    "restart")
		restart
		;;
	    "install")
		install "$@"
		;;
	    "uninstall")
		uninstall "$@"
		;;
	    	*)
		echo "Invalid option: $1"
		help
		;;
	esac
fi

