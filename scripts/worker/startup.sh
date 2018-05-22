# set timezone to UTC
ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

# create elasticshift user/group
groupadd --gid 1002 elasticshift \
	&& useradd --uid 1002 --gid elasticshift --shell /bin/bash --create-home elasticshift \
	&& echo 'elasticshift ALL=NOPASSWD: ALL' >> /etc/sudoers.d/elasticshift

su - elasticshift

# download worker


# start worker as entry point
